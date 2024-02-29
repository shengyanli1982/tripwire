package circuitbreaker

type Callback interface {
	// OnAccept is called when the call is successful.
	OnAccept(opterr error)

	// OnReject is called when the call is failed.
	OnReject(opterr, reason error)
}

// empty callback for the breaker
type emptyCallback struct{}

// OnAccept is nop called when the call is successful.
func (emptyCallback) OnAccept(opterr error) {}

// OnReject is nop called when the call is failed.
func (emptyCallback) OnReject(opterr, reason error) {}

// NewEmptyCallback returns a new empty callback for the breaker.
func NewEmptyCallback() Callback {
	return &emptyCallback{}
}
