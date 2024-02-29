package utils

import (
	"math"
	"math/rand"
	"time"
)

// FindNextPowerOfTwo returns the next power of 2 that is greater than or equal to n.
func FindNextPowerOfTwo(n int) int {
	// If n is already a power of 2, return n.
	if (n & (n - 1)) == 0 {
		return n
	}

	// Set the most significant bit to 1 and all other bits to 0.
	for i := 1; i <= 32; i <<= 1 {
		n |= n >> i
	}

	// Increment n by 1 to get the next power of 2.
	return n + 1
}

// Round returns f rounded to n decimal places.
func Round(f float64, n int) float64 {
	// Calculate 10^n.
	pow := math.Pow(10, float64(n))

	// Round f to n decimal places.
	return math.Round(f*pow) / pow
}

// random is a random number generator.
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomRatio returns a random float64 between 0 and 1.
func GenerateRandomRatio() float64 {
	return random.Float64()
}
