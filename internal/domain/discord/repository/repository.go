package repository

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Repository contains functions for discord domain.
type Repository interface {
	Run() error
	Close() error

	AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready))
	AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate))

	SendInteractionEmbedMessage(ctx context.Context, i *discordgo.Interaction, msgs []*discordgo.MessageEmbed, components []discordgo.MessageComponent, isEdit ...bool) error
	SendInteractionErrorMessage(ctx context.Context, i *discordgo.Interaction, msg string, isEdit ...bool) error
	SendEmbedMessage(ctx context.Context, channelID string, msg *discordgo.MessageEmbed) error
	SendErrorMessage(ctx context.Context, channelID string, msg string) error

	InitPlayer(guildID string)
	Stream(ctx context.Context, guildID, path string) error

	SetIsPlayerExist(guildID string, value bool)
	GetIsPlayerExist(guildID string) bool
	GetIsStopped(guildID string) bool

	Pause(guildID string)
	Resume(guildID string)
	Next(guildID string)
	Stop(guildID string)

	JoinVoiceChannel(ctx context.Context, i *discordgo.Interaction) error
	LeaveVoiceChannel(ctx context.Context, i *discordgo.Interaction) error
}
