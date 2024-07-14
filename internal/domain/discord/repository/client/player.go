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

	voice            *discordgo.VoiceConnection
	channelID        string
	isInVoiceChannel bool

	isPlayerExist bool
	isPlaying     bool
	isPaused      bool
	isStopped     bool

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

// GetIsPlayerExist to get isPlayerExist.
func (c *client) GetIsPlayerExist(guildID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].isPlayerExist
}

// SetIsPlayerExist to set isPlayerExist.
func (c *client) SetIsPlayerExist(guildID string, value bool) {
	c.Lock()
	defer c.Unlock()
	c.players[guildID].isPlayerExist = value
}

// Pause to pause song.
func (c *client) Pause(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		!c.players[guildID].isPlayerExist ||
		!c.players[guildID].isPlaying ||
		c.players[guildID].isPaused ||
		c.players[guildID].isStopped {
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
		!c.players[guildID].isPlayerExist ||
		!c.players[guildID].isPlaying ||
		!c.players[guildID].isPaused ||
		c.players[guildID].isStopped {
		return
	}

	c.players[guildID].isPaused = false
	c.players[guildID].streamSession.SetPaused(false)
}

// Next to go next song.
func (c *client) Next(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		!c.players[guildID].isPlayerExist ||
		!c.players[guildID].isPlaying {
		return
	}

	if c.players[guildID].encodeSession != nil {
		c.players[guildID].encodeSession.Cleanup()
	}
}

// Stop to stop song.
func (c *client) Stop(guildID string) {
	c.Lock()
	defer c.Unlock()

	if !c.players[guildID].isInVoiceChannel ||
		!c.players[guildID].isPlayerExist ||
		!c.players[guildID].isPlaying {
		return
	}

	c.players[guildID].isStopped = true
	if c.players[guildID].encodeSession != nil {
		c.players[guildID].encodeSession.Cleanup()
	}
}

// GetIsStopped to get isStopped.
func (c *client) GetIsStopped(guildID string) bool {
	c.Lock()
	defer c.Unlock()
	return c.players[guildID].isStopped
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
	// options.Bitrate = 96
	// options.Application = "lowdelay"

	c.players[guildID].encodeSession, err = dca.EncodeFile(path, options)
	if err != nil {
		return stack.Wrap(ctx, err)
	}
	defer c.players[guildID].encodeSession.Cleanup()

	done := make(chan error)

	c.players[guildID].streamSession = dca.NewStream(c.players[guildID].encodeSession, c.players[guildID].voice, done)

	if err = <-done; err != nil && err != io.EOF {
		return stack.Wrap(ctx, err)
	}

	return nil
}
