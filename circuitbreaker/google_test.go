package circuitbreaker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleBreaker_Accept(t *testing.T) {
	// Test case 1: Fuse ratio <= random ratio
	config := NewConfig().WithK(0.2).WithProtected(5).WithStateWindow(10)
	breaker := NewGoogleBreaker(config)
	for i := 0; i < 20; i++ {
		_ = breaker.rwin.Add(1)
	}
	for i := 0; i < 2; i++ {
		_ = breaker.rwin.Add(0)
	}
	err := breaker.accept(0.3)
	assert.NoError(t, err, "Expected no error")

	// Test case 2: Fuse ratio > random ratio
	config = NewConfig().WithK(2).WithProtected(5).WithStateWindow(10)
	breaker = NewGoogleBreaker(config)
	for i := 0; i < 2; i++ {
		_ = breaker.rwin.Add(1)
	}
	for i := 0; i < 20; i++ {
		_ = breaker.rwin.Add(0)
	}
	err = breaker.accept(0.5)
	assert.ErrorIs(t, err, ErrorServiceUnavailable, "Expected error")
}

func TestGoogleBreaker_Allow(t *testing.T) {
	breaker := NewGoogleBreaker(nil)

	// Test allowing execution
	notifier, err := breaker.Allow()
	assert.NoError(t, err, "Unexpected error")

	// Verify the returned notifier
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// Test rejecting execution
	breaker.Reject(errors.New("test"))

	// Test values execution
	v, c, _ := breaker.history()
	assert.Equal(t, float64(0), v, "Expected 1, but got %v", v)
	assert.Equal(t, uint64(1), c, "Expected 1, but got %v", c)

	// Test allowing execution
	notifier, err = breaker.Allow()
	assert.NoError(t, err, "Unexpected error")

	// Verify the returned notifier
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// Test accepting execution
	breaker.Accept()

	// Test values execution
	v, c, _ = breaker.history()
	assert.Equal(t, float64(1), v, "Expected 1, but got %v", v)
	assert.Equal(t, uint64(2), c, "Expected 1, but got %v", c)
}

func TestGoogleBreaker_DoWithFallbackAcceptable(t *testing.T) {
	var (
		execError = errors.New("execution error")
		// fbError   = errors.New("fallback error")
	)

	breaker := NewGoogleBreaker(nil)

	// Test case 1: Successful execution with acceptable result
	fn := func() error {
		return nil
	}
	fallback := func(err error) error {
		return err
	}
	acceptable := func(err error) bool {
		return err == nil
	}
	err := breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.NoError(t, err, "Unexpected error")

	// Test case 2: Successful execution with unacceptable result
	fn = func() error {
		return nil
	}
	acceptable = func(err error) bool {
		return err != nil
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.NoError(t, err, "Unexpected error")

	// Test case 3: Failed execution with acceptable result
	fn = func() error {
		return execError
	}
	acceptable = func(err error) bool {
		return err == execError
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, execError, "Unexpected error")

	// Test case 4: Failed execution with unacceptable result
	fn = func() error {
		return execError
	}
	acceptable = func(err error) bool {
		return err != execError
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, execError, "Unexpected error")
}

func TestGoogleBreaker_DoWithFallback(t *testing.T) {
	breaker := NewGoogleBreaker(nil)

	// Test case 1: Successful execution
	fn := func() error {
		return nil
	}
	fallback := func(err error) error {
		return err
	}
	err := breaker.DoWithFallback(fn, fallback)
	assert.NoError(t, err, "Unexpected error")

	// Test case 2: Failed execution
	execError := errors.New("execution error")
	fn = func() error {
		return execError
	}
	err = breaker.DoWithFallback(fn, fallback)
	assert.ErrorIs(t, err, execError, "Unexpected error")
}

func TestGoogleBreaker_DoWithAcceptable(t *testing.T) {
	breaker := NewGoogleBreaker(nil)

	// Test case 1: Successful execution
	fn := func() error {
		return nil
	}
	acceptable := func(err error) bool {
		return err == nil
	}
	err := breaker.DoWithAcceptable(fn, acceptable)
	assert.NoError(t, err, "Unexpected error")

	// Test case 2: Failed execution
	execError := errors.New("execution error")
	fn = func() error {
		return execError
	}
	acceptable = func(err error) bool {
		return err != execError
	}
	err = breaker.DoWithAcceptable(fn, acceptable)
	assert.ErrorIs(t, err, execError, "Unexpected error")
}

func TestGoogleBreaker_Do(t *testing.T) {
	breaker := NewGoogleBreaker(nil)

	// Test case 1: Successful execution
	fn := func() error {
		return nil
	}
	err := breaker.Do(fn)
	assert.NoError(t, err, "Unexpected error")

	// Test case 2: Failed execution
	execError := errors.New("execution error")
	fn = func() error {
		return execError
	}
	err = breaker.Do(fn)
	assert.ErrorIs(t, err, execError, "Unexpected error")
}
