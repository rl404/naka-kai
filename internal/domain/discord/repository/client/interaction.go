package client

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/discord/entity"
)

// SendInteractionMessage to send interaction message.
func (c *client) SendInteractionMessage(ctx context.Context, i *discordgo.Interaction, data entity.Message) error {
	if !data.DisableAutoDelete {
		defer c.deleteInteraction(i)
	}

	_type := discordgo.InteractionResponseChannelMessageWithSource
	if data.IsEdit {
		_type = discordgo.InteractionResponseUpdateMessage
	}

	return stack.Wrap(ctx, c.session.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: _type,
		Data: &discordgo.InteractionResponseData{
			Embeds:     data.Messages,
			Components: data.Components,
			Flags:      discordgo.MessageFlagsSuppressNotifications,
		},
	}))
}

func (c *client) deleteInteraction(i *discordgo.Interaction) {
	if c.deleteTime == 0 {
		return
	}

	go func() {
		time.Sleep(c.deleteTime)
		c.session.InteractionResponseDelete(i)
	}()
}
