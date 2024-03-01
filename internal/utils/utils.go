package utils

import (
	"math"
	"math/rand"
	"time"
)

// FindNextPowerOfTwo 函数找到大于或等于 n 的最小的 2 的幂
// The FindNextPowerOfTwo function finds the smallest power of two that is greater than or equal to n
func FindNextPowerOfTwo(n int) int {
	// 如果 n 是 2 的幂，直接返回 n
	// If n is a power of 2, return n directly
	if (n & (n - 1)) == 0 {
		return n
	}

	// 从最高位开始，每次右移一位，直到 n 的所有位都是 1
	// Starting from the highest bit, shift right one bit at a time until all bits of n are 1
	for i := 1; i <= 32; i <<= 1 {
		n |= n >> i
	}

	// n + 1 是 2 的幂
	// n + 1 is a power of 2
	return n + 1
}

// Round 函数将浮点数 f 四舍五入到小数点后 n 位
// The Round function rounds the float number f to n decimal places
func Round(f float64, n int) float64 {
	// pow 是 10 的 n 次方
	// pow is 10 to the power of n
	pow := math.Pow(10, float64(n))

	// 返回四舍五入后的结果
	// Return the result after rounding
	return math.Round(f*pow) / pow
}

// 创建一个新的随机数生成器
// Create a new random number generator
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateRandomRatio 函数生成一个 [0,1) 之间的随机浮点数
// The GenerateRandomRatio function generates a random float number between [0,1)
func GenerateRandomRatio() float64 {
	return random.Float64()
}
