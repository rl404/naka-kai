package client

import (
	"context"
	"io"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rl404/fairy/errors/stack"
)

type player struct {
	sync.Mutex

	voice      *discordgo.VoiceConnection
	channelID  string
	messageID  string
	queueIndex int

	isInVoiceChannel bool
	isPlaying        bool
	isPaused         bool
	autoNext         bool

	encodeSession *dca.EncodeSession
	streamSession *dca.StreamingSession
}

// InitPlayer to init player.
func (c *client) InitPlayer(guildID string) {
	c.Lock()
	defer c.Unlock()
	if c.players[guildID] == nil {
		c.players[guildID] = new(player)
	}
}

// GetChannelIDMessageID to get channelID & messageID.
func (c *client) GetChannelIDMessageID(guildID string) (string, string) {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].channelID, c.players[guildID].messageID
}

// SetChannelIDMessageID to set messageID.
func (c *client) SetChannelIDMessageID(guildID string, channelID, messageID string) {
	c.Lock()
	defer c.Unlock()
	c.players[guildID].channelID, c.players[guildID].messageID = channelID, messageID
}

// GetMessageID to get messageID.
func (c *client) GetMessageID(guildID string) string {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].messageID
}

// GetQueueIndex to get queueIndex.
func (c *client) GetQueueIndex(guildID string) int {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].queueIndex
}

// SetQueueIndex to set queueIndex.
func (c *client) SetQueueIndex(guildID string, value int) {
	c.Lock()
	defer c.Unlock()
	c.players[guildID].queueIndex = value
}

// GetAutoNext to get autoNext.
func (c *client) GetAutoNext(guildID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].autoNext
}

// SetAutoNext to set autoNext.
func (c *client) SetAutoNext(guildID string, value bool) {
	c.Lock()
	defer c.Unlock()
	c.players[guildID].autoNext = value
}

// Pause to pause song.
func (c *client) Pause(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		c.players[guildID].messageID == "" ||
		!c.players[guildID].isPlaying ||
		c.players[guildID].isPaused {
		return
	}

	c.players[guildID].isPaused = true
	c.players[guildID].streamSession.SetPaused(true)
}

// Resume to resume song.
func (c *client) Resume(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		c.players[guildID].messageID == "" ||
		!c.players[guildID].isPlaying ||
		!c.players[guildID].isPaused {
		return
	}

	c.players[guildID].isPaused = false
	c.players[guildID].streamSession.SetPaused(false)
}

// Stop to stop song.
func (c *client) Stop(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		c.players[guildID].messageID == "" ||
		!c.players[guildID].isPlaying {
		return
	}

	c.cleanUp(guildID)
}

func (c *client) setIsPlaying(guildID string, value bool) {
	c.Lock()
	defer c.Unlock()
	c.players[guildID].isPlaying = value
}

// Stream to stream a song.
func (c *client) Stream(ctx context.Context, guildID, path string) (err error) {
	c.setIsPlaying(guildID, true)
	defer c.setIsPlaying(guildID, false)

	options := dca.StdEncodeOptions
	options.RawOutput = true

	c.players[guildID].encodeSession, err = dca.EncodeFile(path, options)
	if err != nil {
		return stack.Wrap(ctx, err)
	}
	defer c.cleanUp(guildID)

	done := make(chan error)

	c.players[guildID].streamSession = dca.NewStream(c.players[guildID].encodeSession, c.players[guildID].voice, done)

	if err = <-done; err != nil && err != io.EOF {
		return stack.Wrap(ctx, err)
	}

	return nil
}

func (c *client) cleanUp(guildID string) {
	if c.players[guildID].encodeSession != nil {
		c.players[guildID].encodeSession.Cleanup()
	}
}
