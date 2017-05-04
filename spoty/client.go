package spoty

import (
	"errors"
	"fmt"
	"strings"

	"github.com/parnurzeal/gorequest"
)

// StatusResult represent the result of a status call
type StatusResult struct {
	ClientVersion string `json:"client_version"`
	Version       int    `json:"version"`

	Running bool `json:"running"`
	Playing bool `json:"playing"`
	Shuffle bool `json:"shuffle"`
	Repeat  bool `json:"repeat"`

	Track struct {
		TrackResource struct {
			Name     string `json:"name"`
			URI      string `json:"uri"`
			Location struct {
				OG string `json:"og"`
			} `json:"location"`
		} `json:"track_resource"`
		ArtistResource struct {
			Name     string `json:"name"`
			URI      string `json:"uri"`
			Location struct {
				OG string `json:"og"`
			} `json:"location"`
		} `json:"artist_resource"`
		AlbumResource struct {
			Name     string `json:"name"`
			URI      string `json:"uri"`
			Location struct {
				OG string `json:"og"`
			} `json:"location"`
		} `json:"album_resource"`
	} `json:"track"`
}

// PlayResult represent the result of a play call
type PlayResult struct{}

// Session is the session data
type Session struct {
	CSRF  string
	OAuth string
}

var session *Session

var defaultReturnOn = []string{"login", "logout", "play", "pause", "error", "ap"}

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

	fmt.Printf("Session: %+v\n", session)

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
	_, body, errors := request.Get(getURL(fmt.Sprintf("/remote/status.json?%s", data))).
		Set("Origin", "https://open.spotify.com").
		Send(data).
		EndStruct(result)

	if len(errors) > 0 {
		return nil, fmt.Errorf("Can't fetch status: %+v", errors)
	}

	fmt.Printf("BODY %+v\n", string(body))

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
	request.Get("https://open.spotify.com/token").
		EndStruct(&oauthToken)

	return oauthToken.Token, nil
}

func getCSRFToken() (string, error) {
	authToken := struct {
		Token string `json:"token"`
	}{}

	request := gorequest.New()
	request.Get(getURL("/simplecsrf/token.json")).
		Set("Origin", "https://open.spotify.com").
		EndStruct(&authToken)

	return authToken.Token, nil
}

func getURL(path string) string {
	return fmt.Sprintf("https://%s:%d%s", generateLocalHostname(), 4370, path)
}

func generateLocalHostname() string {
	return "lucasefe.spotilocal.com"
}
