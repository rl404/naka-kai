package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// SendInteractionErrorMessage to send interaction error message.
func (c *client) SendInteractionErrorMessage(ctx context.Context, i *discordgo.Interaction, msg string, isEdit ...bool) error {
	_type := discordgo.InteractionResponseChannelMessageWithSource
	if len(isEdit) > 0 && isEdit[0] {
		_type = discordgo.InteractionResponseUpdateMessage
	}

	return stack.Wrap(ctx, c.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: _type,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Error",
					Description: msg,
					Color:       entity.ColorRed,
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	}))
}

// SendInteractionEmbedMessage to send interaction embed message.
func (c *client) SendInteractionEmbedMessage(ctx context.Context, i *discordgo.Interaction, msgs []*discordgo.MessageEmbed, components []discordgo.MessageComponent, isEdit ...bool) error {
	_type := discordgo.InteractionResponseChannelMessageWithSource
	if len(isEdit) > 0 && isEdit[0] {
		_type = discordgo.InteractionResponseUpdateMessage
	}

	return stack.Wrap(ctx, c.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: _type,
		Data: &discordgo.InteractionResponseData{
			Embeds:     msgs,
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	}))
}
