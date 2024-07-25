package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleCommandQueue to handle command queue.
func (s *service) HandleCommandQueue(ctx context.Context, i *discordgo.Interaction) error {
	args := s.parseArgs(i)
	if len(args) > 0 {
		return stack.Wrap(ctx, s.search(ctx, i, args, false))
	}

	// Get queue.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue)), err)
	}

	videos := make([]entity.Video, len(queue))
	for i, v := range queue {
		videos[i] = entity.Video{
			VideoTitle:   v.VideoTitle,
			ChannelTitle: v.ChannelTitle,
			Order:        v.Order,
			AddedBy:      v.AddedBy,
		}
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Queue(videos, s.discord.GetQueueIndex(i.GuildID))))
}
