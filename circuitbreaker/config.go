package circuitbreaker

const (
	// DefaultKValue is the default value of k.
	DefaultKValue = 1.5

	// DefaultProtected is the default value of protected.
	DefaultProtected = 5

	// DefaultStateWindow is the default value of state window.
	DefaultStateWindow = 10
)

// Config is the configuration for the breaker.
type Config struct {
	k           float64
	protected   int
	callback    Callback
	stateWindow int
}

// NewConfig returns a new configuration for the breaker.
func NewConfig() *Config {
	return &Config{
		k:           DefaultKValue,
		protected:   DefaultProtected,
		callback:    NewEmptyCallback(),
		stateWindow: DefaultStateWindow,
	}
}

// DefaultConfig returns the default configuration for the breaker.
func DefaultConfig() *Config {
	return NewConfig()
}

func (c *Config) WithCallback(callback Callback) *Config {
	c.callback = callback
	return c
}

// WithK sets the k value of the configuration.
func (c *Config) WithK(k float64) *Config {
	c.k = k
	return c
}

// WithProtected sets the protected value of the configuration.
func (c *Config) WithProtected(protected int) *Config {
	c.protected = protected
	return c
}

// WithStateWindow sets the state window of the configuration.
func (c *Config) WithStateWindow(window int) *Config {
	c.protected = window
	return c
}

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
