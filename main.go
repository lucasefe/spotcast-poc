package main

import (
	"github.com/lucasefe/spotcast/spoty"

	"fmt"
	"time"
)

// Wrecking..
const song = "spotify:album:6eWtdQm0hSlTgpkbw4LaBG"

func main() {
	err := spoty.Connect()
	if err != nil {
		panic(fmt.Sprintf("Could not Connect: %+v\n", err))
	}

	result, err := spoty.Status()
	if err != nil {
		panic(fmt.Sprintf("Could not get status: %+v\n", err))
	}

	fmt.Printf("Status: %+v\n", result)

	_, err = spoty.Play(song)
	if err != nil {
		panic(fmt.Sprintf("Could not play song: %+v\n", err))
	}

	result, err = spoty.Status()
	if err != nil {
		panic(fmt.Sprintf("Could not get status: %+v\n", err))
	}

	time.Sleep(15 * time.Second)

	_, err = spoty.Pause()
	if err != nil {
		panic(fmt.Sprintf("Could not pause song: %+v\n", err))
	}

	fmt.Printf("Status: %+v\n", result)
}
