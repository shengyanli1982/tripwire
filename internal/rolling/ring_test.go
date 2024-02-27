package rolling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRing_Push(t *testing.T) {
	r := NewRing(3)
	assert.Equal(t, 0, r.Len(), "Len() = %d, expected %d", r.Len(), 0)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)

	// first loop

	r.Push(1)
	assert.Equal(t, 1, r.Len(), "Len() = %d, expected %d", r.Len(), 1)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 1, r.Head(), "Head() = %d, expected %d", r.Head(), 1)
	assert.Equal(t, 1, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 1)

	r.Push(2)
	assert.Equal(t, 2, r.Len(), "Len() = %d, expected %d", r.Len(), 2)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 1, r.Head(), "Head() = %d, expected %d", r.Head(), 1)
	assert.Equal(t, 2, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 2)

	r.Push(3)
	assert.Equal(t, 3, r.Len(), "Len() = %d, expected %d", r.Len(), 3)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 1, r.Head(), "Head() = %d, expected %d", r.Head(), 1)
	assert.Equal(t, 3, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 3)

	r.Push(4)
	assert.Equal(t, 4, r.Len(), "Len() = %d, expected %d", r.Len(), 4)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 1, r.Head(), "Head() = %d, expected %d", r.Head(), 1)
	assert.Equal(t, 4, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 4)

	// second loop

	r.Push(5)
	assert.Equal(t, 4, r.Len(), "Len() = %d, expected %d", r.Len(), 4)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 2, r.Head(), "Head() = %d, expected %d", r.Head(), 2)
	assert.Equal(t, 5, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 5)

	r.Push(6)
	assert.Equal(t, 4, r.Len(), "Len() = %d, expected %d", r.Len(), 4)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 3, r.Head(), "Head() = %d, expected %d", r.Head(), 3)
	assert.Equal(t, 6, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 6)

	// print buffer

	for i, v := range r.buffer {
		t.Logf("r.buffer[%d] = %v", i, v)
	}
}
func TestRing_Reset(t *testing.T) {
	r := NewRing(3)

	r.Push(1)
	r.Push(2)
	r.Push(3)
	r.Push(4)

	assert.Equal(t, 4, r.Len(), "Len() = %d, expected %d", r.Len(), 4)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 1, r.Head(), "Head() = %d, expected %d", r.Head(), 1)
	assert.Equal(t, 4, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 4)

	r.Reset()

	assert.Equal(t, 0, r.Len(), "Len() = %d, expected %d", r.Len(), 0)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, nil, r.Head(), "Head() = %d, expected %d", r.Head(), nil)
	assert.Equal(t, nil, r.Tail(), "Tail() = %d, expected %d", r.Tail(), nil)

	r.Push(7)
	r.Push(8)
	r.Push(9)
	r.Push(10)

	assert.Equal(t, 4, r.Len(), "Len() = %d, expected %d", r.Len(), 4)
	assert.Equal(t, 4, r.Cap(), "Cap() = %d, expected %d", r.Cap(), 4)
	assert.Equal(t, 7, r.Head(), "Head() = %d, expected %d", r.Head(), 7)
	assert.Equal(t, 10, r.Tail(), "Tail() = %d, expected %d", r.Tail(), 10)
}
