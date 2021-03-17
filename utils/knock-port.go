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
				embedMessage := &discordgo.MessageEmbed{
					Title: "📜 Server Status",
					Color: 0xe74c3c,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "📖 Current Status",
							Value: "UserStyles.world is currently offline.",
						},
						{
							Name:  "❓ Help",
							Value: "Please be patient, admins are looking into it.",
						},
						{
							Name:  "💡 How long will it take",
							Value: "Most of the time, this means the server is updating and should take a couple of minutes.",
						},
					},
				}

				s.ChannelMessageSendEmbed(StatusChannelID, embedMessage)
				LastUptime = time.Now()
				IsDown = true
			}
			if conn != nil && IsDown {
				embedMessage := &discordgo.MessageEmbed{
					Title: "📜 Server Status",
					Color: 0x2ecc71,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "📖 Current Status",
							Value: "UserStyles.world is currently back online!",
						},
						{
							Name:  "⏲️ Duration",
							Value: "The server was out for: " + time.Since(LastUptime).Round(time.Second).String(),
						},
						{
							Name:  "💡 Note",
							Value: "Thank you for being patient.",
						},
						{
							Name:  "🖥️ Website",
							Value: "https://www.v1.userstyles.world/",
						},
					},
				}
				s.ChannelMessageSendEmbed(StatusChannelID, embedMessage)
				LastUptime = time.Now()
			}
			IsDown = conn == nil
			time.Sleep(time.Second * 5)
		}
	}()
}
