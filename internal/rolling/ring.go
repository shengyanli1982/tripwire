package rolling

type Ring struct {
	// The ring buffer.
	buffer []any
	// The index of the first element in the ring buffer.
	head uint64
	// The index of the last element in the ring buffer.
	tail uint64
	// The number of elements in the ring buffer.
	count int
	// The maximum number of elements in the ring buffer. This is must be a power of 2.
	cap int
	// The default value for new elements in the ring buffer.
	optval uint64
}

// NewRing returns a new ring buffer with the specified capacity.
func NewRing(cap int) *Ring {
	// Ensure that the capacity is a power of 2.
	cap = FindNextPowerOfTwo(cap)

	// Create and return the ring buffer.
	return &Ring{
		buffer: make([]any, cap), // The ring buffer is a slice of any.
		cap:    cap,              // cap is a power of 2.
		optval: uint64(cap - 1),  // optval is cap - 1. This is used to calculate the modulo of the indices.
	}
}

// Reset resets the ring buffer.
func (r *Ring) Reset() {
	r.head = 0
	r.tail = 0
	r.count = 0
}

// Len returns the number of elements in the ring buffer.
func (r *Ring) Len() int {
	// If the ring buffer is not full, return the number of elements in the ring buffer.
	if r.count < r.cap {
		return r.count
	}

	// If the ring buffer is full, return the maximum number of elements in the ring buffer.
	return r.cap
}

// Cap returns the maximum number of elements in the ring buffer.
func (r *Ring) Cap() int {
	return r.cap
}

// Head returns the first element in the ring buffer.
func (r *Ring) Head() any {
	// If the ring buffer is empty, return nil.
	if r.count == 0 {
		return nil
	}
	// Return the first element in the ring buffer.
	return r.buffer[r.head]
}

// Tail returns the last element in the ring buffer.
func (r *Ring) Tail() any {
	// If the ring buffer is empty, return nil.
	if r.count == 0 {
		return nil
	}
	// Return the last element in the ring buffer.
	return r.buffer[(r.tail-1)&r.optval]
}

// Push appends a new element to the ring buffer.
func (r *Ring) Push(v any) {
	// Append the new element to the ring buffer.
	r.buffer[r.tail] = v

	// Increment the tail index and the count.
	r.tail = (r.tail + 1) & r.optval

	// Increment the count.
	if r.count < r.cap {
		r.count++
	} else {
		// Increment the head index. This overwrites the oldest element in the ring buffer. Let r.head & r.optval == 0
		r.head = (r.head + 1) & r.optval
	}
}
