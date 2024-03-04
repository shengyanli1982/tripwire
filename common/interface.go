package common

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

	// Notifier 是一个通知 Breaker 执行结果的接口。
	// Notifier is an interface that notifies the Breaker of the execution result.
	Notifier = interface {
		// MarkSuccess 告诉 Breaker 调用成功。
		// MarkSuccess tells the Breaker that the call is successful.
		MarkSuccess()

		// MarkFailure 告诉 Breaker 调用失败。
		// MarkFailure tells the Breaker that the call is failed.
		MarkFailure(reason error)
	}

	// Breaker 是一个表示熔断器的接口。
	// Breaker is an interface that represents a circuit breaker.
	Breaker = interface {
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

		// Stop 停止熔断器。
		// Stop stops the circuit breaker.
		Stop()
	}

	// RetryResult 接口定义了执行结果的相关方法
	// The RetryResult interface defines methods related to execution results
	RetryResult = interface {
		// Data 方法返回执行结果的数据
		// The Data method returns the data of the execution result
		Data() any

		// TryError 方法返回尝试执行时的错误
		// The TryError method returns the error when trying to execute
		TryError() error

		// ExecErrors 方法返回所有执行错误的列表
		// The ExecErrors method returns a list of all execution errors
		ExecErrors() []error

		// IsSuccess 方法返回执行是否成功
		// The IsSuccess method returns whether the execution was successful
		IsSuccess() bool

		// LastExecError 方法返回最后一次执行的错误
		// The LastExecError method returns the error of the last execution
		LastExecError() error

		// FirstExecError 方法返回第一次执行的错误
		// The FirstExecError method returns the error of the first execution
		FirstExecError() error

		// ExecErrorByIndex 方法返回指定索引处的执行错误
		// The ExecErrorByIndex method returns the execution error at the specified index
		ExecErrorByIndex(idx int) error

		// Count 方法返回执行的次数
		// The Count method returns the number of executions
		Count() int64
	}

	// Retry 是一个接口，定义了一个方法 TryOnConflictVal。
	// Retry is an interface that defines a method TryOnConflictVal.
	Retry = interface {
		// TryOnConflictVal 执行函数并返回重试结果。
		// TryOnConflictVal executes the function and returns the retry result.
		TryOnConflictVal(fn RetryableFunc) RetryResult
	}
)
