package rolling

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRollingWindow_Span(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// timeSleep
	timeSleep := 1

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 0; i < rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
	}

	// Call the span method.
	span := rw.span()

	// Check if the span is correct.
	assert.Equal(t, 0, span, "Span mismatch")

	// Long time waiting for the next test
	time.Sleep(time.Duration(timeSleep) * time.Second)

	// Call the span method.
	span = rw.span()

	// Check if the span is correct.
	assert.Equal(t, timeSleep*int(time.Second/DefaultRollingWindowSlotInterval), span, "Span mismatch")
}

func TestRollingWindow_UpdateSlots(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 0; i < rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
	}

	// Call the updateSlots method.
	rw.updateOffset()

	// Check if the slots have been reset.
	for i := 0; i < 5; i++ {
		bucket := rw.ring.At(i + 1).(*Bucket)
		assert.Equal(t, uint64(0), bucket.Count(), "Bucket count mismatch")
	}
}

func TestRollingWindow_Sum(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 1; i <= rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
	}

	// Calculate the expected sum.
	expectedSum := 1.0 + 2.0 + 3.0 + 4.0 + 5.0

	// Call the Sum method.
	sum, err := rw.Sum()
	assert.NoError(t, err, "Unexpected error")

	// Print the slots.
	for i := 0; i < len(rw.ring.slots); i++ {
		fmt.Printf("Slot %d: %v\n", i, rw.ring.slots[i])
	}

	// Check if the sum matches the expected sum.
	assert.Equal(t, expectedSum, sum, "Sum mismatch")
}

func TestRollingWindow_SumWithIdleSleep(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// timeSleep
	timeSleep := 500

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 1; i <= rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
		time.Sleep(time.Duration(timeSleep) * time.Millisecond)
	}

	// Calculate the expected sum.
	expectedSum := 1.0 + 2.0 + 3.0 + 4.0 + 5.0

	// Call the Sum method.
	sum, err := rw.Sum()
	assert.NoError(t, err, "Unexpected error")

	// Print the slots.
	for i := 0; i < len(rw.ring.slots); i++ {
		fmt.Printf("Slot %d: %v\n", i, rw.ring.slots[i])
	}

	// Check if the sum matches the expected sum.
	assert.Equal(t, expectedSum, sum, "Sum mismatch")
}

func TestRollingWindow_Avg(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 1; i <= rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
	}

	// Call the Avg method.
	avg, err := rw.Avg()
	assert.NoError(t, err, "Unexpected error")

	// Print the slots.
	for i := 0; i < len(rw.ring.slots); i++ {
		fmt.Printf("Slot %d: %v\n", i, rw.ring.slots[i])
	}

	// Calculate the expected average.
	expectedAvg := (1.0 + 2.0 + 3.0 + 4.0 + 5.0) / float64(rwSize)

	// Check if the average matches the expected average.
	assert.Equal(t, expectedAvg, avg, "Average mismatch")
}

func TestRollingWindow_AvgWithIdleSleep(t *testing.T) {
	// rolling window size.
	rwSize := 5

	// timeSleep
	timeSleep := 500

	// Create a new rolling rw.
	rw := NewRollingWindow(rwSize)
	defer rw.Stop()

	// Add some values to the rolling window.
	var err error
	for i := 1; i <= rwSize; i++ {
		err = rw.Add(float64(i))
		assert.NoError(t, err, "Unexpected error")
		time.Sleep(time.Duration(timeSleep) * time.Millisecond)
	}

	// Call the Avg method.
	avg, err := rw.Avg()
	assert.NoError(t, err, "Unexpected error")

	// Print the slots.
	for i := 0; i < len(rw.ring.slots); i++ {
		fmt.Printf("Slot %d: %v\n", i, rw.ring.slots[i])
	}

	// Calculate the expected average.
	expectedAvg := (1.0 + 2.0 + 3.0 + 4.0 + 5.0) / float64(rwSize)

	// Check if the average matches the expected average.
	assert.Equal(t, expectedAvg, avg, "Average mismatch")
}
