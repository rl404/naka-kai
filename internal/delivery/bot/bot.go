package bot

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/naka-kai/internal/service"
	"github.com/rl404/naka-kai/internal/utils"
)

// Bot contains functions for bot.
type Bot struct {
	service service.Service
}

// New to create new bot.
func New(service service.Service) *Bot {
	return &Bot{
		service: service,
	}
}

// Run to run bot.
func (b *Bot) Run() error {
	return b.service.Run()
}

// Stop to stop bot.
func (b *Bot) Stop() error {
	return b.service.Stop()
}

// AddHandler to add handlers.
func (b *Bot) AddHandler(nrApp *newrelic.Application) {
	b.service.AddReadyHandler(b.readyHandler())
	b.service.AddInteractionHandler(b.commandHandler(nrApp))
	b.service.AddInteractionHandler(b.componentHandler(nrApp))
}

func (b *Bot) readyHandler() func(*discordgo.Session, *discordgo.Ready) {
	return func(s *discordgo.Session, _ *discordgo.Ready) {
		// b.registerCommand(s)
	}
}

func (b *Bot) log(ctx context.Context) {
	if rvr := recover(); rvr != nil {
		stack.Wrap(ctx, fmt.Errorf("%s", debug.Stack()), fmt.Errorf("%v", rvr), fmt.Errorf("panic"))
	}

	errStack := stack.Get(ctx)
	if len(errStack) > 0 {
		utils.Log(map[string]interface{}{
			"level": utils.ErrorLevel,
			"error": errStack,
		})
	}
}
