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
			panic(fmt.Sprintf("Could not play song: %+v\n", err))
		}
		fmt.Printf("Status: %+v", result)
	case "pause":
		result, err := spoty.Pause()
		if err != nil {
			panic(fmt.Sprintf("Could not get status: %+v\n", err))
		}

		fmt.Printf("Status: %+v", result)
	case "status":
		result, err := spoty.Status()
		if err != nil {
			panic(fmt.Sprintf("Could not get status: %+v\n", err))
		}

		fmt.Printf("Status: %+v", result)
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
