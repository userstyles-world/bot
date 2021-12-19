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

func serverDown() {
	LastUptime = time.Now()

	video := "https://cdn.discordapp.com/attachments/821455365274075136/904434361829564416/server.webm"
	session.Discord.ChannelMessageSend(AnnouncementsID, video)

	/*
	embedMessage := NewEmbed().
		SetTitle("📜 Server Status").
		SetColor(0xe74c3c).
		AddField("📖 Current Status", "UserStyles.world is currently offline.").
		AddField("❓ Help", "Please be patient, admins are looking into it.").
		AddField("💡 Duration", "Most of the time, this means the server is updating and should take a couple of minutes.")
	session.Discord.ChannelMessageSendEmbed(StatusChannelID, embedMessage.MessageEmbed)
	*/
}

func serverUp() {
	IsDown = false

	embedMessage := NewEmbed().
		SetTitle("📜 Server Status").
		SetColor(0x2ecc71).
		AddField("📖 Current Status", "UserStyles.world is back online!").
		AddField("⏲️ Duration", "The server was out for: "+time.Since(LastUptime).Round(time.Second).String()).
		AddField("💡 Note", "Thank you for being patient.").
		AddField("🖥️ Website", "https://userstyles.world/")
	session.Discord.ChannelMessageSendEmbed(AnnouncementsID, embedMessage.MessageEmbed)
}

func Initalize() {
	// Update bot status.
	go func() {
		for {
			start, err := getStartTime()
			if err != nil {
				IsDown = true

				log.Println(err)
				session.Discord.UpdateGameStatus(0, "USw is offline 📉")
			} else {
				IsDown = false
				LastUptime = start

				uptime := time.Since(start).Round(time.Second).String()
				session.Discord.UpdateGameStatus(0, "for "+uptime+" 📈")
			}

			if IsDown {
				serverDown()

				for {
					if _, err := getStartTime(); err == nil {
						serverUp()
						break
					}

					time.Sleep(time.Second)
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()
}
