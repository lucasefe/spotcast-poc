package spoty

// Session is the spotify session interface
type Session interface {
	Pause() (*Result, error)
	Play(string) (*Result, error)
	Resume() (*Result, error)
	SetVerbose()
	Status() (*Result, error)
}

// FakedSession represents a spotify session that does not connect to the local client
// Useful for dev mode
type FakedSession struct{}

// Pause implements Session.Pause
func (*FakedSession) Pause() (*Result, error) { return &Result{}, nil }

// Play implements Session.Play
func (*FakedSession) Play(string) (*Result, error) { return &Result{}, nil }

// Resume implements Session.Resume
func (*FakedSession) Resume() (*Result, error) { return &Result{}, nil }

// Status implements Session.Status
func (*FakedSession) Status() (*Result, error) { return &Result{}, nil }

// SetVerbose implements Session.SetVerbose
func (*FakedSession) SetVerbose() {}
