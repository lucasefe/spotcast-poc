package main

import (
	"flag"
	"fmt"
	"os"

	"bitbucket.org/lucasefe/spotcast/spoty"
)

const defaultSong = "spotify:album:6eWtdQm0hSlTgpkbw4LaBG"

func main() {
	spoty.EnableVerbose()
	spoty.Connect()

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
		result, err := spoty.Play(song)
		if err != nil {
			fmt.Printf("Could not play song: %+v\n", err)
			os.Exit(1)
		}

		printPlaying(result)

	case "pause":
		result, err := spoty.Pause()
		if err != nil {
			fmt.Printf("Could not pause: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status: %+v", result)
	case "resume":
		result, err := spoty.Resume()
		if err != nil {
			fmt.Printf("Could not resume: %+v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status: %+v", result)
	case "status":
		result, err := spoty.Status()
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
