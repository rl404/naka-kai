package service

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	queueEntity "github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// HandleComponentAdd to handle component add.
func (s *service) HandleComponentAdd(ctx context.Context, i *discordgo.Interaction, customIDs []string) error {
	// Check queue size.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue, true)), err)
	}

	if len(queue) >= s.queue.GetQueueSize() {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.ReachedQueueLimit(true)))
	}

	// Get video data.
	videoID := customIDs[0]
	play, _ := strconv.ParseBool(customIDs[1])

	video, err := s.youtube.GetVideo(ctx, videoID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidYoutubeURL, true)), err)
	}

	// Add to queue.
	queuedVideo := queueEntity.Queue{
		GuildID:      i.GuildID,
		ChannelTitle: video.ChannelTitle,
		ChannelURL:   s.youtube.GenerateChannelURL(video.ChannelID),
		VideoID:      video.ID,
		VideoTitle:   video.Title,
		VideoURL:     s.youtube.GenerateVideoURL(videoID),
		ImageURL:     video.Image,
		Duration:     video.Duration,
		View:         video.View,
		Like:         video.Like,
		AddedBy:      i.Member.User.ID,
		Order:        len(queue) + 1,
	}

	if err := s.queue.Add(ctx, queuedVideo); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrAddQueue, true)), err)
	}

	if err := s.discord.SendInteractionMessage(ctx, i, s.template.AddedVideo([]entity.Video{{
		ChannelTitle: queuedVideo.ChannelTitle,
		ChannelURL:   queuedVideo.ChannelURL,
		VideoTitle:   queuedVideo.VideoTitle,
		VideoURL:     queuedVideo.VideoURL,
		Image:        queuedVideo.ImageURL,
		Duration:     queuedVideo.Duration,
		View:         queuedVideo.View,
		Like:         queuedVideo.Like,
		AddedBy:      queuedVideo.AddedBy,
		Order:        queuedVideo.Order,
		QueueSize:    len(queue) + 1,
	}}, true)); err != nil {
		return stack.Wrap(ctx, err)
	}

	if play {
		return stack.Wrap(ctx, s.play(ctx, i))
	}

	return nil
}
