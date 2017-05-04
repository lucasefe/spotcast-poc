package spoty

import (
	"errors"
	"fmt"
	"strings"

	"github.com/parnurzeal/gorequest"
)

var defaultReturnOn = []string{"login", "logout", "play", "pause", "error", "ap"}

const (
	openSpotifyURL = "https://open.spotify.com"
)

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

	session = &Session{CSRF: csfrToken, OAuth: oauthToken}

	return nil
}

// Status fetches the current status
func Status() (*StatusResult, error) {
	if session == nil {
		return nil, errors.New("Not connected")
	}

	dataFormat := `oauth=%s&csrf=%s&returnafter=1&returnon=%v}`
	data := fmt.Sprintf(dataFormat, session.OAuth, session.CSRF, strings.Join(defaultReturnOn, ","))

	result := &StatusResult{}

	request := gorequest.New()
	_, _, errors := request.Get(getURL(fmt.Sprintf("/remote/status.json?%s", data))).
		Set("Origin", "https://open.spotify.com").
		Send(data).
		EndStruct(result)

	if len(errors) > 0 {
		return nil, fmt.Errorf("Can't fetch status: %+v", errors)
	}

	// fmt.Printf("BODY %+v\n", string(body))

	return result, nil
}

// Play plays a song in the local spotify, provided it's open.
func Play(song string) (*PlayResult, error) {
	return &PlayResult{}, nil
}

// Private stuff
func getOauthToken() (string, error) {
	oauthToken := struct {
		Token string `json:"t"`
	}{}

	request := gorequest.New()
	_, _, errs := request.Get(fmt.Sprintf("%s/token", openSpotifyURL)).EndStruct(&oauthToken)
	if len(errs) > 0 {
		return "", fmt.Errorf("Can't fetch csrf status: %+v", errs)
	}

	return oauthToken.Token, nil
}

func getCSRFToken() (string, error) {
	authToken := struct {
		Token string `json:"token"`
	}{}

	request := gorequest.New()
	_, _, errs := request.Get(getURL("/simplecsrf/token.json")).
		Set("Origin", openSpotifyURL).
		EndStruct(&authToken)

	if len(errs) > 0 {
		return "", fmt.Errorf("Can't fetch csrf status: %+v", errs)
	}

	return authToken.Token, nil
}

func getURL(path string) string {
	return fmt.Sprintf("https://%s:%d%s", generateLocalHostname(), 4370, path)
}

// TODO: It needs to be dynamic.
func generateLocalHostname() string {
	return "lucasefe.spotilocal.com"
}
