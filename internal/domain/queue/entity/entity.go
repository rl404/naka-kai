package entity

import "time"

// Queue is queue entity.
type Queue struct {
	ID           int64
	GuildID      string
	Title        string
	URL          string
	ChannelTitle string
	ChannelURL   string
	ImageURL     string
	Duration     time.Duration
	View         int
	Like         int
	SourceURL    string
}
