package spoty

// StatusResult represent the result of a status call
type StatusResult struct{}

// PlayResult represent the result of a play call
type PlayResult struct{}

// Status fetches the current status
func Status() (*StatusResult, error) {
	return &StatusResult{}, nil
}

// Play plays a song in the local spotify, provided it's open.
func Play(song string) (*PlayResult, error) {
	return &PlayResult{}, nil
}
