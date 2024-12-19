package bandwidthbot

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	tokenBucketStore  = make(map[string]*tokenBucket)
	loggerTB          = log.New(os.Stdout, "[SERVER-TOKEN]: ", log.LstdFlags)
	tokenBucketMutex  sync.Mutex
	tokenBucketExipry = 60 * time.Second // Maximum time to keep and ip in ipAdderStore
)

type tokenBucket struct {
	ipAddress   string // Ip adderess of the current bucket
	capacity    int    // Maximum size of tokens
	tokenAmount int    // Current size of of bucket
	fillRate    int    // Amount of tokens to fill the bucket with per second
	removeRate  int    // Amopunt of tokens to remove from the bucket when a request is recieved

	lastRequest time.Time // Time of the last request made
	startTime   time.Time // Time of when the token bucket start for current ip
}

// RunTokenBucket execute the token bucket algorithm
// for an ip address for rate limiting. It returns the
// http status code of the operation. i.e 429, 200 ...
func RunTokenBucket(ip string) int {
	tokenBucketMutex.Lock()
	defer tokenBucketMutex.Unlock()

	removeExpiredIpTokenBucket()

	bucket, exists := tokenBucketStore[ip]
	if !exists {
		tokenBucketStore[ip] = &tokenBucket{
			ipAddress:   ip,
			capacity:    10,
			tokenAmount: 10,
			fillRate:    1,
			removeRate:  1,

			lastRequest: time.Now(),
			startTime:   time.Now(),
		}

		return http.StatusOK
	}

	if bucket.tokenAmount > 0 {
		tokenBucketStore[ip].tokenAmount -= bucket.removeRate
	}

	if bucket.tokenAmount < bucket.capacity {
		amountToFillBucket := int(time.Since(bucket.lastRequest).Seconds()) * bucket.fillRate

		if amountToFillBucket+bucket.capacity > bucket.capacity {
			tokenBucketStore[bucket.ipAddress].tokenAmount = bucket.capacity
		} else {
			tokenBucketStore[bucket.ipAddress].tokenAmount += amountToFillBucket
		}
	}

	tokenBucketStore[ip].lastRequest = time.Now()

	if bucket.tokenAmount <= 0 {
		loggerTB.Printf("%s bucket size is too low (%d) - bad request\n", bucket.ipAddress, bucket.tokenAmount)
		return http.StatusTooManyRequests
	}

	loggerTB.Printf("%s bucket\n\trunning for %.2f seconds\n\ttoken amount %d", bucket.ipAddress, time.Since(bucket.startTime).Seconds(), bucket.tokenAmount)

	return http.StatusOK
}

func removeExpiredIpTokenBucket() {
	now := time.Now()
	for ip, bucket := range tokenBucketStore {
		if now.Sub(bucket.lastRequest) > tokenBucketExipry {
			delete(tokenBucketStore, ip)
			loggerTB.Printf("%s time expired, removing from memory\n", bucket.ipAddress)
		}
	}
}
