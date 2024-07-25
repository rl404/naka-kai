package service

import (
	"context"
	"sort"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleCommandRemove to handle command remove.
func (s *service) HandleCommandRemove(ctx context.Context, i *discordgo.Interaction) error {
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue)))
	}

	args := s.parseArgs(i)
	songNumbers := make([]int, len(args))
	for j, arg := range args {
		no, err := strconv.Atoi(arg)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidSongNumber)))
		}
		if no <= 0 || no > len(queue) {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidSongNumber)))
		}
		songNumbers[j] = no
	}

	sort.Ints(songNumbers)

	if err := s.queue.DeleteByGuildIDAndOrders(ctx, i.GuildID, songNumbers); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrRemoveSong)))
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Removed(songNumbers)))
}

// HandleCommandPurge to handle command purge.
func (s *service) HandleCommandPurge(ctx context.Context, i *discordgo.Interaction) error {
	if err := s.queue.DeleteByGuildID(ctx, i.GuildID); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrPurgeQueue)), err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Purged()))
}
