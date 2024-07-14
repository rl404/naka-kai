package service

import (
	"context"
	_errors "errors"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/errors"
)

// HandleCommandPlay to handle play command.
func (s *service) HandleCommandPlay(ctx context.Context, i *discordgo.Interaction) error {
	args := s.parseArgs(i)
	if len(args) != 0 {
		return stack.Wrap(ctx, s.search(ctx, i, args, true))
	}

	// Player already exist.
	if s.discord.GetIsPlayerExist(i.GuildID) {
		return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetStillPlaying()}, nil))
	}

	// Check queue.
	queueCount, err := s.queue.CountByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	if queueCount == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetEmptyQueue()}, nil))
	}

	// Join voice channel.
	if err := stack.Wrap(ctx, s.discord.JoinVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.NotInVC), err)
		}
		return stack.Wrap(ctx, err)
	}

	if err := s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetStartPlaying()}, nil); err != nil {
		return stack.Wrap(ctx, err)
	}

	return stack.Wrap(ctx, s.play(ctx, i))
}

func (s *service) play(ctx context.Context, i *discordgo.Interaction) error {
	// Player already exist.
	if s.discord.GetIsPlayerExist(i.GuildID) {
		return nil
	}

	// Check queue.
	queueCount, err := s.queue.CountByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, err)
	}

	if queueCount == 0 {
		return stack.Wrap(ctx, s.discord.SendEmbedMessage(ctx, i.ChannelID, s.template.GetEmptyQueue()))
	}

	// Join voice channel.
	if err := stack.Wrap(ctx, s.discord.JoinVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendErrorMessage(ctx, i.ChannelID, entity.NotInVC), err)
		}
		return stack.Wrap(ctx, err)
	}

	// Loop the queue.
	go func() error {
		s.discord.SetIsPlayerExist(i.GuildID, true)
		defer s.discord.SetIsPlayerExist(i.GuildID, false)

		for {
			// Get queue.
			queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
			if err != nil {
				return stack.Wrap(ctx, s.discord.SendErrorMessage(ctx, i.ChannelID, entity.ErrGetQueue), err)
			}

			if len(queue) == 0 {
				return stack.Wrap(ctx, s.discord.SendEmbedMessage(ctx, i.ChannelID, s.template.GetEndQueue()))
			}

			song := queue[0]

			if err := s.discord.SendEmbedMessage(ctx, i.ChannelID, s.template.GetSong(entity.Video{
				Title:        song.Title,
				URL:          song.URL,
				ChannelTitle: song.ChannelTitle,
				ChannelURL:   song.ChannelURL,
				Image:        song.ImageURL,
				Duration:     song.Duration,
				View:         song.View,
				Like:         song.Like,
				QueueCnt:     len(queue),
			}, true)); err != nil {
				return stack.Wrap(ctx, err)
			}

			// Start stream.
			if err := s.discord.Stream(ctx, i.GuildID, song.SourceURL); err != nil {
				return stack.Wrap(ctx, s.discord.SendErrorMessage(ctx, i.ChannelID, entity.ErrPlaySong), err)
			}

			// Go to next queue.
			if err := s.queue.DeleteByID(ctx, song.ID); err != nil {
				return stack.Wrap(ctx, err)
			}

			if s.discord.GetIsStopped(i.GuildID) {
				return nil
			}
		}
	}()

	return nil
}
