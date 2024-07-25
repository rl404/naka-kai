package repository

import (
	discordEntity "github.com/rl404/naka-kai/internal/domain/discord/entity"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
)

// Repository contains functions for template.
type Repository interface {
	Error(err string, isEdit ...bool) discordEntity.Message
	Joined() discordEntity.Message
	Left() discordEntity.Message
	ReachedQueueLimit(isEdit ...bool) discordEntity.Message
	AddedVideo(data []entity.Video, isEdit ...bool) discordEntity.Message
	VideoSearchResults(data []entity.Video, play bool) discordEntity.Message
	Queue(data []entity.Video, index int) discordEntity.Message
	Stopped() discordEntity.Message
	Player(data entity.Video, playing bool) discordEntity.Message
	Playing() discordEntity.Message
	Paused() discordEntity.Message
	PausedButtons() discordEntity.Message
	PlayingButtons() discordEntity.Message
	Purged() discordEntity.Message
	Next() discordEntity.Message
	Previous() discordEntity.Message
	Jumped(songNumber int) discordEntity.Message
	Removed(songNumbers []int) discordEntity.Message
}
