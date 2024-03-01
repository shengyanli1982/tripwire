package tripwire

type (
	// AcceptableFunc 是一个检查错误是否可接受的函数。
	// AcceptableFunc is a function that checks if the error is acceptable.
	AcceptableFunc = func(err error) bool

	// FallbackFunc 是一个处理降级逻辑的函数。
	// FallbackFunc is a function that handles the fallback logic.
	FallbackFunc = func(err error) error

	// HandleFunc 是一个处理执行的函数。
	// HandleFunc is a function that handles the execution.
	HandleFunc = func() error

	// RetryableFunc 是一个处理重试逻辑的函数。
	// RetryableFunc is a function that handles the retry logic.
	RetryableFunc = func() (any, error)
)

// Notifier 是一个通知 Breaker 执行结果的接口。
// Notifier is an interface that notifies the Breaker of the execution result.
type Notifier interface {
	// MarkSuccess 告诉 Breaker 调用成功。
	// MarkSuccess tells the Breaker that the call is successful.
	MarkSuccess()

	// MarkFailure 告诉 Breaker 调用失败。
	// MarkFailure tells the Breaker that the call is failed.
	MarkFailure(reason error)
}

// Breaker 是一个表示熔断器的接口。
// Breaker is an interface that represents a circuit breaker.
type Breaker interface {
	// Allow 检查熔断器是否允许执行。
	// Allow checks if the circuit breaker allows the execution.
	Allow() (Notifier, error)

	// Do 执行函数并返回错误。
	// Do executes the function and returns the error.
	Do(fn HandleFunc) error

	// DoWithAcceptable 执行函数并返回错误。
	// DoWithAcceptable executes the function and returns the error.
	DoWithAcceptable(fn HandleFunc, acceptable AcceptableFunc) error

	// DoWithFallback 执行函数并返回错误。
	// DoWithFallback executes the function and returns the error.
	DoWithFallback(fn HandleFunc, fallback FallbackFunc) error

	// DoWithFallbackAcceptable 执行函数并返回错误。
	// DoWithFallbackAcceptable executes the function and returns the error.
	DoWithFallbackAcceptable(fn HandleFunc, fallback FallbackFunc, acceptable AcceptableFunc) error
}

type RetryResult interface {
	// Data 返回结果的数据。
	// Data returns the data of the result.
	Data() any

	// TryError 返回结果的重试错误。
	// TryError returns the retry error of the result.
	TryError() error

	// IsSuccess 检查重试是否成功。
	// IsSuccess checks if the retry is successful.
	IsSuccess() bool

	// Count 返回结果的重试次数。
	// Count returns the retry count of the result.
	Count() int64
}

type Retry interface {
	// TryOnConflict 执行函数并返回重试结果。
	// TryOnConflict executes the function and returns the retry result.
	TryOnConflict(fn RetryableFunc) RetryResult
}

type Throttle interface {
	// Allow 检查节流器是否允许执行。
	// Allow checks if the throttle allows the execution.
	Allow() (Notifier, error)

	// Do 执行函数并返回错误。
	// Do executes the function and returns the error.
	Do(fn HandleFunc, fallback FallbackFunc, acceptable AcceptableFunc) error
}
