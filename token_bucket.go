package bandwidthbot

import (
	"log"
	"net/http"
	"os"
	"time"
)

var (
	ipAdderStore = make(map[string]*Bucket)
	loggerTB     = log.New(os.Stdout, "[SERVER-TOKEN]: ", log.LstdFlags)
)

type Bucket struct {
	IpAdder     string // Ip adderess of the current bucket
	capacity    int    // Maximum size of tokens
	tokenAmount int    // Current size of of bucket
	fillRate    int    // Amount of tokens to fill the bucket with per second
	removeRate  int    // Amopunt of tokens to remove from the bucket when a request is recieved
	httpStatus  int    // HTTP status of the bucket
}

func InitializeTokenBucket(ip string) *Bucket {
	b, ok := ipAdderStore[ip]
	if ok {
		if b.tokenAmount <= 0 {
			loggerFW.Printf("%s bucket size is too low (%d) - bad request\n", b.IpAdder, b.tokenAmount)
			ipAdderStore[ip].httpStatus = http.StatusTooManyRequests
		} else if b.tokenAmount > 0 {
			ipAdderStore[ip].httpStatus = http.StatusOK
			ipAdderStore[ip].tokenAmount -= b.removeRate
		}

		loggerFW.Printf("ip address in memory (%s) has (%d) tokens in bucket\n", b.IpAdder, b.tokenAmount)
		return b
	}

	newBucket := &Bucket{
		IpAdder:     ip,
		capacity:    10,
		fillRate:    1,
		removeRate:  1,
		tokenAmount: 10,
		httpStatus:  http.StatusOK,
	}

	go newBucket.fillBucket()

	ipAdderStore[ip] = newBucket
	loggerFW.Printf("ip adddress created in memory (%s)\n", newBucket.IpAdder)

	return newBucket
}

func (b *Bucket) fillBucket() {
	loggerFW.Printf("Fill bucket process started for %s\n", b.IpAdder)

	ticker := time.NewTicker(time.Second)
	start := time.Now()

	go func() {
		for {
			select {
			case <-ticker.C:
				timeOfFillBucketOperation := time.Since(start).Seconds()
				loggerFW.Printf("%s running fill bucket for %.2f seconds\n", b.IpAdder, timeOfFillBucketOperation)

				if b.tokenAmount < b.capacity {
					ipAdderStore[b.IpAdder].tokenAmount += b.fillRate
				}

				if timeOfFillBucketOperation > 60.00 {
					loggerFW.Printf("removing %s from memory\n", b.IpAdder)
					delete(ipAdderStore, b.IpAdder)
					return
				}
			}
		}
	}()
}

// Get the http status code of the current bucket
func (b *Bucket) GetHTTPStatus() int {
	return b.httpStatus
}
