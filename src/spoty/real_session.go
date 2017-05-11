package spoty

import (
	"fmt"
	"net/url"
	"strings"

	"util"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/gommon/log"
	"github.com/parnurzeal/gorequest"
)

const openSpotifyURL = "https://open.spotify.com"

var defaultReturnOn = []string{"login", "logout", "play", "pause", "error", "ap"}

// RealSession is the session data
type RealSession struct {
	csrfToken  string
	oauthToken string
	log        *logrus.Logger
}

// NewSession creates a new Real Session
func NewSession() (*RealSession, error) {
	log := util.NewLogger()

	oauthToken, err := getoauthToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get oauthToken: %+v", err)
	}

	csfrToken, err := getcsrfToken()
	if err != nil {
		return nil, fmt.Errorf("Could not get csrfToken: %+v", err)
	}

	session := &RealSession{
		csrfToken:  csfrToken,
		oauthToken: oauthToken,
		log:        log,
	}

	log.Debugf("Session: %+v\n", session)

	return session, nil
}

// SetVerbose enables verbose level on the logger
func (s *RealSession) SetVerbose() {
	s.log.Level = logrus.DebugLevel
}

// Status fetches the current status
func (s *RealSession) Status() (*Result, error) {
	params := getAuthParams(s)
	params.Set("returnafter", "1")
	params.Set("returnon", strings.Join(defaultReturnOn, ","))

	result, err := getResult("/remote/status.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Play plays a song in the local spotify, provided it's open.
func (s *RealSession) Play(song string) (*Result, error) {
	params := getAuthParams(s)
	params.Set("context", song)
	params.Set("uri", song)

	result, err := getResult("/remote/play.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Pause plays a song in the local spotify, provided it's open.
func (s *RealSession) Pause() (*Result, error) {
	params := getAuthParams(s)
	params.Set("pause", "true")

	result, err := getResult("/remote/pause.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Resume plays a song in the local spotify, provided it's open.
func (s *RealSession) Resume() (*Result, error) {
	params := getAuthParams(s)
	params.Set("pause", "false")

	result, err := getResult("/remote/pause.json", params)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Private stuff
func getoauthToken() (string, error) {
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

func getcsrfToken() (string, error) {
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

func getAuthParams(session *RealSession) *url.Values {
	v := &url.Values{}
	v.Set("oauth", session.oauthToken)
	v.Set("csrf", session.csrfToken)

	return v
}

func getResult(path string, params *url.Values) (*Result, error) {
	log.Debugf("Requesting Path: %s Params: %s", path, params.Encode())

	result := &Result{}
	request := gorequest.New()
	_, _, errors := request.Get(getURL(path)).
		Set("Origin", openSpotifyURL).
		Query(params.Encode()).
		EndStruct(result)

	if len(errors) > 0 {
		return nil, fmt.Errorf("Error getting %s => %+v", path, errors)
	}

	if err := errorOnResult(result); err != nil {
		return nil, fmt.Errorf("Error getting %s => %+v", path, err)
	}

	return result, nil
}

func errorOnResult(result *Result) error {
	if result.Error.Type != "" {
		return fmt.Errorf("Type: %s Message: %s", result.Error.Type, result.Error.Message)
	}

	return nil
}

func getURL(path string) string {
	return fmt.Sprintf("https://%s:%d%s", generateLocalHostname(), 4370, path)
}

// TODO: It needs to be dynamic.
func generateLocalHostname() string {
	return "lucassa1fe.spotilocal.com"
}
