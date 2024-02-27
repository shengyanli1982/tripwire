package rolling

type Bucket struct {
	// The sum of the values in the bucket.
	sum float64
	// The number of values in the bucket.
	count uint64
}

// Reset resets the bucket.
func (b *Bucket) Reset() {
	b.sum = 0
	b.count = 0
}

// Add adds a value to the bucket.
func (b *Bucket) Add(value float64) {
	b.sum += value
	b.count++
}

// Count returns the number of values in the bucket.
func (b *Bucket) Count() uint64 {
	return b.count
}

// Avg returns the average of the values in the bucket.
func (b *Bucket) Avg() float64 {
	// If the bucket is empty, return 0.
	if b.count == 0 {
		return 0
	}

	// Return the average of the values in the bucket.
	return b.sum / float64(b.count)
}

// Count returns the number of values in the bucket.
func (b *Bucket) Sum() float64 {
	return b.sum
}

// NewBucket returns a new bucket.
func NewBucket() *Bucket {
	return &Bucket{}
}
