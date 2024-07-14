package service

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	queueEntity "github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleComponentAddSong to handle add song component.
func (s *service) HandleComponentAddSong(ctx context.Context, i *discordgo.Interaction, customIDs []string) error {
	videoID := customIDs[0]
	play, _ := strconv.ParseBool(customIDs[1])

	sourceURL, err := s.youtube.GetSourceURLByID(ctx, videoID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL, true), err)
	}

	video, err := s.youtube.GetVideo(ctx, videoID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL, true), err)
	}

	queuedVideo := queueEntity.Queue{
		GuildID:      i.GuildID,
		Title:        video.Title,
		URL:          s.youtube.GenerateVideoURL(videoID),
		ChannelTitle: video.ChannelTitle,
		ChannelURL:   s.youtube.GenerateChannelURL(video.ChannelID),
		ImageURL:     video.Image,
		Duration:     video.Duration,
		View:         video.View,
		Like:         video.Like,
		SourceURL:    sourceURL,
	}

	if err := s.queue.Add(ctx, queuedVideo); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrAddQueue, true), err)
	}

	queueCount, err := s.queue.CountByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrGetQueue, true), err)
	}

	msg := s.template.GetSong(entity.Video{
		Title:        queuedVideo.Title,
		URL:          queuedVideo.URL,
		ChannelTitle: queuedVideo.ChannelTitle,
		ChannelURL:   queuedVideo.ChannelURL,
		Image:        queuedVideo.ImageURL,
		Duration:     queuedVideo.Duration,
		View:         queuedVideo.View,
		Like:         queuedVideo.Like,
		QueueI:       queueCount,
		QueueCnt:     queueCount,
	})

	if err := s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{msg}, nil, true); err != nil {
		return stack.Wrap(ctx, err)
	}

	if play {
		return stack.Wrap(ctx, s.play(ctx, i))
	}

	return nil
}
