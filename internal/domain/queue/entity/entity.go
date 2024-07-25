package entity

import (
	"time"
)

// Queue is queue entity.
type Queue struct {
	ID           int64
	GuildID      string
	ChannelTitle string
	ChannelURL   string
	VideoID      string
	VideoTitle   string
	VideoURL     string
	ImageURL     string
	Duration     time.Duration
	View         int
	Like         int
	AddedBy      string
	Order        int
}
