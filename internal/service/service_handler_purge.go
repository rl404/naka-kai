package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleCommandPurge to handle purge command.
func (s *service) HandleCommandPurge(ctx context.Context, i *discordgo.Interaction) error {
	if err := s.queue.DeleteByGuildID(ctx, i.GuildID); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrPurgeQueue), err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetPurged()}, nil))
}

// HandleComponentPurge to handle purge component.
func (s *service) HandleComponentPurge(ctx context.Context, i *discordgo.Interaction) error {
	if err := s.queue.DeleteByGuildID(ctx, i.GuildID); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrPurgeQueue, true), err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetPurged()}, nil, true))
}
