package repository

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka-kai/internal/domain/discord/entity"
)

// Repository contains functions for discord domain.
type Repository interface {
	Run() error
	Close() error

	AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready))
	AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate))

	SendInteractionMessage(ctx context.Context, i *discordgo.Interaction, data entity.Message) error
	SendMessage(ctx context.Context, channelID string, data entity.Message) (string, string, error)
	EditMessage(ctx context.Context, channelID, messageID string, data entity.Message) error

	JoinVoiceChannel(ctx context.Context, i *discordgo.Interaction) error
	LeaveVoiceChannel(ctx context.Context, i *discordgo.Interaction) error

	InitPlayer(guildID string)
	Stream(ctx context.Context, guildID, path string) error

	SetChannelIDMessageID(guildID string, channelID, messageID string)
	GetChannelIDMessageID(guildID string) (string, string)
	GetMessageID(guildID string) string
	SetAutoNext(guildID string, value bool)
	GetAutoNext(guildID string) bool
	SetQueueIndex(guildID string, value int)
	GetQueueIndex(guildID string) int

	Pause(guildID string)
	Resume(guildID string)
	Stop(guildID string)
}
