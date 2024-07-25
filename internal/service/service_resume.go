package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// HandleCommandResume to handle command resume.
func (s *service) HandleCommandResume(ctx context.Context, i *discordgo.Interaction) error {
	if s.discord.GetMessageID(i.GuildID) == "" {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Playing()))
	}

	s.discord.Resume(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.PlayingButtons()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Playing()))
}

// HandleComponentResume to handle component resume.
func (s *service) HandleComponentResume(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Resume(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.PlayingButtons()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Playing()))
}
