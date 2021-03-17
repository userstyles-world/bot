package handlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"bot/utils"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, utils.Prefix) {
		return
	}
	content := strings.TrimPrefix(m.Content, utils.Prefix)

	if content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
	if content == "uptime" {
		uptime := time.Since(utils.LastUptime).Round(time.Second).String()
		if utils.IsDown {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server has been **offline** for %s.", uptime))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Server has been **online** for %s.", uptime))
		}
	}
}
