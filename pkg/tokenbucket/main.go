package tokenbucket

type Bucket struct {
	capacity int // Maximum size of tokens
	size     int // Current size of of bucket
	fillRate int // amount of tokens to add to the bucket per second
}

func NewBucket() *Bucket {
	return &Bucket{
		capacity: 10,
		size:     0,
		fillRate: 1,
	}
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
