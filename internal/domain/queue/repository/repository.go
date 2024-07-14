package repository

import (
	"context"

	"github.com/rl404/naka-kai/internal/domain/queue/entity"
)

// Repository contains functions for queue domain.
type Repository interface {
	Add(ctx context.Context, data entity.Queue) error
	GetByGuildID(ctx context.Context, guildID string) ([]entity.Queue, error)
	CountByGuildID(ctx context.Context, guildID string) (int, error)
	DeleteByGuildID(ctx context.Context, guildID string) error
	DeleteByID(ctx context.Context, id int64) error
}
