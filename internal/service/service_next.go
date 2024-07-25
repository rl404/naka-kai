package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// HandleCommandNext to handle command next.
func (s *service) HandleCommandNext(ctx context.Context, i *discordgo.Interaction) error {
	if s.discord.GetMessageID(i.GuildID) == "" {
		s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)+1)
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Next()))
	}

	s.discord.Stop(i.GuildID)

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Next()))
}

// HandleComponentNext to handle component next.
func (s *service) HandleComponentNext(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Stop(i.GuildID)
	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Next()))
}
