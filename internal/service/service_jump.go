package service

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleCommandJump to handle command jump.
func (s *service) HandleCommandJump(ctx context.Context, i *discordgo.Interaction) error {
	args := s.parseArgs(i)
	no, err := strconv.Atoi(args[0])
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidSongNumber)))
	}

	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue)))
	}

	if len(queue) == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Queue(nil, 0)))
	}

	if no <= 0 || no > len(queue) {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidSongNumber)))
	}

	if s.discord.GetMessageID(i.GuildID) == "" {
		s.discord.SetQueueIndex(i.GuildID, no-1)
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Jumped(no)))
	}

	s.discord.SetQueueIndex(i.GuildID, no-2)
	s.discord.Stop(i.GuildID)

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Jumped(no)))
}
