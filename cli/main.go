package main

import (
	"fmt"
	"time"

	"bitbucket.org/lucasefe/spotcast/spoty"
)

const song = "spotify:album:6eWtdQm0hSlTgpkbw4LaBG"

func main() {
	spoty.EnableVerbose()

	err := spoty.Connect()
	if err != nil {
		panic(fmt.Sprintf("Could not Connect: %+v\n", err))
	}

	_, err = spoty.Status()
	if err != nil {
		panic(fmt.Sprintf("Could not get status: %+v\n", err))
	}

	_, err = spoty.Play(song)
	if err != nil {
		panic(fmt.Sprintf("Could not play song: %+v\n", err))
	}

	_, err = spoty.Status()
	if err != nil {
		panic(fmt.Sprintf("Could not get status: %+v\n", err))
	}

	time.Sleep(15 * time.Second)

	_, err = spoty.Pause()
	if err != nil {
		panic(fmt.Sprintf("Could not pause song: %+v\n", err))
	}
}
