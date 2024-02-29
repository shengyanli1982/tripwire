package breaker

import (
	"errors"
	"math"
	"sync"

	rw "github.com/shengyanli1982/tripwire/internal/rolling"
	"github.com/shengyanli1982/tripwire/internal/utils"
)

const (
	// The default floating-point precision is set to 2.
	DefaultFloatingPrecision = 2
)

var ErrorServiceUnavailable = errors.New("service unavailable")

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
	return b.rwin.Sum()
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

// Reject rejects the execution.
func (b *Breaker) Reject(reason error) {
	b.config.callback.OnReject(b.rwin.Add(0), reason)
}

// Accept accepts the execution.
func (b *Breaker) Accept() {
	b.config.callback.OnAccept(b.rwin.Add(1))
}

// Allow checks if the circuit breaker allows the execution.
func (b *Breaker) Allow() (Result, error) {
	// Accept the execution.
	if err := b.accept(); err != nil {
		return nil, err
	}

	// Return the result notifier.
	return b, nil
}

// Do executes the given function with circuit breaker protection.
func (b *Breaker) Do(fn HandleFunc, fallback FallbackFunc, acceptable AcceptableFunc) error {
	var err error

	// If accept returns an error, reject the execution and return the error.
	if err = b.accept(); err != nil {
		b.Reject(err)
		if fallback != nil {
			return fallback(err)
		}
		return err
	}

	// Exec the handle function, if the error is acceptable, accept the execution, otherwise reject the execution.
	err = fn()
	if acceptable(err) {
		b.Accept()
	} else {
		b.Reject(err)
	}

	// Return the error.
	return err
}
