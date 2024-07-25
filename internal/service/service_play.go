package service

import (
	"context"
	_errors "errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/errors"
)

// HandleCommandPlay to handle command play.
func (s *service) HandleCommandPlay(ctx context.Context, i *discordgo.Interaction) error {
	args := s.parseArgs(i)
	if len(args) > 0 {
		return stack.Wrap(ctx, s.search(ctx, i, args, true))
	}

	if err := s.play(ctx, i); err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Error(entity.ErrPlayingSong)))
	}

	return stack.Wrap(ctx, s.discord.SendInteractionMessage(ctx, i, s.template.Playing()))
}

func (s *service) play(ctx context.Context, i *discordgo.Interaction) error {
	// Player already exist.
	if s.discord.GetMessageID(i.GuildID) != "" {
		return nil
	}

	// Check queue.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		_, _, err2 := s.discord.SendMessage(ctx, i.ChannelID, s.template.Error(entity.ErrGetQueue))
		return stack.Wrap(ctx, err2, err)
	}

	if len(queue) == 0 {
		_, _, err := s.discord.SendMessage(ctx, i.ChannelID, s.template.Queue(nil, 0))
		return stack.Wrap(ctx, err)
	}

	// Join voice channel.
	if err := stack.Wrap(ctx, s.discord.JoinVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			_, _, err2 := s.discord.SendMessage(ctx, i.ChannelID, s.template.Error(entity.NotInVC))
			return stack.Wrap(ctx, err2, err)
		}
		return stack.Wrap(ctx, err)
	}

	s.discord.SetAutoNext(i.GuildID, true)

	// Loop the queue.
	go func() error {
		channelID, messageID, err := s.discord.SendMessage(ctx, i.ChannelID, s.template.Player(entity.Video{
			ChannelTitle: queue[0].ChannelTitle,
			ChannelURL:   queue[0].ChannelURL,
			VideoTitle:   queue[0].VideoTitle,
			VideoURL:     queue[0].VideoURL,
			Image:        queue[0].ImageURL,
			Duration:     queue[0].Duration,
			View:         queue[0].View,
			Like:         queue[0].Like,
			AddedBy:      queue[0].AddedBy,
			Order:        queue[0].Order,
			QueueSize:    len(queue),
		}, true))
		if err != nil {
			return stack.Wrap(ctx, err)
		}

		s.discord.SetChannelIDMessageID(i.GuildID, channelID, messageID)
		defer s.discord.SetChannelIDMessageID(i.GuildID, "", "")

		for {
			// Get queue.
			queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
			if err != nil {
				return stack.Wrap(ctx, s.discord.EditMessage(ctx, channelID, messageID, s.template.Error(entity.ErrGetQueue)), err)
			}

			if len(queue) == 0 {
				return stack.Wrap(ctx, s.discord.EditMessage(ctx, channelID, messageID, s.template.Queue(nil, 0)))
			}

			queueIndex := s.discord.GetQueueIndex(i.GuildID)
			if queueIndex >= len(queue) {
				queueIndex = 0
			}
			if queueIndex < 0 {
				queueIndex = len(queue) - 1
			}
			s.discord.SetQueueIndex(i.GuildID, queueIndex)
			song := queue[queueIndex]

			if err := s.discord.EditMessage(ctx, channelID, messageID, s.template.Player(entity.Video{
				ChannelTitle: song.ChannelTitle,
				ChannelURL:   song.ChannelURL,
				VideoTitle:   song.VideoTitle,
				VideoURL:     song.VideoURL,
				Image:        song.ImageURL,
				Duration:     song.Duration,
				View:         song.View,
				Like:         song.Like,
				AddedBy:      song.AddedBy,
				Order:        song.Order,
				QueueSize:    len(queue),
			}, true)); err != nil {
				return stack.Wrap(ctx, err)
			}

			// Get video source.
			sourceURL, err := s.youtube.GetSourceURLByID(ctx, song.VideoID)
			if err != nil {
				return stack.Wrap(ctx, s.discord.EditMessage(ctx, channelID, messageID, s.template.Error(entity.InvalidYoutubeURL)), err)
			}

			// Start stream.
			if err := s.discord.Stream(ctx, i.GuildID, sourceURL); err != nil {
				return stack.Wrap(ctx, s.discord.EditMessage(ctx, channelID, messageID, s.template.Error(entity.ErrPlayingSong)), err)
			}

			// Go to next queue.
			s.discord.SetQueueIndex(i.GuildID, s.discord.GetQueueIndex(i.GuildID)+1)

			// Stopped.
			if !s.discord.GetAutoNext(i.GuildID) {
				return stack.Wrap(ctx, s.discord.EditMessage(ctx, channelID, messageID, s.template.Stopped()))
			}
		}
	}()

	return nil
}
