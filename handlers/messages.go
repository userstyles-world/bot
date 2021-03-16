package handlers

import (
	"strings"

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
}
