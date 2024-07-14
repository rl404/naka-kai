package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/domain/queue/repository"
)

type client struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new queue cache.
func New(cacher cache.Cacher, repo repository.Repository) *client {
	return &client{
		cacher: cacher,
		repo:   repo,
	}
}

// Add to add to queue.
func (c *client) Add(ctx context.Context, data entity.Queue) error {
	return c.repo.Add(ctx, data)
}

// GetByGuildID to get by guild id.
func (c *client) GetByGuildID(ctx context.Context, guildID string) ([]entity.Queue, error) {
	return c.repo.GetByGuildID(ctx, guildID)
}

// CountByGuildID to get count by guild id.
func (c *client) CountByGuildID(ctx context.Context, guildID string) (int, error) {
	return c.repo.CountByGuildID(ctx, guildID)
}

// DeleteByGuildID to delete by guild id.
func (c *client) DeleteByGuildID(ctx context.Context, guildID string) error {
	return c.repo.DeleteByGuildID(ctx, guildID)
}

// DeleteByID to delete by id.
func (c *client) DeleteByID(ctx context.Context, id int64) error {
	return c.repo.DeleteByID(ctx, id)
}
