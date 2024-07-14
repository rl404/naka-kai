package client

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/utils"
)

type client struct{}

// New to create new template client.
func New() *client {
	return &client{}
}

// GetSong to get song template.
func (c *client) GetSong(data entity.Video, play ...bool) *discordgo.MessageEmbed {
	footer := fmt.Sprintf("Added to queue %d/%d", data.QueueI, data.QueueCnt)
	if len(play) > 0 && play[0] {
		footer = fmt.Sprintf("Playing queue 1/%d", data.QueueCnt)
	}

	return &discordgo.MessageEmbed{
		Title: data.Title,
		Color: entity.ColorBlue,
		URL:   data.URL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: data.Image,
		},
		Author: &discordgo.MessageEmbedAuthor{
			Name: data.ChannelTitle,
			URL:  data.ChannelURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: footer,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Duration",
				Value:  data.Duration.String(),
				Inline: true,
			},
			{
				Name:   "View",
				Value:  utils.Thousands(data.View),
				Inline: true,
			},
			{
				Name:   "Like",
				Value:  utils.Thousands(data.Like),
				Inline: true,
			},
		},
	}
}

// GetSearch to get search template.
func (c *client) GetSearch(data []entity.Video, play bool) (*discordgo.MessageEmbed, discordgo.MessageComponent) {
	fields := make([]*discordgo.MessageEmbedField, len(data))
	buttons := make([]discordgo.MessageComponent, len(data))
	for i, video := range data {
		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("#%d | %s", i+1, video.Title),
			Value: video.ChannelTitle,
		}
		buttons[i] = discordgo.Button{
			Label:    fmt.Sprintf("#%d", i+1),
			CustomID: fmt.Sprintf("add-%s-%v", video.ID, play),
		}
	}

	return &discordgo.MessageEmbed{
			Color: entity.ColorGreyLight,
			Title: "Search Results",
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Click the button below to choose the song.",
			},
			Fields: fields,
		},
		discordgo.ActionsRow{
			Components: buttons,
		}
}

// GetQueue to get queue template.
func (c *client) GetQueue(data []entity.Video) (*discordgo.MessageEmbed, discordgo.MessageComponent) {
	cnt := len(data)
	if cnt == 0 {
		return &discordgo.MessageEmbed{
			Color:       entity.ColorGreyLight,
			Title:       "Queued Songs",
			Description: "Empty queue. Go add some songs.",
		}, nil
	}

	if len(data) > 25 {
		data = data[:25]
	}

	fields := make([]*discordgo.MessageEmbedField, len(data))
	for i, video := range data {
		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("#%d | %s", i+1, video.Title),
			Value: video.ChannelTitle,
		}
	}

	return &discordgo.MessageEmbed{
			Color: entity.ColorGreyLight,
			Title: "Queued Songs",
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("queue 1/%d", cnt),
			},
			Fields: fields,
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Next",
					CustomID: "next",
				},
				discordgo.Button{
					Label:    "Purge",
					CustomID: "purge",
					Style:    discordgo.DangerButton,
				},
			},
		}
}

// GetPurged to get purge template.
func (c *client) GetPurged() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Queue has been purged.",
	}
}

// GetJoined to get join template.
func (c *client) GetJoined() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Joined voice channel.",
	}
}

// GetLeft to get leave template.
func (c *client) GetLeft() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Left voice channel.",
	}
}

// GetPaused to get pause template.
func (c *client) GetPaused() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Song paused.",
	}
}

// GetResumed to get resume template.
func (c *client) GetResumed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Song resumed.",
	}
}

// GetStopped to get stop template.
func (c *client) GetStopped() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Song stopped.",
	}
}

// GetEmptyQueue to get empty queue template.
func (c *client) GetEmptyQueue() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Empty queue. Go add some songs.",
	}
}

// GetEndQueue to get end of queue template.
func (c *client) GetEndQueue() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "End of queue.",
	}
}

// GetStillPlaying to get still playing template.
func (c *client) GetStillPlaying() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Still playing a song.",
	}
}

// GetStartPlaying to get start playing template.
func (c *client) GetStartPlaying() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Start playing a song.",
	}
}

// GetNext to get next template.
func (c *client) GetNext() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Color: entity.ColorGreyLight,
		Title: "Go to next song.",
	}
}
