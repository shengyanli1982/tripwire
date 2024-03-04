package circuitbreaker

// 定义默认的常量值
// Define the default constant values
const (
	// DefaultKValue 是 k 的默认值。
	// DefaultKValue is the default value of k.
	DefaultKValue = 1.5

	// DefaultProtected 是 protected 的默认值。
	// DefaultProtected is the default value of protected.
	DefaultProtected = 5

	// DefaultStateWindow 是 state window 的默认值。
	// DefaultStateWindow is the default value of state window.
	DefaultStateWindow = 10
)

// Config 是熔断器的配置。
// Config is the configuration for the breaker.
type Config struct {
	k           float64
	protected   int
	callback    Callback
	stateWindow int
}

// NewConfig 返回熔断器的新配置。
// NewConfig returns a new configuration for the breaker.
func NewConfig() *Config {
	return &Config{
		k:           DefaultKValue,
		protected:   DefaultProtected,
		callback:    NewEmptyCallback(),
		stateWindow: DefaultStateWindow,
	}
}

// DefaultConfig 返回熔断器的默认配置。
// DefaultConfig returns the default configuration for the breaker.
func DefaultConfig() *Config {
	return NewConfig()
}

// WithCallback 设置配置的回调函数。
// WithCallback sets the callback of the configuration.
func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

// WithK 设置配置的 k 值。
// WithK sets the k value of the configuration.
func (c *Config) WithK(k float64) *Config {
	c.k = k
	return c
}

// WithProtected 设置配置的 protected 值。
// WithProtected sets the protected value of the configuration.
func (c *Config) WithProtected(protected int) *Config {
	c.protected = protected
	return c
}

// WithStateWindow 设置配置的 state window 值。
// WithStateWindow sets the state window of the configuration.
func (c *Config) WithStateWindow(window int) *Config {
	c.protected = window
	return c
}

// isConfigValid 检查配置是否有效。
// isConfigValid checks if the configuration is valid.
func isConfigValid(conf *Config) *Config {
	if conf != nil {
		if conf.k < 1 || conf.k >= 5 {
			conf.k = DefaultKValue
		}
		if conf.protected < 0 {
			conf.protected = DefaultProtected
		}
		if conf.callback == nil {
			conf.callback = NewEmptyCallback()
		}
		if conf.stateWindow <= 0 {
			conf.stateWindow = DefaultStateWindow
		}
	} else {
		conf = DefaultConfig()
	}

	return conf
}
