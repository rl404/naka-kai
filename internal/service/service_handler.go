package service

import (
	"context"
	_errors "errors"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/errors"
)

// AddReadyHandler to add discord bot ready handler.
func (s *service) AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready)) {
	s.discord.AddReadyHandler(fn)
}

// AddInteractionHandler to add discord bot interaction handler.
func (s *service) AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate)) {
	s.discord.AddInteractionHandler(fn)
}

// HandleCommandJoin to handle join command.
func (s *service) HandleCommandJoin(ctx context.Context, i *discordgo.Interaction) error {
	if err := stack.Wrap(ctx, s.discord.JoinVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.NotInVC), err)
		}
		return stack.Wrap(ctx, err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetJoined()}, nil))
}

// HandleCommandLeave to handle leave command.
func (s *service) HandleCommandLeave(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Stop(i.GuildID)
	if err := stack.Wrap(ctx, s.discord.LeaveVoiceChannel(ctx, i)); err != nil {
		if _errors.Is(err, errors.ErrNotInVC) {
			return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.NotInVC), err)
		}
		return stack.Wrap(ctx, err)
	}
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetLeft()}, nil))
}

// HandleCommandPause to handle pause command.
func (s *service) HandleCommandPause(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Pause(i.GuildID)
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetPaused()}, nil))
}

// HandleCommandResume to handle resume command.
func (s *service) HandleCommandResume(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Resume(i.GuildID)
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetResumed()}, nil))
}

// HandleCommandNext to handle next command.
func (s *service) HandleCommandNext(ctx context.Context, i *discordgo.Interaction) error {
	// Get queue.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrGetQueue), err)
	}

	if len(queue) == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetEndQueue()}, nil))
	}

	// Go to next queue.
	if err := s.queue.DeleteByID(ctx, queue[0].ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	s.discord.Next(i.GuildID)

	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetNext()}, nil))
}

// HandleComponentNext to handle next component.
func (s *service) HandleComponentNext(ctx context.Context, i *discordgo.Interaction) error {
	// Get queue.
	queue, err := s.queue.GetByGuildID(ctx, i.GuildID)
	if err != nil {
		return stack.Wrap(ctx, s.discord.SendInteractionErrorMessage(ctx, i, entity.ErrGetQueue, true), err)
	}

	if len(queue) == 0 {
		return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetEndQueue()}, nil, true))
	}

	// Go to next queue.
	if err := s.queue.DeleteByID(ctx, queue[0].ID); err != nil {
		return stack.Wrap(ctx, err)
	}

	s.discord.Next(i.GuildID)

	return stack.Wrap(ctx, s.getQueue(ctx, i, true))
}

// HandleCommandStop to handle stop command.
func (s *service) HandleCommandStop(ctx context.Context, i *discordgo.Interaction) error {
	s.discord.Stop(i.GuildID)
	return stack.Wrap(ctx, s.discord.SendInteractionEmbedMessage(ctx, i, []*discordgo.MessageEmbed{s.template.GetStopped()}, nil))
}

func (s *service) parseArgs(i *discordgo.Interaction) []string {
	var args []string
	opts := i.ApplicationCommandData().Options
	if len(opts) > 0 {
		args = strings.Fields(opts[0].StringValue())
	}
	return args
}
