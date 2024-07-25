package repository

import (
	"context"

	"github.com/rl404/naka-kai/internal/domain/queue/entity"
)

// Repository contains functions for queue domain.
type Repository interface {
	GetQueueSize() int
	GetByGuildID(ctx context.Context, guildID string) ([]entity.Queue, error)
	Add(ctx context.Context, data entity.Queue) error
	DeleteByGuildID(ctx context.Context, guildID string) error
	DeleteByGuildIDAndOrders(ctx context.Context, guildID string, orders []int) error
}
