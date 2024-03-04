package tripwire

import (
	"errors"
	"testing"

	com "github.com/shengyanli1982/tripwire/common"
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_DoWithFallbackAcceptable(t *testing.T) {
	var execError = errors.New("execution error")

	// create a new GoogleBreaker
	breaker := New(nil)
	defer breaker.Stop()

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
	assert.NoError(t, err)

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

func TestCircuitBreaker_DoWithFallbackAcceptable_FallbackTrigger(t *testing.T) {
	var (
		execError = errors.New("execution error")
		fbError   = errors.New("fallback error")
	)

	breaker := New(nil)
	defer breaker.Stop()

	// Simulate running 100 times, failed
	for i := 0; i < 100; i++ {
		_ = breaker.Do(func() error {
			return execError
		})
	}

	// Simulate running 1 time, success
	err := breaker.Do(func() error {
		return nil
	})
	assert.ErrorIs(t, err, com.ErrorServiceUnavailable, "Unexpected error")

	// Test case 1: Successful execution with fallback error and acceptable result
	fn := func() error {
		return nil
	}
	fallback := func(err error) error {
		return fbError
	}
	acceptable := func(err error) bool {
		return err == nil
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, fbError, "Unexpected error")

	// Test case 2: Successful execution with fallback error and unacceptable result
	acceptable = func(err error) bool {
		return err != nil
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, fbError, "Unexpected error")

	// Test case 3: Failed execution with fallback error and acceptable result
	fn = func() error {
		return execError
	}
	acceptable = func(err error) bool {
		return err == execError
	}
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, fbError, "Unexpected error")

	// Test case 4: Failed execution with fallback error and unacceptable result
	acceptable = func(err error) bool {
		return err != execError
	}

	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, fbError, "Unexpected error")
}

func TestCircuitBreaker_DoAfterStop(t *testing.T) {
	var execError = errors.New("execution error")

	// Test case 1: Successful execution
	breaker := New(nil)

	// Simulate running 100 times, success
	for i := 0; i < 100; i++ {
		_ = breaker.Do(func() error {
			return nil
		})
	}

	// Stop the breaker
	breaker.Stop()

	// Exec
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
	assert.ErrorIs(t, err, com.ErrorRollingWindowStopped, "Unexpected error")

	// Test case 2: Failed execution
	breaker = New(nil)

	// Simulate running 100 times, failed
	for i := 0; i < 100; i++ {
		_ = breaker.Do(func() error {
			return execError
		})
	}

	// Stop the breaker
	breaker.Stop()

	// Exec
	err = breaker.DoWithFallbackAcceptable(fn, fallback, acceptable)
	assert.ErrorIs(t, err, com.ErrorRollingWindowStopped, "Unexpected error")
}
