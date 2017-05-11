package main

import (
	"encoding/json"
)

// Action ..
type Action struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

const (
	pauseActionType    = "PAUSE"
	playSongActionType = "PLAY_SONG"
	stopActionType     = "STOP"
)

func handleAction(message string) {
	action, err := parseAction(message)
	if err != nil {
		logger.Warningf("Could not parse message %+v. Got error: %s", message, err)
	}

	switch action.Type {
	case playSongActionType:
		playSong(action.Data)
		break
	case stopActionType, pauseActionType:
		player.Pause()
		break
	default:
		logger.Warningf("Don't know what this means: %+v", message)
	}
}

func parseAction(message string) (*Action, error) {
	action := &Action{}

	if err := json.Unmarshal([]byte(message), &action); err != nil {
		return nil, err
	}

	return action, nil
}

func newPlayAction(songURI string) *Action {
	return &Action{
		Type: playSongActionType,
		Data: map[string]string{"song": songURI},
	}
}

func pauseAction() ([]byte, error) {
	action := &Action{
		Type: pauseActionType,
	}

	message, err := json.Marshal(action)

	if err != nil {
		return nil, err
	}

	return message, nil
}

func playSongAction(song string) ([]byte, error) {
	// action := newPlayAction("spotify:artist:08td7MxkoHQkXnWAYD8d6Q")
	action := newPlayAction(song)
	message, err := json.Marshal(action)

	if err != nil {
		return nil, err
	}

	return message, nil
}
