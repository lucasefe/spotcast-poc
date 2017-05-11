package spoty

import "fmt"

type Track struct {
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
}

func (t *Track) CurrentSongURI() string {
	if t != nil {
		return t.TrackResource.URI
	}

	return ""
}

func (t *Track) CurrentSongTitle() string {
	if t != nil {
		return fmt.Sprintf("%s - %s by %s", t.TrackResource.Name, t.AlbumResource.Name, t.ArtistResource.Name)
	}

	return ""
}
