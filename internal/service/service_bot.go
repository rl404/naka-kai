package service

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Run to run discord bot.
func (s *service) Run() error {
	return s.discord.Run()
}

// Stop to stop discord bot.
func (s *service) Stop() error {
	return s.discord.Close()
}

// AddReadyHandler to add discord bot ready handler.
func (s *service) AddReadyHandler(fn func(*discordgo.Session, *discordgo.Ready)) {
	s.discord.AddReadyHandler(fn)
}

// AddInteractionHandler to add discord bot interaction handler.
func (s *service) AddInteractionHandler(fn func(*discordgo.Session, *discordgo.InteractionCreate)) {
	s.discord.AddInteractionHandler(fn)
}

// InitPlayer to init player.
func (s *service) InitPlayer(guildID string) {
	s.discord.InitPlayer(guildID)
}

func (s *service) parseArgs(i *discordgo.Interaction) []string {
	var args []string
	opts := i.ApplicationCommandData().Options
	if len(opts) > 0 {
		args = strings.Fields(opts[0].StringValue())
	}
	return args
}
