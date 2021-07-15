package utils

import (
	"bot/modules/config"
	"bot/modules/session"
	"log"
	"net"
	"os/exec"
	"time"
)

// Get the current PID that's listening on the defined port.
func getPID() (string, error) {
	command := exec.Command("lsof", "-n", "sport = :"+config.ServerPort, "|", "grep", "-Po", "(?<=pid=)\\d+")
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func Initalize() {
	log.Println(getPID())
	go func() {
		for {
			conn, _ := net.Dial("tcp", config.ServerPort)
			if conn == nil && !IsDown {
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
			if conn != nil && IsDown {
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
			IsDown = conn == nil
			time.Sleep(time.Second * 1)
		}
	}()
}
