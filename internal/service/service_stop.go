package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// HandleCommandStop to handle command stop.
func (s *service) HandleCommandStop(ctx context.Context, i *discordgo.Interaction) error {
	if s.discord.GetMessageID(i.GuildID) == "" {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Stopped()))
	}

	s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)-1)
	s.discord.SetAutoNext(i.GuildID, false)
	s.discord.Stop(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.Stopped()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Stopped()))
}

// HandleComponentStop to handle component stop.
func (s *service) HandleComponentStop(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)-1)
	s.discord.SetAutoNext(i.GuildID, false)
	s.discord.Stop(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.Stopped()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Stopped()))
}
