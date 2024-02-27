package rolling

type Bucket struct {
	// The sum of the values in the bucket.
	sum float64
	// The number of values in the bucket.
	count uint64
}

// Add adds a value to the bucket.
func (b *Bucket) Add(value float64) {
	b.sum += value
	b.count++
}

// Avg returns the average of the values in the bucket.
func (b *Bucket) Avg() float64 {
	if b.count == 0 {
		return 0
	}
	return b.sum / float64(b.count)
}

// NewBucket returns a new bucket.
func NewBucket() *Bucket {
	return &Bucket{}
}
