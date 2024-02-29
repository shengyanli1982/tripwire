package breaker

// See Client-Side Throttling section in
// https://landing.google.com/sre/sre-book/chapters/handling-overload/

import (
	"errors"
	"math"
	"sync"

	rw "github.com/shengyanli1982/tripwire/internal/rolling"
	"github.com/shengyanli1982/tripwire/internal/utils"
)

// The default floating-point precision is set to 2.
const DefaultFloatingPrecision = 2

var (
	ErrorServiceUnavailable = errors.New("service unavailable")
)

// Breaker is a circuit breaker that opens when the error rate is high.
type Breaker struct {
	config *Config
	rwin   *rw.RollingWindow
	once   sync.Once
}

// NewBreaker returns a new breaker.
func NewBreaker(conf *Config) *Breaker {
	conf = isConfigValid(conf)
	return &Breaker{
		config: conf,
		rwin:   rw.NewRollingWindow(conf.protected),
		once:   sync.Once{},
	}
}

// Stop stops the breaker.
func (b *Breaker) Stop() {
	b.once.Do(func() {
		b.rwin.Stop()
	})
}

// history returns the history of the breaker. Sum of accepted and total, and error if any
func (b *Breaker) history() (float64, uint64, error) {
	accepted, total, err := b.rwin.Sum()
	return accepted, total, err
}

// Accept accepts a execution.
func (b *Breaker) accept() error {
	// Get the history state of the breaker.
	accepted, total, err := b.history()
	if err != nil {
		return err
	}

	// Calculate the weighted accepts.
	weightedAccepted := b.config.k * accepted

	// Calculate the fuse ratio.
	fuseRatio := math.Max(0, (float64(total-uint64(b.config.protected))-weightedAccepted)/float64(total+1))

	// If the fuse ratio is less than or equal to 0, or if the fuse ratio is less than a random float64 between 0 and 1, return nil.
	if fuseRatio <= 0 || fuseRatio < utils.Round(utils.GenerateRandomRatio(), DefaultFloatingPrecision) {
		return nil
	}

	// Otherwise, trigger the breaker.
	return ErrorServiceUnavailable
}
