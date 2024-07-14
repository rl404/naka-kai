package repository

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// Repository contains functions for template.
type Repository interface {
	GetSong(data entity.Video, play ...bool) *discordgo.MessageEmbed
	GetSearch(data []entity.Video, play bool) (*discordgo.MessageEmbed, discordgo.MessageComponent)
	GetQueue(data []entity.Video) (*discordgo.MessageEmbed, discordgo.MessageComponent)
	GetPurged() *discordgo.MessageEmbed
	GetJoined() *discordgo.MessageEmbed
	GetLeft() *discordgo.MessageEmbed
	GetPaused() *discordgo.MessageEmbed
	GetResumed() *discordgo.MessageEmbed
	GetStopped() *discordgo.MessageEmbed
	GetEmptyQueue() *discordgo.MessageEmbed
	GetEndQueue() *discordgo.MessageEmbed
	GetStillPlaying() *discordgo.MessageEmbed
	GetStartPlaying() *discordgo.MessageEmbed
	GetNext() *discordgo.MessageEmbed
	// GetJumped(i int) string
	// GetRemoved(i []string) string
}
