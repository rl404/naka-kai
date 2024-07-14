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
	Title        string
	URL          string
	ChannelTitle string
	ChannelURL   string
	ImageURL     string
	Duration     int64 // time.Duration
	View         int
	Like         int
	SourceURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

func (q *Queue) toEntity() entity.Queue {
	return entity.Queue{
		ID:           q.ID,
		GuildID:      q.GuildID,
		Title:        q.Title,
		URL:          q.URL,
		ChannelTitle: q.ChannelTitle,
		ChannelURL:   q.ChannelURL,
		ImageURL:     q.ImageURL,
		Duration:     time.Duration(q.Duration),
		View:         q.View,
		Like:         q.Like,
		SourceURL:    q.SourceURL,
	}
}

func (sql *sql) fromEntity(data entity.Queue) *Queue {
	return &Queue{
		ID:           data.ID,
		GuildID:      data.GuildID,
		Title:        data.Title,
		URL:          data.URL,
		ChannelTitle: data.ChannelTitle,
		ChannelURL:   data.ChannelURL,
		ImageURL:     data.ImageURL,
		Duration:     int64(data.Duration),
		View:         data.View,
		Like:         data.Like,
		SourceURL:    data.SourceURL,
	}
}
