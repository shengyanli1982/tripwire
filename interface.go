package tripwire

type (
	// AcceptableFunc is a function that checks if the error is acceptable.
	AcceptableFunc = func(err error) bool

	// FallbackFunc is a function that handles the fallback logic.
	FallbackFunc = func(err error) error

	// HandleFunc is a function that handles the execution.
	HandleFunc = func() error
)

// ResultNotifier is an interface that notifies the Breaker of the execution result.
type ResultNotifier interface {
	// Accept tells the Breaker that the call is successful.
	Accept()

	// Reject tells the Breaker that the call is failed.
	Reject(reason error)
}

// Breaker is an interface that represents a circuit breaker.
type Breaker interface {
	// Allow checks if the circuit breaker allows the execution.
	Allow() (ResultNotifier, error)

	// Do executes the function and returns the error.
	Do(fn HandleFunc) error

	// DoWithAcceptable executes the function and returns the error.
	//acceptable - 支持自定义判定执行结果
	DoWithAcceptable(fn HandleFunc, acceptable AcceptableFunc) error

	// DoWithFallback executes the function and returns the error.
	//fallback - 支持自定义快速失败
	DoWithFallback(fn HandleFunc, fallback FallbackFunc) error

	// DoWithFallbackAcceptable executes the function and returns the error.
	//fallback - 支持自定义快速失败
	//acceptable - 支持自定义判定执行结果
	DoWithFallbackAcceptable(fn HandleFunc, fallback FallbackFunc, acceptable AcceptableFunc) error
}
