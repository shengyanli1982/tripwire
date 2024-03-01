package circuitbreaker

import (
	"math"
	"sync"

	com "github.com/shengyanli1982/tripwire/common"
	rw "github.com/shengyanli1982/tripwire/internal/rolling"
)

const (
	// 定义默认的浮点数精度为 3
	// Define the default floating-point precision as 3
	DefaultFloatingPrecision = 3
)

// DefaultAcceptableFunc 是默认的可接受函数。
// DefaultAcceptableFunc is the default acceptable function.
func DefaultAcceptableFunc(err error) bool { return err == nil } // 如果错误为nil，则返回true

// DefaultFallbackFunc 是默认的回退函数。
// DefaultFallbackFunc is the default fallback function.
func DefaultFallbackFunc(err error) error { return err } // 直接返回错误

// GoogleBreaker 是一个当错误率高时打开的熔断器。
// GoogleBreaker is a circuit breaker that opens when the error rate is high.
type GoogleBreaker struct {
	config *Config           // 熔断器的配置 Config of the breaker
	rwin   *rw.RollingWindow // 滚动窗口 Rolling window
	once   sync.Once         // 用于确保某个操作只执行一次 The sync.Once to ensure that an operation is executed only once
	sr     *SafeRandom       // 安全的随机数生成器 Safe random number generator
}

// NewGoogleBreaker 返回一个新的熔断器。
// NewGoogleBreaker returns a new breaker.
func NewGoogleBreaker(conf *Config) *GoogleBreaker {
	conf = isConfigValid(conf)
	return &GoogleBreaker{
		config: conf,
		sr:     NewSafeRandom(),
		rwin:   rw.NewRollingWindow(conf.stateWindow),
		once:   sync.Once{},
	}
}

// Stop 停止熔断器。
// Stop stops the breaker.
func (b *GoogleBreaker) Stop() {
	b.once.Do(func() {
		b.rwin.Stop() // 停止滚动窗口
	})
}

// history 返回熔断器的历史。接受和总计的和，以及任何错误
// history returns the history of the breaker. Sum of accepted and total, and error if any
func (b *GoogleBreaker) history() (float64, uint64, error) {
	return b.rwin.Sum() // 返回滚动窗口的和
}

// Accept 接受一个执行。
// Accept accepts a execution.
func (b *GoogleBreaker) accept(ratio float64) error {
	// 获取熔断器的历史状态。
	// Get the history state of the breaker.
	accepted, total, err := b.history()
	if err != nil {
		return err
	}

	// 计算加权接受。
	// Calculate the weighted accepts.
	weightedAcceptes := b.config.k * accepted

	// 计算熔丝比率。
	// Calculate the fuse ratio.
	refFactor := (float64(int64(total)-int64(b.config.protected)) - weightedAcceptes) / float64(total+1)
	fuseRatio := math.Max(0, refFactor)

	// 如果熔丝比率小于或等于0，或者熔丝比率大于等于0和1之间的随机浮点数，返回nil。
	// If the fuse ratio is less than or equal to 0, or if the fuse ratio is greater than or equal a random float64 between 0 and 1, return nil.
	if fuseRatio <= 0 || ratio >= fuseRatio {
		b.config.callback.OnAccept(nil, refFactor)
		return nil
	}

	// 如果熔丝比率大于随机浮点数，返回服务不可用的错误。
	// If the fuse ratio is greater than the random float64, return the error of service unavailable.
	b.config.callback.OnAccept(com.ErrorServiceUnavailable, refFactor)
	return com.ErrorServiceUnavailable
}

// MarkFailure 标记一个失败的执行，并调用失败回调
// MarkFailure marks a failed execution and calls the failure callback
func (b *GoogleBreaker) MarkFailure(reason error) {
	b.config.callback.OnFailure(b.rwin.Add(0), reason) // 添加一个失败的执行，并调用失败回调
	// Add a failed execution and call the failure callback
}

// MarkSuccess 标记一个成功的执行，并调用成功回调
// MarkSuccess marks a successful execution and calls the success callback
func (b *GoogleBreaker) MarkSuccess() {
	b.config.callback.OnSuccess(b.rwin.Add(1)) // 添加一个成功的执行，并调用成功回调
	// Add a successful execution and call the success callback
}

// Allow 检查熔断器是否允许执行。
// Allow checks if the circuit breaker allows the execution.
func (b *GoogleBreaker) Allow() (com.Notifier, error) {
	// 接受执行。
	// Accept the execution.
	if err := b.accept(b.sr.Float64()); err != nil {
		return nil, err
	}

	// 返回结果通知器。
	// Return the result notifier.
	return b, nil
}

// do 使用熔断器保护执行给定的函数。
// do executes the given function with circuit breaker protection.
func (b *GoogleBreaker) do(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	var err error

	// 如果 accept 返回错误，拒绝执行并返回错误。
	// If accept returns an error, reject the execution and return the error.
	if err = b.accept(b.sr.Float64()); err != nil {
		// 标记执行失败
		// Mark the execution as failed
		b.MarkFailure(err)

		// 如果提供了回退函数，执行回退函数。
		// If a fallback function is provided, execute the fallback function.
		if fallback != nil {
			return fallback(err)
		}

		// 返回错误。
		// Return the error.
		return err
	}

	// 执行函数
	// Execute the function
	err = fn()

	// 如果错误可接受，标记执行成功，否则标记执行失败并返回错误。
	// If the error is acceptable, mark the execution as successful, otherwise mark the execution as failed and return the error.
	if acceptable(err) {
		// 标记执行成功
		// Mark the execution as successful
		b.MarkSuccess()

		// 正常返回
		// Return nil
		return nil
	} else {
		// 标记执行失败
		// Mark the execution as failed
		b.MarkFailure(err)

		// 返回错误。
		// Return the error.
		return err
	}
}

// Do 执行函数并返回错误。
// Do executes the function and returns the error.
func (b *GoogleBreaker) Do(fn com.HandleFunc) error {
	return b.do(fn, nil, DefaultAcceptableFunc)
}

// DoWithAcceptable 使用给定的可接受函数执行函数并返回错误。
// DoWithAcceptable executes the function with the given acceptable function and returns the error.
func (b *GoogleBreaker) DoWithAcceptable(fn com.HandleFunc, acceptable com.AcceptableFunc) error {
	return b.do(fn, nil, acceptable)
}

// DoWithFallback 使用给定的回退函数执行函数并返回错误。
// DoWithFallback executes the function with the given fallback function and returns the error.
func (b *GoogleBreaker) DoWithFallback(fn com.HandleFunc, fallback com.FallbackFunc) error {
	return b.do(fn, fallback, DefaultAcceptableFunc)
}

// DoWithFallbackAcceptable 使用给定的回退和可接受函数执行函数并返回错误。
// DoWithFallbackAcceptable executes the function with the given fallback and acceptable functions and returns the error.
func (b *GoogleBreaker) DoWithFallbackAcceptable(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	return b.do(fn, fallback, acceptable)
}
