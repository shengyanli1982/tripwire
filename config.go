package tripwire

import (
	cb "github.com/shengyanli1982/tripwire/circuitbreaker"
	com "github.com/shengyanli1982/tripwire/common"
)

// 定义配置结构体
// Define the Config struct
type Config struct {
	retry   com.Retry   // 重试策略 Retry strategy
	breaker com.Breaker // 断路器 Circuit breaker
}

// NewConfig 函数创建并返回一个新的配置实例
// The NewConfig function creates and returns a new instance of Config
func NewConfig() *Config {
	return &Config{
		// 使用默认的重试策略
		// Use the default retry strategy
		retry: NewEmptyRetry(),

		// 使用默认的断路器
		// Use the default circuit breaker
		breaker: cb.NewGoogleBreaker(cb.DefaultConfig()),
	}
}

// WithRetry 方法设置配置的重试策略，并返回配置本身
// The WithRetry method sets the retry strategy of the Config and returns the Config itself
func (c *Config) WithRetry(retry com.Retry) *Config {
	c.retry = retry
	return c
}

// WithBreaker 方法设置配置的断路器，并返回配置本身
// The WithBreaker method sets the circuit breaker of the Config and returns the Config itself
func (c *Config) WithBreaker(breaker com.Breaker) *Config {
	c.breaker = breaker
	return c
}

// isConfigValid 函数检查配置是否有效，如果无效则使用默认值，最后返回配置
// The isConfigValid function checks whether the Config is valid, uses default values if invalid, and finally returns the Config
func isConfigValid(conf *Config) *Config {
	if conf != nil {
		// 如果重试策略为空，则使用默认的重试策略
		// If the retry strategy is nil, use the default retry strategy
		if conf.retry == nil {
			conf.retry = NewEmptyRetry()
		}

		// 如果断路器为空，则使用默认的断路器
		// If the circuit breaker is nil, use the default circuit breaker
		if conf.breaker == nil {
			conf.breaker = cb.NewGoogleBreaker(cb.DefaultConfig())
		}
	} else {
		// 如果配置为空，则创建一个新的配置
		// If the Config is nil, create a new Config
		conf = NewConfig()
	}

	return conf
}
