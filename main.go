package main

import (
	"github.com/lucasefe/spotcast/spoty"

	"fmt"
)

// Wrecking..
const song = "0TB7xPRIQ6sZqH8q50maWh"

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

	fmt.Printf("Status: %+v\n", result)
}
