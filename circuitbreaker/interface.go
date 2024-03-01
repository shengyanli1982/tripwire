package circuitbreaker

type Callback interface {
	// OnSuccess is called when the call is successful.
	OnSuccess(opterr error)

	// OnFailed is called when the call is failed.
	OnFailed(opterr, reason error)

	// OnAccept is called when the call is accepted.
	OnAccept(reason error, refFactor float64)
}

// empty callback for the breaker
type emptyCallback struct{}

// OnSuccess is nop called when the call is successful.
func (emptyCallback) OnSuccess(opterr error) {}

// OnFailed is nop called when the call is failed.
func (emptyCallback) OnFailed(opterr, reason error) {}

// OnAccept is nop called when the call is accepted.
func (emptyCallback) OnAccept(reason error, refFactor float64) {}

// NewEmptyCallback returns a new empty callback for the breaker.
func NewEmptyCallback() Callback {
	return &emptyCallback{}
}
