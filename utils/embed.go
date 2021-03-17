package utils

import "github.com/bwmarrin/discordgo"

type EmbedBuilder struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *EmbedBuilder {
	return &EmbedBuilder{&discordgo.MessageEmbed{}}
}

func (eb *EmbedBuilder) SetTitle(name string) *EmbedBuilder {
	eb.Title = name
	return eb
}

func (eb *EmbedBuilder) SetColor(color int) *EmbedBuilder {
	eb.Color = color
	return eb
}

func (eb *EmbedBuilder) AddField(name, value string) *EmbedBuilder {
	eb.Fields = append(eb.Fields, &discordgo.MessageEmbedField{
		Name:  name,
		Value: value,
	})
	return eb
}
