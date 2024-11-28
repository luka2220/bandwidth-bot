package tokenbucket

type Bucket struct {
	capacity int	// Maximum size of tokens 
	id string 		// ID of the bucket
	size int		// Current size of iof bucket
	fillRate int	// amount of tokens to add to the bucket per second
}

func newBucket(id string) *Bucket {
	return &Bucket{
		capacity: 10,
		id: id,
		size: 0,
		fillRate: 1,
	}
}

// Consistently adds a token to the bucket at a rate of Bucket.fillRate
func (b *Bucket) addToken() {

}