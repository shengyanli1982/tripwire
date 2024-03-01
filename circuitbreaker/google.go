package circuitbreaker

import (
	"errors"
	"math"
	"sync"

	com "github.com/shengyanli1982/tripwire/common"
	rw "github.com/shengyanli1982/tripwire/internal/rolling"
	"github.com/shengyanli1982/tripwire/internal/utils"
)

const (
	// The default floating-point precision is set to 2.
	DefaultFloatingPrecision = 2
)

var ErrorServiceUnavailable = errors.New("service unavailable")

// DefaultAcceptableFunc is the default acceptable function.
func DefaultAcceptableFunc(err error) bool { return err == nil }

// DefaultFallbackFunc is the default fallback function.
func DefaultFallbackFunc(err error) error { return err }

// GoogleBreaker is a circuit breaker that opens when the error rate is high.
type GoogleBreaker struct {
	config *Config
	rwin   *rw.RollingWindow
	once   sync.Once
}

// NewGoogleBreaker returns a new breaker.
func NewGoogleBreaker(conf *Config) *GoogleBreaker {
	conf = isConfigValid(conf)
	return &GoogleBreaker{
		config: conf,
		rwin:   rw.NewRollingWindow(conf.stateWindow),
		once:   sync.Once{},
	}
}

// Stop stops the breaker.
func (b *GoogleBreaker) Stop() {
	b.once.Do(func() {
		b.rwin.Stop()
	})
}

// history returns the history of the breaker. Sum of accepted and total, and error if any
func (b *GoogleBreaker) history() (float64, uint64, error) {
	return b.rwin.Sum()
}

// Accept accepts a execution.
func (b *GoogleBreaker) accept(ratio float64) error {
	// Get the history state of the breaker.
	accepted, total, err := b.history()
	if err != nil {
		return err
	}

	// Calculate the weighted accepts.
	weightedAccepted := b.config.k * accepted

	// Calculate the fuse ratio.
	refFactor := utils.Round((float64(total-uint64(b.config.protected))-weightedAccepted)/float64(total+1), DefaultFloatingPrecision)
	fuseRatio := math.Max(0, refFactor)

	// If the fuse ratio is less than or equal to 0, or if the fuse ratio is less than a random float64 between 0 and 1, return nil.
	if fuseRatio <= 0 || fuseRatio >= utils.Round(ratio, DefaultFloatingPrecision) {
		b.config.callback.OnAccept(nil, refFactor)
		return nil
	}

	// Otherwise, trigger the breaker.
	b.config.callback.OnAccept(ErrorServiceUnavailable, refFactor)
	return ErrorServiceUnavailable
}

// Reject rejects the execution.
func (b *GoogleBreaker) Reject(reason error) {
	b.config.callback.OnFailed(b.rwin.Add(0), reason)
}

// Accept accepts the execution.
func (b *GoogleBreaker) Accept() {
	b.config.callback.OnSuccess(b.rwin.Add(1))
}

// Allow checks if the circuit breaker allows the execution.
func (b *GoogleBreaker) Allow() (com.Notifier, error) {
	// Accept the execution.
	if err := b.accept(utils.GenerateRandomRatio()); err != nil {
		return nil, err
	}

	// Return the result notifier.
	return b, nil
}

// do executes the given function with circuit breaker protection.
func (b *GoogleBreaker) do(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	var err error

	// If accept returns an error, reject the execution and return the error.
	if err = b.accept(utils.GenerateRandomRatio()); err != nil {
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

// Do executes the function and returns the error.
func (b *GoogleBreaker) Do(fn com.HandleFunc) error {
	return b.do(fn, nil, DefaultAcceptableFunc)
}

// DoWithAcceptable executes the function with the given acceptable function and returns the error.
func (b *GoogleBreaker) DoWithAcceptable(fn com.HandleFunc, acceptable com.AcceptableFunc) error {
	return b.do(fn, nil, acceptable)
}

// DoWithFallback executes the function with the given fallback function and returns the error.
func (b *GoogleBreaker) DoWithFallback(fn com.HandleFunc, fallback com.FallbackFunc) error {
	return b.do(fn, fallback, DefaultAcceptableFunc)
}

// DoWithFallbackAcceptable executes the function with the given fallback and acceptable functions and returns the error.
func (b *GoogleBreaker) DoWithFallbackAcceptable(fn com.HandleFunc, fallback com.FallbackFunc, acceptable com.AcceptableFunc) error {
	return b.do(fn, fallback, acceptable)
}
