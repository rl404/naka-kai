package service

import (
	"context"

	"github.com/bwmarrin/discordgo"
	discordRepository "github.com/rl404/naka-kai/internal/domain/discord/repository"
	queueRepository "github.com/rl404/naka-kai/internal/domain/queue/repository"
	templateRepository "github.com/rl404/naka-kai/internal/domain/template/repository"
	youtubeRepository "github.com/rl404/naka-kai/internal/domain/youtube/repository"
)

// Service contains functions for service.
type Service interface {
	Run() error
	Stop() error

	AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready))
	AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate))

	InitPlayer(guildID string)

	HandleCommandPlay(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandQueue(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandRemove(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandPurge(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandJoin(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandLeave(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandPause(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandResume(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandPrevious(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandNext(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandJump(ctx context.Context, i *discordgo.Interaction) error
	HandleCommandStop(ctx context.Context, i *discordgo.Interaction) error

	HandleComponentAdd(ctx context.Context, i *discordgo.Interaction, customIDs []string) error
	HandleComponentPause(ctx context.Context, i *discordgo.Interaction) error
	HandleComponentResume(ctx context.Context, i *discordgo.Interaction) error
	HandleComponentPrevious(ctx context.Context, i *discordgo.Interaction) error
	HandleComponentNext(ctx context.Context, i *discordgo.Interaction) error
	HandleComponentStop(ctx context.Context, i *discordgo.Interaction) error
}

type service struct {
	discord  discordRepository.Repository
	youtube  youtubeRepository.Repository
	queue    queueRepository.Repository
	template templateRepository.Repository
}

// New to create new service.
func New(
	discord discordRepository.Repository,
	youtube youtubeRepository.Repository,
	queue queueRepository.Repository,
	template templateRepository.Repository,
) Service {
	return &service{
		discord:  discord,
		youtube:  youtube,
		queue:    queue,
		template: template,
	}
}
