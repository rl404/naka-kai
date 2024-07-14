package entity

import "time"

// Video is video entity.
type Video struct {
	ID           string
	Title        string
	URL          string
	ChannelTitle string
	ChannelURL   string
	Image        string
	Duration     time.Duration
	View         int
	Like         int

	QueueI   int
	QueueCnt int
}
