package tripwire

import (
	"sync"

	com "github.com/shengyanli1982/tripwire/common"
)

// CircuitBreaker 结构体包含了熔断器和重试机制
// The CircuitBreaker struct contains the circuit breaker and retry mechanism
type CircuitBreaker struct {
	config *Config
	once   sync.Once
}

// New 创建一个新的熔断器，如果没有提供熔断器或重试机制，会使用默认的
// New creates a new circuit breaker, if no breaker or retry mechanism is provided, the default ones will be used
func New(config *Config) *CircuitBreaker {
	config = isConfigValid(config)
	return &CircuitBreaker{
		config: config,
		once:   sync.Once{},
	}
}

// Stop 停止熔断器的运行
// Stop stops the operation of the circuit breaker
func (c *CircuitBreaker) Stop() {
	c.once.Do(func() {
		c.config.breaker.Stop()
	})
}

// Allow 方法手动操作熔断器是否允许请求通过 (纯手动，不建议使用)
// The Allow method manually operates whether the circuit breaker allows requests to pass through (pure manual, not recommended)
func (c *CircuitBreaker) Allow() (com.Notifier, error) {
	// 调用配置中的熔断器的 Allow 方法
	// Call the Allow method of the breaker in the configuration
	return c.config.breaker.Allow()
}

// DoWithFallbackAcceptable 使用回退和可接受函数执行函数
// DoWithFallbackAcceptable executes the function with fallback and acceptable functions
func (c *CircuitBreaker) DoWithFallbackAcceptable(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	result := c.config.retry.TryOnConflictVal(func() (any, error) {
		return nil, c.config.breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	})
	return result.TryError()
}

// DoWithFallback 使用回退函数执行函数
// DoWithFallback executes the function with fallback function
func (c *CircuitBreaker) DoWithFallback(fn com.HandleFunc, fallback com.FallbackFunc) error {
	result := c.config.retry.TryOnConflictVal(func() (any, error) {
		return nil, c.config.breaker.DoWithFallback(fn, fallback)
	})
	return result.TryError()
}

// DoWithAcceptable 使用可接受函数执行函数
// DoWithAcceptable executes the function with acceptable function
func (c *CircuitBreaker) DoWithAcceptable(fn com.HandleFunc, acceptable com.AcceptableFunc) error {
	result := c.config.retry.TryOnConflictVal(func() (any, error) {
		return nil, c.config.breaker.DoWithAcceptable(fn, acceptable)
	})
	return result.TryError()
}

// Do 执行函数
// Do executes the function
func (c *CircuitBreaker) Do(fn com.HandleFunc) error {
	result := c.config.retry.TryOnConflictVal(func() (any, error) {
		return nil, c.config.breaker.Do(fn)
	})
	return result.TryError()
}
