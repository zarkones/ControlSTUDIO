package main

import (
	"common/profiles"
	"fmt"
	"time"
)

// handleErr handles errors of the main function only.
// It handles unrecoverable errors and for most cases you don't
// need to use this function.
func handleErr(err error) {
	if err == nil {
		return
	}
	panic(err)
}

func main() {
	profile, err := profiles.Load(rawProfile)
	handleErr(err)

	for range time.Tick(tickAmount(profile.Receive.SleepMin, profile.Receive.SleepMax)) {
		instruction, err := receive(&profile)
		if err != nil {
			fmt.Println("receive error:", err)
			continue
		}

		fmt.Println("instruction:", instruction)
	}
}
