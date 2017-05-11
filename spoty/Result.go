package spoty

// Result represent the result of a status call
type Result struct {
	ClientVersion string `json:"client_version"`
	Version       int    `json:"version"`

	Running bool `json:"running"`
	Playing bool `json:"playing"`
	Shuffle bool `json:"shuffle"`
	Repeat  bool `json:"repeat"`

	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	}

	Track Track `json:"track"`
}
