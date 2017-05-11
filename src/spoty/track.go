package spoty

import "fmt"

// Track is a track from a result, and includes artist, album and song data
type Track struct {
	TrackResource  trackResource  `json:"track_resource"`
	ArtistResource artistResource `json:"artist_resource"`
	AlbumResource  albumResource  `json:"album_resource"`
}

type trackResource struct {
	Name     string `json:"name"`
	URI      string `json:"uri"`
	Location struct {
		OG string `json:"og"`
	} `json:"location"`
}

type artistResource struct {
	Name     string `json:"name"`
	URI      string `json:"uri"`
	Location struct {
		OG string `json:"og"`
	} `json:"location"`
}
type albumResource struct {
	Name     string `json:"name"`
	URI      string `json:"uri"`
	Location struct {
		OG string `json:"og"`
	} `json:"location"`
}

// CurrentSongURI returns the URI from the song being played right now on the player
func (t *Track) CurrentSongURI() string {
	if t != nil {
		return t.TrackResource.URI
	}

	return ""
}

// CurrentSongTitle returns the title from the song being played right now on the player
func (t *Track) CurrentSongTitle() string {
	if t != nil {
		return fmt.Sprintf("%s - %s by %s", t.TrackResource.Name, t.AlbumResource.Name, t.ArtistResource.Name)
	}

	return ""
}
