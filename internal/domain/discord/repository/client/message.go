package client

import (
	"context"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/discord/entity"
)

// SendMessage to send message.
func (c *client) SendMessage(ctx context.Context, channelID string, data entity.Message) (string, string, error) {
	msg, err := c.session.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embeds:     data.Messages,
		Components: data.Components,
	})
	if err != nil {
		return "", "", stack.Wrap(ctx, err)
	}

	if data.AutoDelete {
		defer c.deleteMessage(msg.ChannelID, msg.ID)
	}

	return msg.ChannelID, msg.ID, nil
}

// EditMessage to edit message.
func (c *client) EditMessage(ctx context.Context, channelID, messageID string, data entity.Message) error {
	if _, err := c.session.ChannelMessageEditComplex(&discordgo.MessageEdit{
		ID:         messageID,
		Channel:    channelID,
		Embeds:     &data.Messages,
		Components: &data.Components,
	}); err != nil {
		return stack.Wrap(ctx, err)
	}

	if data.AutoDelete {
		defer c.deleteMessage(channelID, messageID)
	}

	return nil
}

func (c *client) deleteMessage(channelID, messageID string) {
	if c.deleteTime == 0 {
		return
	}

	go func() {
		time.Sleep(c.deleteTime)
		c.session.ChannelMessageDelete(channelID, messageID)
	}()
}
