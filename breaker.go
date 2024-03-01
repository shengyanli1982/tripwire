package tripwire

import (
	"sync"

	cb "github.com/shengyanli1982/tripwire/circuitbreaker"
	com "github.com/shengyanli1982/tripwire/common"
)

// CircuitBreaker 结构体包含了熔断器和重试机制
// The CircuitBreaker struct contains the circuit breaker and retry mechanism
type CircuitBreaker struct {
	breaker com.Breaker
	retry   com.Retry
	once    sync.Once
}

// NewCircuitBreaker 创建一个新的熔断器，如果没有提供熔断器或重试机制，会使用默认的
// NewCircuitBreaker creates a new circuit breaker, if no breaker or retry mechanism is provided, the default ones will be used
func NewCircuitBreaker(breaker com.Breaker, retry com.Retry) *CircuitBreaker {
	if breaker == nil {
		breaker = cb.NewGoogleBreaker(nil)
	}
	if retry == nil {
		retry = NewEmptyRetry()
	}
	return &CircuitBreaker{
		breaker: breaker,
		retry:   retry,
		once:    sync.Once{},
	}
}

// Stop 停止熔断器的运行
// Stop stops the operation of the circuit breaker
func (c *CircuitBreaker) Stop() {
	c.once.Do(func() {
		c.breaker.Stop()
	})
}

// DoWithFallbackAcceptable 使用回退和可接受函数执行函数
// DoWithFallbackAcceptable executes the function with fallback and acceptable functions
func (c *CircuitBreaker) DoWithFallbackAcceptable(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	result := c.retry.TryOnConflict(func() (any, error) {
		return nil, c.breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	})
	return result.TryError()
}

// DoWithFallback 使用回退函数执行函数
// DoWithFallback executes the function with fallback function
func (c *CircuitBreaker) DoWithFallback(fn com.HandleFunc, fallback com.FallbackFunc) error {
	result := c.retry.TryOnConflict(func() (any, error) {
		return nil, c.breaker.DoWithFallback(fn, fallback)
	})
	return result.TryError()
}

// DoWithAcceptable 使用可接受函数执行函数
// DoWithAcceptable executes the function with acceptable function
func (c *CircuitBreaker) DoWithAcceptable(fn com.HandleFunc, acceptable com.AcceptableFunc) error {
	result := c.retry.TryOnConflict(func() (any, error) {
		return nil, c.breaker.DoWithAcceptable(fn, acceptable)
	})
	return result.TryError()
}

// Do 执行函数
// Do executes the function
func (c *CircuitBreaker) Do(fn com.HandleFunc) error {
	result := c.retry.TryOnConflict(func() (any, error) {
		return nil, c.breaker.Do(fn)
	})
	return result.TryError()
}
