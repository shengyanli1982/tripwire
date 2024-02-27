package rolling

import (
	"errors"
	"sync"
	"time"
)

const (
	// The default size of the rolling window. 1 slot equal 1s, 10 slots equal 10s.
	DefaultRollingWindowSize = 10

	// The default duration of each slot in the rolling window.
	DefaultRollingWindowSlotInterval = time.Millisecond * 200
)

var (
	minRollingWindowSize = 5       // 5 slots, 5s
	maxRollingWindowSize = 10 * 60 // 600 slots, 10 minutes
)

var (
	ErrorRollingWindowStopped = errors.New("rolling window stopped")
)

// RollingWindow is a rolling window.
type RollingWindow struct {
	// The ring buffer.
	ring *Ring

	// The size of the rolling window.
	size int

	// The duration of each slot.
	interval time.Duration

	// The time when the rolling window slot was last updated.
	updateAt time.Time

	// The mutex to protect the rolling window.
	lock sync.Mutex

	// The flag to indicate if the rolling window is running.
	runing bool

	// The sync.Once to ensure that the rolling window is stopped only once.
	once sync.Once
}

// NewRollingWindow returns a new rolling window with the specified size and slot duration.
func NewRollingWindow(size int) *RollingWindow {
	// If the size is less than the minimum size or greater than the maximum size, use the default size.
	if size < minRollingWindowSize || size > maxRollingWindowSize {
		size = DefaultRollingWindowSize
	}

	// Create and return the rolling window.
	rw := RollingWindow{
		ring:     NewRing(size),
		size:     size,
		interval: DefaultRollingWindowSlotInterval,
		// ignoreCurrent: ignore,
		lock:     sync.Mutex{},
		runing:   true,
		once:     sync.Once{},
		updateAt: time.Now(),
	}

	// Initialize the rolling window.
	for i := 0; i < rw.ring.Cap(); i++ {
		rw.ring.Push(NewBucket())
	}

	// Return the rolling window.
	return &rw
}

// Stop stops the rolling window.
func (w *RollingWindow) Stop() {
	w.once.Do(func() {
		w.lock.Lock()
		w.runing = false
		w.ring.Reset()
		w.lock.Unlock()
	})
}

// span returns the number of slots that have elapsed since the rolling window was last updated.
func (w *RollingWindow) span() int {
	return int(time.Since(w.updateAt) / w.interval)
}

// updateSlot updates the rolling window.
func (w *RollingWindow) updateSlots() {
	// Calculate the number of slots that have elapsed since the rolling window was last updated.
	n := w.span()

	// If the rolling window has not been updated, return.
	if n == 0 {
		return
	}

	// If the rolling window has been move forward, reset the slots that have elapsed.
	for i := 0; i < n; i++ {
		bucket := w.ring.At(i).(*Bucket)
		bucket.Reset()
	}

	// Update the time when the rolling window slot was last updated.
	w.updateAt = time.Now()
}

// Add adds a value to the rolling window.
func (w *RollingWindow) Add(value float64) error {
	w.lock.Lock()
	defer w.lock.Unlock()

	// If the rolling window is not running, return an error.
	if !w.runing {
		return ErrorRollingWindowStopped
	}

	// Update the rolling window.
	w.updateSlots()

	// Add the value to the current slot.
	bucket := w.ring.Head().(*Bucket)
	bucket.Add(value)

	return nil
}

// Avg returns the average of the values in the rolling window.
func (w *RollingWindow) Avg() (float64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// If the rolling window is not running, return an error.
	if !w.runing {
		return 0, ErrorRollingWindowStopped
	}

	// Update the rolling window.
	w.updateSlots()

	// Calculate the average of the values in the rolling window.
	var sum float64
	var count uint64
	for i := 0; i < w.ring.Len(); i++ {
		bucket := w.ring.At(i).(*Bucket)
		sum += bucket.Sum()
		count += bucket.Count()
	}

	// If the count is 0, return 0.
	if count == 0 {
		return 0, nil
	}

	return sum / float64(count), nil
}

// Sum returns the sum of the values in the rolling window.
func (w *RollingWindow) Sum() (float64, error) {
	w.lock.Lock()
	defer w.lock.Unlock()

	// If the rolling window is not running, return an error.
	if !w.runing {
		return 0, ErrorRollingWindowStopped
	}

	// Update the rolling window.
	w.updateSlots()

	// Calculate the sum of the values in the rolling window.
	var sum float64
	for i := 0; i < w.ring.Len(); i++ {
		bucket := w.ring.At(i).(*Bucket)
		sum += bucket.Sum()
	}

	return sum, nil
}
