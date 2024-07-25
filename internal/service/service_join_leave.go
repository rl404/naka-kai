package service

import (
	"context"
	_errors "errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/errors"
)

// HandleCommandJoin to handle command join.
func (s *service) HandleCommandJoin(ctx context.Context, i *discordgo.Interaction) error {
	if err := stack.Wrap(ctx, s.discord.JoinVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.NotInVC)), err)
		}
		return stack.Wrap(ctx, err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Joined()))
}

// HandleCommandLeave to handle command leave.
func (s *service) HandleCommandLeave(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.SetAutoNext(i.GuildID, false)
	s.discord.Stop(i.GuildID)

	if err := stack.Wrap(ctx, s.discord.LeaveVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.NotInVC)), err)
		}
		return stack.Wrap(ctx, err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Left()))
}
