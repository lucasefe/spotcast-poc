package spoty

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/parnurzeal/gorequest"
)

var defaultReturnOn = []string{"login", "logout", "play", "pause", "error", "ap"}

const (
	openSpotifyURL = "https://open.spotify.com"
)

// Session is the session data
type Session struct {
	CSRFToken  string
	OAuthToken string
}

var session *Session

// Connect connects to the local server and gets session data
func Connect() error {
	if session != nil {
		return errors.New("Session already created.")
	}

	oauthToken, err := getOauthToken()
	if err != nil {
		return fmt.Errorf("Could not get OAuthToken: %+v", err)
	}

	csfrToken, err := getCSRFToken()
	if err != nil {
		return fmt.Errorf("Could not get CSRFToken: %+v", err)
	}

	session = &Session{CSRFToken: csfrToken, OAuthToken: oauthToken}

	fmt.Printf("Session: %+v\n", session)
	return nil
}

// Status fetches the current status
func Status() (*Result, error) {
	if session == nil {
		return nil, errors.New("Not connected")
	}

	params := getAuthParams(session)
	params.Set("returnafter", "1")
	params.Set("returnon", strings.Join(defaultReturnOn, ","))

	result, err := getResult("/remote/status.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Play plays a song in the local spotify, provided it's open.
func Play(song string) (*Result, error) {
	if session == nil {
		return nil, errors.New("Not connected")
	}

	params := getAuthParams(session)
	params.Set("context", song)
	params.Set("uri", song)

	result, err := getResult("/remote/play.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Pause plays a song in the local spotify, provided it's open.
func Pause() (*Result, error) {
	if session == nil {
		return nil, errors.New("Not connected")
	}

	params := getAuthParams(session)
	params.Set("pause", "true")

	result, err := getResult("/remote/pause.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Private stuff
func getOauthToken() (string, error) {
	data := struct {
		Token string `json:"t"`
	}{}

	request := gorequest.New()
	_, _, errs := request.Get(fmt.Sprintf("%s/token", openSpotifyURL)).EndStruct(&data)
	if len(errs) > 0 {
		return "", fmt.Errorf("Can't fetch csrf status: %+v", errs)
	}

	return data.Token, nil
}

func getCSRFToken() (string, error) {
	data := struct {
		Token string `json:"token"`
	}{}

	request := gorequest.New()
	_, _, errs := request.Get(getURL("/simplecsrf/token.json")).
		Set("Origin", openSpotifyURL).
		EndStruct(&data)

	if len(errs) > 0 {
		return "", fmt.Errorf("Can't fetch csrf status: %+v", errs)
	}

	return data.Token, nil
}

func getAuthParams(session *Session) *url.Values {
	v := &url.Values{}
	v.Set("oauth", session.OAuthToken)
	v.Set("csrf", session.CSRFToken)

	return v
}

func getResult(path string, params *url.Values) (*Result, error) {
	result := &Result{}
	request := gorequest.New()
	_, _, errors := request.Get(getURL(path)).
		Set("Origin", openSpotifyURL).
		Query(params.Encode()).
		EndStruct(result)

	if len(errors) > 0 {
		return nil, fmt.Errorf("Error getting %s => %+v", path, errors)
	}

	return result, nil
}

func getURL(path string) string {
	return fmt.Sprintf("https://%s:%d%s", generateLocalHostname(), 4370, path)
}

// TODO: It needs to be dynamic.
func generateLocalHostname() string {
	return "lucassa1fe.spotilocal.com"
}
