package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	discordEntity "github.com/rl404/naka-kai/internal/domain/discord/entity"
	"github.com/rl404/naka-kai/internal/domain/template/entity"
	"github.com/rl404/naka-kai/internal/utils"
)

type client struct{}

// New to create new template client.
func New() *client {
	return &client{}
}

// Error is error template.
func (c *client) Error(err string, isEdit ...bool) discordEntity.Message {
	var edit bool
	if len(isEdit) > 0 && isEdit[0] {
		edit = isEdit[0]
	}

	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title:       "Error",
			Description: err,
			Color:       entity.ColorRed,
		}},
		Components: []discordgo.MessageComponent{},
		IsEdit:     edit,
		AutoDelete: true,
	}
}

// Joined is joined template.
func (c *client) Joined() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Joined voice channel.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// Left is left template.
func (c *client) Left() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Left voice channel.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// ReachedQueueLimit is reached queue limit template.
func (c *client) ReachedQueueLimit(isEdit ...bool) discordEntity.Message {
	var edit bool
	if len(isEdit) > 0 && isEdit[0] {
		edit = isEdit[0]
	}

	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title:       "Reached Queue Limit.",
			Description: "Please remove some songs in queue.",
			Color:       entity.ColorGreyLight,
		}},
		Components: []discordgo.MessageComponent{},
		IsEdit:     edit,
	}
}

// AddedVideo is added video template.
func (c *client) AddedVideo(data []entity.Video, isEdit ...bool) discordEntity.Message {
	var edit bool
	if len(isEdit) > 0 && isEdit[0] {
		edit = isEdit[0]
	}

	msgs := make([]*discordgo.MessageEmbed, len(data))
	for i, d := range data {
		msgs[i] = &discordgo.MessageEmbed{
			Title: d.VideoTitle,
			Color: entity.ColorBlue,
			URL:   d.VideoURL,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: d.Image,
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name: d.ChannelTitle,
				URL:  d.ChannelURL,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Added to queue %d/%d", d.Order, d.QueueSize),
			},
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Duration",
					Value:  d.Duration.String(),
					Inline: true,
				},
				{
					Name:   "View",
					Value:  utils.Thousands(d.View),
					Inline: true,
				},
				{
					Name:   "Like",
					Value:  utils.Thousands(d.Like),
					Inline: true,
				},
			},
		}
	}

	return discordEntity.Message{
		Messages:   msgs,
		Components: []discordgo.MessageComponent{},
		IsEdit:     edit,
	}
}

// VideoSearchResults is video search results.
func (c *client) VideoSearchResults(data []entity.Video, play bool) discordEntity.Message {
	if len(data) == 0 {
		return discordEntity.Message{
			Messages: []*discordgo.MessageEmbed{{
				Title: "No song found.",
				Color: entity.ColorGreyLight,
			}},
		}
	}

	fields := make([]*discordgo.MessageEmbedField, len(data))
	buttons := make([]discordgo.MessageComponent, len(data))
	for i, video := range data {
		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("#%d | %s", i+1, video.VideoTitle),
			Value: video.ChannelTitle,
		}
		buttons[i] = discordgo.Button{
			Label:    fmt.Sprintf("#%d", i+1),
			CustomID: fmt.Sprintf("add-%s-%v", video.ID, play),
		}
	}

	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title:  "Search Results",
			Color:  entity.ColorGreyLight,
			Fields: fields,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Click the button below to choose the song.",
			},
		}},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: buttons,
			},
		},
	}
}

// Queue is queue template.
func (c *client) Queue(data []entity.Video, index int) discordEntity.Message {
	if len(data) == 0 {
		return discordEntity.Message{
			Messages: []*discordgo.MessageEmbed{{
				Title:       "Queued Songs",
				Description: "Empty queue. Go add some songs.",
				Color:       entity.ColorGreyLight,
			}},
			Components: []discordgo.MessageComponent{},
		}
	}

	fields := make([]*discordgo.MessageEmbedField, len(data))
	for i, video := range data {
		no := fmt.Sprintf("#%d", video.Order)
		if i == index {
			no = "->"
		}

		fields[i] = &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s | %s", no, video.VideoTitle),
			Value: fmt.Sprintf("%s | *Added by <@%s>*", video.ChannelTitle, video.AddedBy),
		}
	}

	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Color:  entity.ColorBlue,
			Title:  "Queued Songs",
			Fields: fields,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("playing #%d", index+1),
			},
		}},
	}
}

// Stopped is stopped template.
func (c *client) Stopped() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Song stopped.",
			Color: entity.ColorGreyLight,
		}},
		Components: []discordgo.MessageComponent{},
		AutoDelete: true,
	}
}

// Player is player template.
func (c *client) Player(data entity.Video, playing bool) discordEntity.Message {
	var components []discordgo.MessageComponent
	if playing {
		components = c.PlayingButtons().Components
	} else {
		components = c.PausedButtons().Components
	}

	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: data.VideoTitle,
			Color: entity.ColorBlue,
			URL:   data.VideoURL,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: data.Image,
			},
			Author: &discordgo.MessageEmbedAuthor{
				Name: data.ChannelTitle,
				URL:  data.ChannelURL,
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Playing queue %d/%d", data.Order, data.QueueSize),
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
					Name:   "Added by",
					Value:  fmt.Sprintf("<@%s>", data.AddedBy),
					Inline: true,
				},
			},
		}},
		Components: components,
	}
}

// Playing is playing template.
func (c *client) Playing() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Playing song.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// Paused is paused template.
func (c *client) Paused() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Song paused.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// Next is next template.
func (c *client) Next() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Go to next song.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// Previous is previous template.
func (c *client) Previous() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Go to previous.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// Jumped is jumped template.
func (c *client) Jumped(songNumber int) discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: fmt.Sprintf("Go to song #%d.", songNumber),
			Color: entity.ColorGreyLight,
		}},
	}
}

// Removed is removed template.
func (c *client) Removed(songNumbers []int) discordEntity.Message {
	strs := make([]string, len(songNumbers))
	for i, no := range songNumbers {
		strs[i] = fmt.Sprintf("#%s", strconv.Itoa(no))
	}
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: fmt.Sprintf("Removed song %s.", strings.Join(strs, ", ")),
			Color: entity.ColorGreyLight,
		}},
	}
}

// Purged is purged template.
func (c *client) Purged() discordEntity.Message {
	return discordEntity.Message{
		Messages: []*discordgo.MessageEmbed{{
			Title: "Queue purged.",
			Color: entity.ColorGreyLight,
		}},
	}
}

// PausedButtons is paused buttons template.
func (c *client) PausedButtons() discordEntity.Message {
	return discordEntity.Message{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443257983729806",
							Name: "prev",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "previous",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443256058548316",
							Name: "play",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "resume",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443261687431230",
							Name: "stop",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "stop",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443253487439872",
							Name: "next",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "next",
					},
				},
			},
		},
	}
}

// PlayingButtons is playing buttons template.
func (c *client) PlayingButtons() discordEntity.Message {
	return discordEntity.Message{
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443257983729806",
							Name: "prev",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "previous",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443259795669066",
							Name: "pause",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "pause",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443261687431230",
							Name: "stop",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "stop",
					},
					discordgo.Button{
						Emoji: &discordgo.ComponentEmoji{
							ID:   "1264443253487439872",
							Name: "next",
						},
						Style:    discordgo.SecondaryButton,
						CustomID: "next",
					},
				},
			},
		},
	}
}
