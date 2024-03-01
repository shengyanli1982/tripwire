package circuitbreaker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleBreaker_Accept(t *testing.T) {
	// create a new GoogleBreaker
	breaker := NewGoogleBreaker(nil)
	defer breaker.Stop()

	// Simulate running 100 times, failed
	for i := 0; i < 100; i++ {
		err := breaker.rwin.Add(0)
		assert.Nil(t, err)
	}

	// Simulate running 1 time, success
	err := breaker.rwin.Add(1)
	assert.Nil(t, err)

	// Test the Accept function, the success rate is 1/101, trigger the ErrorServiceUnavailable error
	// fuse ratio is 0.926, must greater than 0.4
	err = breaker.accept(0.4)
	assert.ErrorIs(t, err, ErrorServiceUnavailable, "unexpected error returned by accept")

	// Test the Accept function, the success rate is 1/101, trigger the ErrorServiceUnavailable error
	// fuse ratio is 0.926, equal the random float64
	err = breaker.accept(0.926)
	assert.ErrorIs(t, err, ErrorServiceUnavailable, "unexpected error returned by accept")

	// create a new GoogleBreaker
	breaker = NewGoogleBreaker(nil)
	defer breaker.Stop()

	// Simulate running 100 times, success
	for i := 0; i < 100; i++ {
		err := breaker.rwin.Add(1)
		assert.Nil(t, err)
	}

	// Simulate running 1 time, failed
	err = breaker.rwin.Add(0)
	assert.Nil(t, err)

	// Test the Accept function, the success rate is 100/101, no error
	err = breaker.accept(0.4)
	assert.NoError(t, err, "unexpected error returned by accept")
}

func TestGoogleBreaker_Allow(t *testing.T) {
	breaker := NewGoogleBreaker(nil)
	defer breaker.Stop()

	// Test allowing execution
	notifier, err := breaker.Allow()
	assert.NoError(t, err, "Unexpected error")
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// Simulate running 100 times, success
	for i := 0; i < 100; i++ {
		err := breaker.rwin.Add(1)
		assert.Nil(t, err)
	}

	// Simulate running 1 time, failed
	err = breaker.rwin.Add(0)
	assert.Nil(t, err)

	// Test allowing execution
	notifier, err = breaker.Allow()
	assert.NoError(t, err, "Unexpected error")
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// Test rejecting execution
	notifier.MarkFailure(errors.New("test"))

	// Test values execution
	v, c, _ := breaker.history()
	assert.Equal(t, float64(100), v, "Expected 100, but got %v", v)
	assert.Equal(t, uint64(102), c, "Expected 102, but got %v", c)

	// Test allowing execution
	notifier, err = breaker.Allow()
	assert.NoError(t, err, "Unexpected error")

	// Verify the returned notifier
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// Test accepting execution
	notifier.MarkSuccess()

	// Test values execution
	v, c, _ = breaker.history()
	assert.Equal(t, float64(101), v, "Expected 101, but got %v", v)
	assert.Equal(t, uint64(103), c, "Expected 103, but got %v", c)
}

func TestGoogleBreaker_DoWithFallbackAcceptable(t *testing.T) {
	var execError = errors.New("execution error")

	// create a new GoogleBreaker
	breaker := NewGoogleBreaker(nil)
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
	defer breaker.Stop()

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
	defer breaker.Stop()

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
	defer breaker.Stop()

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

func TestGoogleBreaker_DoWithFallbackAcceptable_FallbackTrigger(t *testing.T) {
	var (
		execError = errors.New("execution error")
		fbError   = errors.New("fallback error")
	)

	breaker := NewGoogleBreaker(nil)
	defer breaker.Stop()

	// Simulate running 100 times, failed
	for i := 0; i < 100; i++ {
		err := breaker.rwin.Add(0)
		assert.Nil(t, err)
	}

	// Simulate running 1 time, success
	err := breaker.rwin.Add(1)
	assert.Nil(t, err)

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

type testCallback struct {
	factor float64
	sc, fc int
}

func (t *testCallback) OnSuccess(opterr error) {
	t.sc += 1
}

func (t *testCallback) OnFailure(opterr, reason error) {
	t.fc += 1
}

func (t *testCallback) OnAccept(reason error, refFactor float64) {
	t.factor = refFactor
}

func newTestCallback() Callback {
	return &testCallback{}
}

func TestGoogleBreaker_Callback(t *testing.T) {
	config := NewConfig().WithCallback(newTestCallback())
	breaker := NewGoogleBreaker(config)
	defer breaker.Stop()

	// Test allowing execution
	notifier, err := breaker.Allow()
	assert.NoError(t, err, "Unexpected error")
	assert.NotNil(t, notifier, "Expected a notifier, but got nil")

	// callback
	cb := breaker.config.callback.(*testCallback)

	// Test case 1: OnSuccess
	notifier.MarkSuccess()
	assert.Equal(t, 1, cb.sc, "Unexpected success count")

	// Test case 2: OnFailure
	notifier.MarkFailure(errors.New("test"))
	assert.Equal(t, 1, cb.fc, "Unexpected failed count")

	// Test case 3: OnAccept
	err = breaker.accept(0.5)
	assert.NoError(t, err, "Unexpected error")
	assert.Equal(t, -1.5, cb.factor, "Unexpected reference factor")
}
