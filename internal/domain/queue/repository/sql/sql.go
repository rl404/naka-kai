package sql

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/errors"
	"gorm.io/gorm"
)

type sql struct {
	db        *gorm.DB
	queueSize int
}

// New to create new queue sql database.
func New(db *gorm.DB, queueSize int) *sql {
	return &sql{
		db:        db,
		queueSize: queueSize,
	}
}

// GetQueueSize to get queue size.
func (sql *sql) GetQueueSize() int {
	return sql.queueSize
}

// GetByGuildID to get by guild id.
func (sql *sql) GetByGuildID(ctx context.Context, guildID string) ([]entity.Queue, error) {
	var queue []Queue
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Order(`"order" asc`).Find(&queue).Error; err != nil {
		return nil, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	res := make([]entity.Queue, len(queue))
	for i, q := range queue {
		res[i] = q.toEntity()
	}
	return res, nil
}

// Add to add.
func (sql *sql) Add(ctx context.Context, data entity.Queue) error {
	if err := sql.db.WithContext(ctx).Create(sql.fromEntity(data)).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return nil
}

// DeleteByGuildID to delete by guild id.
func (sql *sql) DeleteByGuildID(ctx context.Context, guildID string) error {
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Delete(&Queue{}).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return nil
}

// DeleteByGuildIDAndOrders to delete by guild id and orders.
func (sql *sql) DeleteByGuildIDAndOrders(ctx context.Context, guildID string, orders []int) error {
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Where(`"order" in ?`, orders).Delete(&Queue{}).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	// Reorder.
	var queue []Queue
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Order(`"order" asc`).Find(&queue).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	for i := range queue {
		queue[i].Order = i + 1
	}

	if err := sql.db.WithContext(ctx).Save(&queue).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return nil
}
