package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
)

// HandleCommandPause to handle command pause.
func (s *service) HandleCommandPause(ctx context.Context, i *discordgo.Interaction) error {
	if s.discord.GetMessageID(i.GuildID) == "" {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Paused()))
	}

	s.discord.Pause(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.PausedButtons()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Paused()))
}

// HandleComponentPause to handle component pause.
func (s *service) HandleComponentPause(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Pause(i.GuildID)

	if err := s.discord.EditMessage(ctx, i.ChannelID, s.discord.GetMessageID(i.GuildID), s.template.PausedButtons()); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Paused()))
}
