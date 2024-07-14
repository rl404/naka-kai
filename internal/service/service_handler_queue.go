package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleCommandQueue to handle queue command.
func (s *service) HandleCommandQueue(ctx context.Context, i *discordgo.Interaction) error {
	args := s.parseArgs(i)
	if len(args) != 0 {
		return stack.Wrap(ctx, s.search(ctx, i, args, false))
	}
	return stack.Wrap(ctx, s.getQueue(ctx, i, false))
}

func (s *service) getQueue(ctx context.Context, i *discordgo.Interaction, edit bool) error {
	queuedVideos, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL, edit))
	}

	videos := make([]entity.Video, len(queuedVideos))
	for i, v := range queuedVideos {
		videos[i] = entity.Video{
			Title:        v.Title,
			ChannelTitle: v.ChannelTitle,
		}
	}

	msg, component := s.template.GetQueue(videos)

	var components []discordgo.MessageComponent
	if component != nil {
		components = append(components, component)
	}

	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{msg}, components, edit))
}
