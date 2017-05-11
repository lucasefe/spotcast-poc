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
type FakedSession struct {
	lastResult *Result
}

var fakeTrack = Track{
	TrackResource:  trackResource{Name: "Fake Track", URI: "spotify:track:0ZfM5XfJTtFPhOxAERRnNY"},
	AlbumResource:  albumResource{Name: "Fake Album", URI: ""},
	ArtistResource: artistResource{Name: "Fake Artist", URI: ""},
}

// NewFakedSession creates a new FakedSession
func NewFakedSession() *FakedSession {
	result := &Result{Running: true, Playing: false, Track: fakeTrack}
	return &FakedSession{lastResult: result}
}

// Pause implements Session.Pause
func (f *FakedSession) Pause() (*Result, error) {
	f.lastResult.Playing = false
	return f.lastResult, nil
}

// Play implements Session.Play
func (f *FakedSession) Play(songURI string) (*Result, error) {
	f.lastResult.Playing = true
	f.lastResult.Track.TrackResource.URI = songURI
	return f.lastResult, nil
}

// Resume implements Session.Resume
func (f *FakedSession) Resume() (*Result, error) {
	f.lastResult.Playing = true
	return f.lastResult, nil
}

// Status implements Session.Status
func (f *FakedSession) Status() (*Result, error) {
	return f.lastResult, nil
}

// SetVerbose implements Session.SetVerbose
func (*FakedSession) SetVerbose() {}
