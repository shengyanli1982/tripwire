package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtils_FindNextPowerOfTwo(t *testing.T) {
	tests := []struct {
		n        int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{10, 16},
		{16, 16},
		{17, 32},
		{100, 128},
	}

	for _, test := range tests {
		result := FindNextPowerOfTwo(test.n)
		assert.Equal(t, test.expected, result, "FindNextPowerOfTwo(%d) = %d, expected %d", test.n, result, test.expected)
	}
}

func TestUtils_Round(t *testing.T) {
	tests := []struct {
		f        float64
		n        int
		expected float64
	}{
		{0.12345, 2, 0.12},
		{0.6789, 3, 0.679},
		{1.23456789, 4, 1.2346},
		{3.14159, 0, 3},
		{5.678, 1, 5.7},
	}

	for _, test := range tests {
		result := Round(test.f, test.n)
		assert.Equal(t, test.expected, result, "Round(%f, %d) = %f, expected %f", test.f, test.n, result, test.expected)
	}
}
