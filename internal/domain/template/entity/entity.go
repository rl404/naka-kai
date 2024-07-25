package entity

import "time"

// Video is video entity.
type Video struct {
	ID           string
	ChannelTitle string
	ChannelURL   string
	VideoTitle   string
	VideoURL     string
	Image        string
	Duration     time.Duration
	View         int
	Like         int
	AddedBy      string
	Order        int
	QueueSize    int
}
