package spoty

import (
	"errors"
	"fmt"
)

// StatusResult represent the result of a status call
type StatusResult struct{}

// PlayResult represent the result of a play call
type PlayResult struct{}

// Session is the session data
type Session struct {
	CSRF  string
	OAuth string
}

var session *Session

// Connect connects to the local server and gets session data
func Connect() error {

	if session != nil {
		return errors.New("Session already created.")
	}

	oauthToken, err := getOauthToken()
	if err != nil {
		return fmt.Errorf("Could not get oauthToken: %+v", err)
	}

	csfrToken, err := getCSRFToken()
	if err != nil {
		return fmt.Errorf("Could not get CSRFToken: %+v", err)
	}

	session := &Session{CSRF: csfrToken, OAuth: oauthToken}

	fmt.Printf("Session: %+v", session)

	return nil

}

// Status fetches the current status
func Status() (*StatusResult, error) {
	return &StatusResult{}, nil
}

// Play plays a song in the local spotify, provided it's open.
func Play(song string) (*PlayResult, error) {
	return &PlayResult{}, nil
}

func getOauthToken() (string, error) {
	return "", nil
}

func getCSRFToken() (string, error) {
	return "", nil
}
