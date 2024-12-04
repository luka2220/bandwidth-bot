package tokenbucket

import (
	"log"
	"os"
)

// TODO: Need to find a way to create the in-memory store abstracted way from user
var (
	ipAdderStore map[string]Bucket
	logger       = log.New(os.Stdout, "[SERVER-TOKEN]: ", log.LstdFlags)
)

type Bucket struct {
	ipAdder  string // Ip adderess of the current bucket
	capacity int    // Maximum size of tokens
	size     int    // Current size of of bucket
	fillRate int    // amount of tokens to add to the bucket per second
}

func GetBucket(ip string) *Bucket {

	b, ok := ipAdderStore[ip]
	if ok {
		logger.Println("ip address in memory")
		return &b
	}

	newBucket := &Bucket{
		ipAdder:  ip,
		capacity: 10,
		size:     0,
		fillRate: 1,
	}
	ipAdderStore[ip] = *newBucket
	logger.Println("ip adddress created in memory")

	return newBucket
}

// Gets the current token size of the bucket
func (b *Bucket) GetBucketSize() int {
	return b.size
}

// Consistently adds a token to the bucket at a rate of Bucket.fillRate
func (b *Bucket) addToken() {

}

// Removes a token from the bucket
func (b *Bucket) RemoveToken() {

}
