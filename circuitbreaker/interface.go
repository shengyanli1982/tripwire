package circuitbreaker

// Callback 是一个接口，定义了熔断器的回调函数。
// Callback is an interface that defines the callback functions of the circuit breaker.
type Callback interface {
	// OnSuccess 在调用成功时被调用。
	// OnSuccess is called when the call is successful.
	OnSuccess(opterr error)

	// OnFailure 在调用失败时被调用。
	// OnFailure is called when the call is failed.
	OnFailure(opterr, reason error)

	// OnAccept 在调用被接受时被调用。
	// OnAccept is called when the call is accepted.
	OnAccept(reason error, refFactor float64)
}

// emptyCallback 是熔断器的空回调。
// emptyCallback is the empty callback for the breaker.
type emptyCallback struct{}

// OnSuccess 是在调用成功时被调用的空操作。
// OnSuccess is nop called when the call is successful.
func (emptyCallback) OnSuccess(opterr error) {}

// OnFailure 是在调用失败时被调用的空操作。
// OnFailure is nop called when the call is failed.
func (emptyCallback) OnFailure(opterr, reason error) {}

// OnAccept 是在调用被接受时被调用的空操作。
// OnAccept is nop called when the call is accepted.
func (emptyCallback) OnAccept(reason error, refFactor float64) {}

// NewEmptyCallback 返回一个新的熔断器空回调。
// NewEmptyCallback returns a new empty callback for the breaker.
func NewEmptyCallback() Callback {
	return &emptyCallback{}
}
