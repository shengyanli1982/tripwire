package rolling

import (
	"sync"
	"time"

	com "github.com/shengyanli1982/tripwire/common"
)

const (
	// 默认的滚动窗口大小。1个插槽等于1秒，10个插槽等于10秒。
	// The default size of the rolling window. 1 slot equals 1s, 10 slots equal 10s.
	DefaultRollingWindowSize = 10

	// 滚动窗口中每个插槽的默认持续时间。
	// The default duration of each slot in the rolling window.
	DefaultRollingWindowSlotInterval = time.Millisecond * 500
)

var (
	// 最小滚动窗口大小为2个插槽，即2秒。
	// The minimum rolling window size is 2 slots, i.e., 2s.
	minRollingWindowSize = 2

	// 最大滚动窗口大小为600个插槽，即10分钟。
	// The maximum rolling window size is 600 slots, i.e., 10 minutes.
	maxRollingWindowSize = 10 * 60
)

// RollingWindow 是一个滚动窗口。
// RollingWindow is a rolling window.
type RollingWindow struct {
	// 环形缓冲区。
	// The ring buffer.
	ring *Ring

	// 滚动窗口的大小。
	// The size of the rolling window.
	size int

	// 滚动窗口的偏移量。写入索引。
	// The offset of the rolling window. Writing index.
	offset int

	// 每个插槽的持续时间。
	// The duration of each slot.
	interval time.Duration

	// 滚动窗口插槽最后更新的时间。
	// The time when the rolling window slot was last updated.
	updateAt time.Time

	// 保护滚动窗口的互斥锁。
	// The mutex to protect the rolling window.
	lock sync.Mutex

	// 指示滚动窗口是否正在运行的标志。
	// The flag to indicate if the rolling window is running.
	runing bool

	// sync.Once 以确保滚动窗口只停止一次。
	// The sync.Once to ensure that the rolling window is stopped only once.
	once sync.Once
}

// NewRollingWindow 返回一个具有指定大小和插槽持续时间的新滚动窗口。
// NewRollingWindow returns a new rolling window with the specified size and slot duration.
func NewRollingWindow(size int) *RollingWindow {
	// 如果大小小于最小大小或大于最大大小，则使用默认大小。
	// If the size is less than the minimum size or greater than the maximum size, use the default size.
	if size < minRollingWindowSize || size > maxRollingWindowSize {
		size = DefaultRollingWindowSize
	}

	// 计算滚动窗口中的插槽数量。
	// Calculate the number of slots in the rolling window.
	slotCount := size * int(time.Second/DefaultRollingWindowSlotInterval)

	// 创建并返回滚动窗口。
	// Create and return the rolling window.
	rw := RollingWindow{
		ring:     NewRing(slotCount),
		size:     slotCount,
		interval: DefaultRollingWindowSlotInterval,
		lock:     sync.Mutex{},
		runing:   true,
		once:     sync.Once{},
		updateAt: time.Now(),
	}

	// 初始化滚动窗口。
	// Initialize the rolling window.
	for i := 0; i < rw.ring.Cap(); i++ {
		rw.ring.Push(NewBucket())
	}

	// 返回滚动窗口。
	// Return the rolling window.
	return &rw
}

// Stop 停止滚动窗口。
// Stop stops the rolling window.
func (w *RollingWindow) Stop() {
	w.once.Do(func() {
		w.lock.Lock()
		w.runing = false
		w.ring.Reset()
		w.lock.Unlock()
	})
}

// span 返回自滚动窗口最后更新以来经过的插槽数量。
// span returns the number of slots that have elapsed since the rolling window was last updated.
func (w *RollingWindow) span() int {
	offset := int(time.Since(w.updateAt) / w.interval)
	if offset >= 0 && offset < w.size {
		return offset
	}
	return w.size
}

// updateOffset 更新滚动窗口。
// updateOffset updates the rolling window.
func (w *RollingWindow) updateOffset() {
	// 计算自滚动窗口最后更新以来经过的插槽数量。
	// Calculate the number of slots that have elapsed since the rolling window was last updated.
	n := w.span()

	// 如果滚动窗口没有被更新，返回。
	// If the rolling window has not been updated, return.
	if n <= 0 {
		return
	}

	// 获取当前的偏移量。
	// Get the current offset.
	offset := w.offset

	// 如果滚动窗口已经向前移动，重置已经过去的插槽。
	// If the rolling window has been moved forward, reset the slots that have elapsed.
	for i := 1; i <= n; i++ {
		bucket := w.ring.At((offset + i) % w.size).(*Bucket)
		bucket.Reset()
	}

	// 更新滚动窗口的偏移量。
	// Update the offset of the rolling window.
	w.offset = (offset + n) % w.size

	// 更新滚动窗口插槽最后更新的时间。
	// Update the time when the rolling window slot was last updated.
	now := time.Now().UnixNano()
	w.updateAt = time.Unix(0, now-(now%int64(w.interval)))
}

// Add 向滚动窗口添加一个值。
// Add adds a value to the rolling window.
func (w *RollingWindow) Add(value float64) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// 如果滚动窗口没有运行，返回一个错误。
	// If the rolling window is not running, return an error.
	if !w.runing {
		return com.ErrorRollingWindowStopped
	}

	// 更新滚动窗口。
	// Update the rolling window.
	w.updateOffset()

	// 将值添加到当前插槽。
	// Add the value to the current slot.
	bucket := w.ring.At(w.offset % w.size).(*Bucket)
	bucket.Add(value)

	// 添加成功。
	// Add success.
	return nil
}

// calculateStats 计算滚动窗口中的值的总和和数量。
// calculateStats calculates the sum and count of the values in the rolling window.
func (w *RollingWindow) calculateStats() (float64, uint64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// 如果滚动窗口没有运行，返回一个错误。
	// If the rolling window is not running, return an error.
	if !w.runing {
		return 0, 0, com.ErrorRollingWindowStopped
	}

	// 更新滚动窗口。
	// Update the rolling window.
	w.updateOffset()

	// 计算滚动窗口中的值的总和和数量。
	// Calculate the sum and count of the values in the rolling window.
	var sum float64
	var count uint64
	for i := 0; i < w.size; i++ {
		bucket := w.ring.At(i).(*Bucket)
		sum += bucket.Sum()
		count += bucket.Count()
	}

	return sum, count, nil
}

// Avg 返回滚动窗口中的值的平均值。
// Avg returns the average of the values in the rolling window.
func (w *RollingWindow) Avg() (float64, uint64, error) {
	// 计算滚动窗口中的值的总和和数量。
	// Calculate the sum and count of the values in the rolling window.
	sum, count, err := w.calculateStats()
	if err != nil {
		return 0, 0, err
	}

	// 返回滚动窗口中的值的平均值。
	// Return the average of the values in the rolling window.
	return sum / float64(count), count, nil
}

// Sum 返回滚动窗口中的值的总和。
// Sum returns the sum of the values in the rolling window.
func (w *RollingWindow) Sum() (float64, uint64, error) {
	sum, count, err := w.calculateStats()
	if err != nil {
		return 0, 0, err
	}

	// 返回滚动窗口中的值的总和。
	// Return the sum of the values in the rolling window.
	return sum, count, nil
}
