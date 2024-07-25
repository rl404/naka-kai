package service

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	queueEntity "github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

func (s *service) search(ctx context.Context, i *discordgo.Interaction, args []string, play bool) error {
	// If using youtube url.
	if s.youtube.IsURLValid(args[0]) {
		return stack.Wrap(ctx, s.searchByURL(ctx, i, args, play))
	}

	// Check queue size.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue)), err)
	}

	if len(queue) >= s.queue.GetQueueSize() {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.ReachedQueueLimit()))
	}

	// Search song in youtube.
	videos, err := s.youtube.GetVideos(ctx, strings.Join(args, " "), 5)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidSearchQuery)), err)
	}

	videoResults := make([]entity.Video, len(videos))
	for i, v := range videos {
		videoResults[i] = entity.Video{
			ID:           v.ID,
			VideoTitle:   v.Title,
			ChannelTitle: v.ChannelTitle,
		}
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.VideoSearchResults(videoResults, play)))
}

func (s *service) searchByURL(ctx context.Context, i *discordgo.Interaction, args []string, play bool) error {
	// Check queue size.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrGetQueue)), err)
	}

	if len(queue)+len(args) >= s.queue.GetQueueSize() {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.ReachedQueueLimit()))
	}

	// Get video data.
	var queuedVideos []queueEntity.Queue
	for _, arg := range args {
		if !s.youtube.IsURLValid(arg) {
			continue
		}

		videoID, err := s.youtube.GetIDFromURL(ctx, arg)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidYoutubeURL)), err)
		}

		video, err := s.youtube.GetVideo(ctx, videoID)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidYoutubeURL)), err)
		}

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
			Order:        len(queue) + len(queuedVideos) + 1,
		}

		if err := s.queue.Add(ctx, queuedVideo); err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrAddQueue)), err)
		}

		queuedVideos = append(queuedVideos, queuedVideo)
	}

	if len(queuedVideos) == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.InvalidYoutubeURL)))
	}

	addedVideos := make([]entity.Video, len(queuedVideos))
	for i, v := range queuedVideos {
		addedVideos[i] = entity.Video{
			ChannelTitle: v.ChannelTitle,
			ChannelURL:   v.ChannelURL,
			VideoTitle:   v.VideoTitle,
			VideoURL:     v.VideoURL,
			Image:        v.ImageURL,
			Duration:     v.Duration,
			View:         v.View,
			Like:         v.Like,
			AddedBy:      v.AddedBy,
			Order:        v.Order,
			QueueSize:    len(queue) + len(queuedVideos),
		}
	}

	if err := s.discord.SendInteractionMessage(ctx, i, s.template.AddedVideo(addedVideos)); err != nil {
		return stack.Wrap(ctx, err)
	}

	if play {
		return stack.Wrap(ctx, s.play(ctx, i))
	}

	return nil
}
