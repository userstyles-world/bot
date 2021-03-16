package handlers

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if !strings.HasPrefix(m.Content, "&") {
		return
	}
	content := strings.TrimPrefix(m.Content, "&")

	if content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
