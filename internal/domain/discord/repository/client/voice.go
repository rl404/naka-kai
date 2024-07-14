package client

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/errors"
)

// JoinVoiceChannel to join voice channel.
func (c *client) JoinVoiceChannel(ctx context.Context, i *discordgo.Interaction) error {
	// Already in voice channel.
	if c.players[i.GuildID].isInVoiceChannel {
		return nil
	}

	guild, err := c.session.State.Guild(i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	// Look for the user who call the command in voice channels.
	for _, vs := range guild.VoiceStates {
		if vs.UserID != i.Member.User.ID {
			continue
		}

		// Join voice channel.
		vc, err := c.session.ChannelVoiceJoin(guild.ID, vs.ChannelID, false, false)
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		c.players[i.GuildID].Lock()
		c.players[i.GuildID].voice = vc
		c.players[i.GuildID].channelID = i.ChannelID
		c.players[i.GuildID].isInVoiceChannel = true
		c.players[i.GuildID].Unlock()
	}

	if !c.players[i.GuildID].isInVoiceChannel {
		return stack.Wrap(ctx, errors.ErrNotInVC)
	}

	return nil
}

// LeaveVoiceChannel to leave voice channel.
func (c *client) LeaveVoiceChannel(ctx context.Context, i *discordgo.Interaction) error {
	// Not in voice channel.
	if !c.players[i.GuildID].isInVoiceChannel {
		return nil
	}

	// Leave voice channel.
	if err := c.players[i.GuildID].voice.Disconnect(); err != nil {
		return stack.Wrap(ctx, err)
	}

	c.players[i.GuildID].Lock()
	c.players[i.GuildID].isInVoiceChannel = false
	c.players[i.GuildID].Unlock()

	return nil
}
