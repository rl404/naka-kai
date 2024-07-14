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

	// Search song in youtube.
	videos, err := s.youtube.GetVideos(ctx, strings.Join(args, " "), 5)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidSearchQuery), err)
	}

	videoResults := make([]entity.Video, len(videos))
	for i, v := range videos {
		videoResults[i] = entity.Video{
			ID:           v.ID,
			Title:        v.Title,
			ChannelTitle: v.ChannelTitle,
		}
	}

	msg, component := s.template.GetSearch(videoResults, play)

	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{msg}, []discordgo.MessageComponent{component}))
}

func (s *service) searchByURL(ctx context.Context, i *discordgo.Interaction, args []string, play bool) error {
	var queuedVideos []queueEntity.Queue
	for _, arg := range args {
		if !s.youtube.IsURLValid(arg) {
			continue
		}

		videoID, err := s.youtube.GetIDFromURL(ctx, arg)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL), err)
		}

		sourceURL, err := s.youtube.GetSourceURLByID(ctx, videoID)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL), err)
		}

		video, err := s.youtube.GetVideo(ctx, videoID)
		if err != nil {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL), err)
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
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrAddQueue), err)
		}

		queuedVideos = append(queuedVideos, queuedVideo)
	}

	if len(queuedVideos) == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.InvalidYoutubeURL))
	}

	queueCount, err := s.queue.CountByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrGetQueue), err)
	}

	msgs := make([]*discordgo.MessageEmbed, len(queuedVideos))
	for i, video := range queuedVideos {
		msgs[i] = s.template.GetSong(entity.Video{
			Title:        video.Title,
			URL:          video.URL,
			ChannelTitle: video.ChannelTitle,
			ChannelURL:   video.ChannelURL,
			Image:        video.ImageURL,
			Duration:     video.Duration,
			View:         video.View,
			Like:         video.Like,
			QueueI:       queueCount - len(queuedVideos) + i + 1,
			QueueCnt:     queueCount,
		})
	}

	if err := s.discord.SendInteractionEmbedMessage(ctx, i, msgs, nil); err != nil {
		return stack.Wrap(ctx, err)
	}

	if play {
		return stack.Wrap(ctx, s.play(ctx, i))
	}

	return nil
}
