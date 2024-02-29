package tripwire

import (
	cb "github.com/shengyanli1982/tripwire/circuitbreaker"
	com "github.com/shengyanli1982/tripwire/common"
)

type Config struct {
	retry   com.Retry
	breaker com.Breaker
}

func NewConfig() *Config {
	return &Config{
		retry:   NewEmptyRetry(),
		breaker: cb.NewGoogleBreaker(cb.DefaultConfig()),
	}
}

func (c *Config) WithRetry(retry com.Retry) *Config {
	c.retry = retry
	return c
}

func (c *Config) WithBreaker(breaker com.Breaker) *Config {
	c.breaker = breaker
	return c
}

func isConfigValid(conf *Config) *Config {
	if conf != nil {
		if conf.retry == nil {
			conf.retry = NewEmptyRetry()
		}
		if conf.breaker == nil {
			conf.breaker = cb.NewGoogleBreaker(cb.DefaultConfig())
		}
	} else {
		conf = NewConfig()
	}

	return conf
}
