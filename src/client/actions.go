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
	resumeActionType   = "RESUME"
	setRoleActiontype  = "SET_ROLE"
)

func handleAction(message string) {
	action, err := parseAction(message)
	if err != nil {
		logger.Warningf("Could not parse message %+v. Got error: %s", message, err)
	}

	if role == LeaderRole {
		logger.Debug("Skipping player action. You're the leader.")
		return
	}

	switch action.Type {
	case setRoleActiontype:
		setRole(action.Data)
		break
	case playSongActionType:
		playSong(action.Data)
		break
	case stopActionType, pauseActionType:
		pausePlayer()
		break
	case resumeActionType:
		resumePlayer()
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

func pauseAction() []byte {
	action := &Action{Type: pauseActionType}
	message, _ := json.Marshal(action)
	return message
}

func resumeAction() []byte {
	action := &Action{Type: resumeActionType}
	message, _ := json.Marshal(action)
	return message
}

// action := newPlayAction("spotify:artist:08td7MxkoHQkXnWAYD8d6Q")
func playSongAction(song string) []byte {
	action := newPlayAction(song)
	message, _ := json.Marshal(action)
	return message
}

func playSong(data map[string]string) {
	song, ok := data["song"]

	if !ok {
		logger.Warningf("Wrong data: %+v", data)
		return
	}

	logger.Infof("Playing song: %s", song)
	player.Play(song)

}

func pausePlayer() {
	logger.Info("Pausing player")
	player.Pause()
}

func resumePlayer() {
	logger.Info("Resuming player")
	player.Resume()
}

func setRole(data map[string]string) {
	if newRole, ok := data["role"]; ok {
		logger.Infof("Setting new role: %s", newRole)
		role = Role(newRole)
	}
}
