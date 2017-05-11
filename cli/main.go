package main

import (
	"flag"
	"fmt"
	"os"

	"gitlab.com/lucasefe/spotcast/spoty"
)

const defaultSong = "spotify:album:6eWtdQm0hSlTgpkbw4LaBG"

var session spoty.Session

func main() {
	session, _ = spoty.NewSession()
	session.SetVerbose()

	if len(os.Args) == 1 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "play":
		song := defaultSong
		if len(os.Args) > 2 {
			song = os.Args[2]
		}
		result, err := session.Play(song)
		if err != nil {
			fmt.Printf("Could not play song: %+v\n", err)
			os.Exit(1)
		}

		printPlaying(result)

	case "pause":
		result, err := session.Pause()
		if err != nil {
			fmt.Printf("Could not pause: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status: %+v", result)
	case "resume":
		result, err := session.Resume()
		if err != nil {
			fmt.Printf("Could not resume: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status: %+v", result)
	case "status":
		result, err := session.Status()
		if err != nil {
			fmt.Printf("Could not get status: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status: %+v", result)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func printPlaying(result *spoty.Result) {
	artist := result.Track.ArtistResource.Name
	album := result.Track.AlbumResource.Name
	song := result.Track.TrackResource.Name

	fmt.Printf("Playing: %s - %s - %s", artist, album, song)
}
