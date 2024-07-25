package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// HandleCommandPrevious to handle command previous.
func (s *service) HandleCommandPrevious(ctx context.Context, i *discordgo.Interaction) error {
	if s.discord.GetMessageID(i.GuildID) == "" {
		s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)-1)
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Previous()))
	}

	s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)-2)
	s.discord.Stop(i.GuildID)

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Previous()))
}

// HandleComponentPrevious to handle component previous.
func (s *service) HandleComponentPrevious(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)-2)
	s.discord.Stop(i.GuildID)
	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Previous()))
}
