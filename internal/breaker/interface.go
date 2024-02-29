package breaker

type (
	// AcceptableFunc is a function that checks if the error is acceptable.
	AcceptableFunc = func(err error) bool

	// FallbackFunc is a function that handles the fallback logic.
	FallbackFunc = func(err error) error

	// HandleFunc is a function that handles the execution.
	HandleFunc = func() error
)

func DefaultAcceptableFunc(err error) bool {
	return err == nil
}

func DefaultFallbackFunc(err error) error {
	return err
}

type Result interface {
	// Accept tells the Breaker that the call is successful.
	Accept()

	// Reject tells the Breaker that the call is failed.
	Reject(reason error)
}

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
