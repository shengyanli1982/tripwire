package breaker

const (
	// DefaultKValue is the default value of k.
	DefaultKValue = 1.5

	// DefaultProtected is the default value of protected.
	DefaultProtected = 5
)

// Config is the configuration for the breaker.
type Config struct {
	k         float64
	protected int
}

// NewConfig returns a new configuration for the breaker.
func NewConfig() *Config {
	return &Config{
		k:         DefaultKValue,
		protected: DefaultProtected,
	}
}

// DefaultConfig returns the default configuration for the breaker.
func DefaultConfig() *Config {
	return NewConfig()
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

// isConfigValid checks if the configuration is valid.
func isConfigValid(conf *Config) *Config {
	if conf != nil {
		if conf.k < 1 || conf.k >= 5 {
			conf.k = DefaultKValue
		}
		if conf.protected < 0 {
			conf.protected = DefaultProtected
		}
	} else {
		conf = DefaultConfig()
	}

	return conf
}
