package utils

import (
	"bot/modules/config"
	"bot/modules/session"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
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
	return string(PID[4:]), nil
}

func getStartTime() (t time.Time, err error) {
	pid, err := getPID()
	if err != nil {
		return t, err
	}

	// This command will return uptime in seconds for a given PID.
	cmd := exec.Command("ps", "-p"+pid, "-o", "etimes", "--no-heading")
	raw, err := cmd.Output()
	if err != nil {
		return t, err
	}

	// Convert uptime to string and trim space.
	out := strings.TrimSpace(string(raw))

	// Convert to int64 to use it with `time.Unix()`.
	secs, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return t, err
	}

	// Finally, calculate the start time.
	epoch := time.Now().Unix() - secs
	start := time.Unix(epoch, 0)

	// Check if for whatever reason we got initial time.
	if start.Equal(t) {
		return t, fmt.Errorf("failed to convert time %v", start)
	}
	t = start

	return t, nil
}

func Initalize() {
	// Update bot status.
	go func() {
		for {
			start, err := getStartTime()
			if err != nil {
				log.Println(err)
				session.Discord.UpdateGameStatus(0, "USw is offline")
			} else {
				uptime := time.Since(start).Round(time.Second).String()
				session.Discord.UpdateGameStatus(0, "USw has been running for "+uptime)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// Check server status.
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
				log.Println("Waiting for USw server to come back up.")
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
				log.Printf("Waiting for process %s(USw server) to exit.\n", pid)
				if err := command.Start(); err != nil {
					log.Println(err)
				}
				if err := command.Wait(); err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
