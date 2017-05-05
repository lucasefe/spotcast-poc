package spoty

// StatusResult represent the result of a status call
type Result struct {
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
