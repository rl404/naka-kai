package client

import (
	"sync"

	"github.com/bwmarrin/discordgo"
)

type client struct {
	sync.Mutex

	session *discordgo.Session
	players map[string]*player
}

// New to create new discord client.
func New(token string) (*client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &client{
		session: session,
		players: make(map[string]*player),
	}, nil
}

// Run to login and start discord bot.
func (c *client) Run() error {
	return c.session.Open()
}

// Close to stop discord bot.
func (c *client) Close() error {
	return c.session.Close()
}

// AddReadyHandler to add discord bot ready handler.
func (c *client) AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready)) {
	c.session.AddHandler(fn)
}

// AddInteractionHandler to add discord bot interaction handler.
func (c *client) AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate)) {
	c.session.AddHandler(fn)
}
