package rolling

// 定义 Bucket 结构体
// Define the Bucket struct
type Bucket struct {
	sum   float64 // 存储值的总和 Sum of the values
	count uint64  // 存储值的数量 Number of values
}

// Reset 方法重置 Bucket 的 sum 和 count
// The Reset method resets the sum and count of the Bucket
func (b *Bucket) Reset() {
	b.sum = 0
	b.count = 0
}

// Add 方法将一个值添加到 Bucket 的 sum，并增加 count
// The Add method adds a value to the sum of the Bucket and increases the count
func (b *Bucket) Add(value float64) {
	b.sum += value
	b.count++
}

// Count 方法返回 Bucket 中的 count
// The Count method returns the count of the Bucket
func (b *Bucket) Count() uint64 {
	return b.count
}

// Avg 方法返回 Bucket 中的平均值，如果 count 为 0，则返回 0
// The Avg method returns the average value of the Bucket, returns 0 if count is 0
func (b *Bucket) Avg() float64 {
	// 如果 count 为 0，则返回 0
	// Return 0 if count is 0
	if b.count == 0 {
		return 0
	}

	// 返回 sum 除以 count
	// Return sum divided by count
	return b.sum / float64(b.count)
}

// Sum 方法返回 Bucket 中的 sum
// The Sum method returns the sum of the Bucket
func (b *Bucket) Sum() float64 {
	return b.sum
}

// NewBucket 函数创建并返回一个新的 Bucket 实例
// The NewBucket function creates and returns a new instance of Bucket
func NewBucket() *Bucket {
	return &Bucket{}
}
