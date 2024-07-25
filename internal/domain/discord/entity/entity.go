package entity

import "github.com/bwmarrin/discordgo"

// Message is entity for message.
type Message struct {
	Messages   []*discordgo.MessageEmbed
	Components []discordgo.MessageComponent
	IsEdit     bool
	AutoDelete bool
}
