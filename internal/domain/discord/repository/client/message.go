package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// SendEmbedMessage to send embed message.
func (c *client) SendEmbedMessage(ctx context.Context, channelID string, msg *discordgo.MessageEmbed) error {
	if _, err := c.session.ChannelMessageSendEmbed(channelID, msg); err != nil {
		return stack.Wrap(ctx, err)
	}
	return nil
}

// SendErrorMessage to send error message.
func (c *client) SendErrorMessage(ctx context.Context, channelID string, msg string) error {
	if _, err := c.session.ChannelMessageSendEmbed(channelID, &discordgo.MessageEmbed{
		Title:       "Error",
		Description: msg,
		Color:       entity.ColorRed,
	}); err != nil {
		return stack.Wrap(ctx, err)
	}
	return nil
}
