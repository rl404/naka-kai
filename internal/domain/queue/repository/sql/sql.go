package sql

import (
	"context"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/domain/queue/entity"
	"github.com/rl404/naka-kai/internal/errors"
	"gorm.io/gorm"
)

type sql struct {
	db *gorm.DB
}

// New to create new queue sql database.
func New(db *gorm.DB) *sql {
	return &sql{
		db: db,
	}
}

// Add to add to queue.
func (sql *sql) Add(ctx context.Context, data entity.Queue) error {
	if err := sql.db.WithContext(ctx).Create(sql.fromEntity(data)).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return nil
}

// GetByGuildID to get by guild id.
func (sql *sql) GetByGuildID(ctx context.Context, guildID string) ([]entity.Queue, error) {
	var data []Queue
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Order("id asc").Find(&data).Error; err != nil {
		return nil, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	res := make([]entity.Queue, len(data))
	for i, q := range data {
		res[i] = q.toEntity()
	}
	return res, nil
}

// CountByGuildID to get count by guild id.
func (sql *sql) CountByGuildID(ctx context.Context, guildID string) (int, error) {
	var cnt int64
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Model(&Queue{}).Count(&cnt).Error; err != nil {
		return 0, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return int(cnt), nil
}

// DeleteByGuildID to delete by guild id.
func (sql *sql) DeleteByGuildID(ctx context.Context, guildID string) error {
	if err := sql.db.WithContext(ctx).Where("guild_id = ?", guildID).Delete(&Queue{}).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return nil
}

// DeleteByID to delete by id.
func (sql *sql) DeleteByID(ctx context.Context, id int64) error {
	if err := sql.db.WithContext(ctx).Where("id = ?", id).Delete(&Queue{}).Error; err != nil {
		return stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return nil
}
