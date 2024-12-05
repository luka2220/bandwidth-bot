package tokenbucket

import (
	"log"
	"net/http"
	"os"
)

// TODO: Need to find a way to create the in-memory store abstracted way from user
var (
	ipAdderStore = make(map[string]Bucket)
	logger       = log.New(os.Stdout, "[SERVER-TOKEN]: ", log.LstdFlags)
)

type Bucket struct {
	IpAdder    string // Ip adderess of the current bucket
	capacity   int    // Maximum size of tokens
	size       int    // Current size of of bucket
	fillRate   int    // amount of tokens to add to the bucket per second
	httpStatus int    // HTTP status of the bucket
}

// Gets the buck associasted with the IP address, if none exists in the map a new but is created and stored
func GetIpAdderBucket(ip string) *Bucket {
	b, ok := ipAdderStore[ip]
	if ok {
		// TODO: Since the ip exists, add logic for removing a token from the bucket
		logger.Printf("ip address in memory (%s)\n", b.IpAdder)
		return &b
	}

	newBucket := &Bucket{
		IpAdder:  ip,
		capacity: 10,
		size:     0,
		fillRate: 1,
	}
	ipAdderStore[ip] = *newBucket
	logger.Printf("ip adddress created in memory (%s)\n", newBucket.IpAdder)

	return newBucket
}

// Get the http status code of the current bucket
func (b *Bucket) GetHTTPStatus() int {
	return http.StatusOK
}

// Gets the current token size of the bucket
func (b *Bucket) GetBucketSize() int {
	return b.size
}

// Consistently adds a token to the bucket at a rate of Bucket.fillRate
func (b *Bucket) addToken() {

}

// Removes a token from the bucket
func (b *Bucket) removeToken() {

}
