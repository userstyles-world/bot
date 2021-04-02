package utils

import (
	"net"
	"time"

	"github.com/bwmarrin/discordgo"
)

var port = ":3000"

func Initalize(s *discordgo.Session) {
	go func() {
		for {
			conn, _ := net.DialTimeout("tcp", port, time.Second*4)
			if conn == nil && !IsDown {
				embedMessage := NewEmbed().
					SetTitle("📜 Server Status").
					SetColor(0xe74c3c).
					AddField("📖 Current Status", "UserStyles.world is currently offline.").
					AddField("❓ Help", "Please be patient, admins are looking into it.").
					AddField("💡 Duration", "Most of the time, this means the server is updating and should take a couple of minutes.")
				s.ChannelMessageSendEmbed(StatusChannelID, embedMessage.MessageEmbed)
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
				s.ChannelMessageSendEmbed(StatusChannelID, embedMessage.MessageEmbed)
				LastUptime = time.Now()
			}
			IsDown = conn == nil
			time.Sleep(time.Second * 5)
		}
	}()
}
