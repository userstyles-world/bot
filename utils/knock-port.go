package utils

import (
	"bot/modules/config"
	"bot/modules/session"
	"errors"
	"log"
	"os/exec"
	"regexp"
	"time"
)

var (
	ssPIDRegex    = regexp.MustCompile(`pid=\d+`)
	ErrNoPIDFound = errors.New("Could not get PID")
)

// Get the current PID that's listening on the defined port.
func getPID() (string, error) {
	// ss -lptn "sport = :PORT"
	command := exec.Command("ss", "-lptn", "sport = :"+config.ServerPort)
	output, err := command.Output()
	if err != nil {
		return "", ErrNoPIDFound
	}
	// Extract the pid=12345 part of the output.
	PID := ssPIDRegex.Find(output)

	// E.g. PID = "pid=" or ""
	if len(PID) == 0 || len(PID) <= 4 {
		return "", ErrNoPIDFound
	}
	// remove the PID from the output
	return string(PID[4 : len(PID)-1]), nil
}

func Initalize() {
	go func() {
		for {
			pid, err := getPID()
			if err != nil {
				log.Println(err)
			}
			if pid == "" && !IsDown {
				embedMessage := NewEmbed().
					SetTitle("📜 Server Status").
					SetColor(0xe74c3c).
					AddField("📖 Current Status", "UserStyles.world is currently offline.").
					AddField("❓ Help", "Please be patient, admins are looking into it.").
					AddField("💡 Duration", "Most of the time, this means the server is updating and should take a couple of minutes.")
				session.Discord.ChannelMessageSendEmbed(StatusChannelID, embedMessage.MessageEmbed)
				LastUptime = time.Now()
				IsDown = true
			}
			if pid != "" && IsDown {
				embedMessage := NewEmbed().
					SetTitle("📜 Server Status").
					SetColor(0x2ecc71).
					AddField("📖 Current Status", "UserStyles.world is currently back online!").
					AddField("⏲️ Duration", "The server was out for: "+time.Since(LastUptime).Round(time.Second).String()).
					AddField("💡 Note", "Thank you for being patient.").
					AddField("🖥️ Website", "https://userstyles.world/")
				session.Discord.ChannelMessageSendEmbed(StatusChannelID, embedMessage.MessageEmbed)
				LastUptime = time.Now()
			}
			IsDown = pid == ""
			if IsDown {
				for {
					if maybePID, _ := getPID(); maybePID != "" {
						break
					}
					time.Sleep(time.Second)
				}
			} else {
				// Wait until the process is killed or exited.
				// tail --pid=$pid -f /dev/null
				command := exec.Command("tail", "--pid="+pid, "-f", "/dev/null")

				if err := command.Run(); err != nil {
					log.Println(err)
				}	
			}
		}
	}()
}
