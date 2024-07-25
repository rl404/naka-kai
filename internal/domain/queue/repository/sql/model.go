package sql

import (
	"time"

	"github.com/rl404/naka-kai/internal/domain/queue/entity"
	"gorm.io/gorm"
)

// Queue is queue table.
type Queue struct {
	ID           int64  `gorm:"primaryKey"`
	GuildID      string `gorm:"index"`
	ChannelTitle string
	ChannelURL   string
	VideoID      string
	VideoTitle   string
	VideoURL     string
	ImageURL     string
	Duration     int64 // time.Duration
	View         int
	Like         int
	AddedBy      string
	Order        int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (q *Queue) toEntity() entity.Queue {
	return entity.Queue{
		ID:           q.ID,
		GuildID:      q.GuildID,
		ChannelTitle: q.ChannelTitle,
		ChannelURL:   q.ChannelURL,
		VideoID:      q.VideoID,
		VideoTitle:   q.VideoTitle,
		VideoURL:     q.VideoURL,
		ImageURL:     q.ImageURL,
		Duration:     time.Duration(q.Duration),
		View:         q.View,
		Like:         q.Like,
		AddedBy:      q.AddedBy,
		Order:        q.Order,
	}
}

func (sql *sql) fromEntity(data entity.Queue) *Queue {
	return &Queue{
		ID:           data.ID,
		GuildID:      data.GuildID,
		ChannelTitle: data.ChannelTitle,
		ChannelURL:   data.ChannelURL,
		VideoID:      data.VideoID,
		VideoTitle:   data.VideoTitle,
		VideoURL:     data.VideoURL,
		ImageURL:     data.ImageURL,
		Duration:     int64(data.Duration),
		View:         data.View,
		Like:         data.Like,
		AddedBy:      data.AddedBy,
		Order:        data.Order,
	}
}
