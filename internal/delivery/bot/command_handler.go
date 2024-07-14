package bot

import (
	"context"

	"github.com/bwmarrin/discordgo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
)

const (
	cmdPlay   = "play"
	cmdQueue  = "queue"
	cmdPurge  = "purge"
	cmdJoin   = "join"
	cmdLeave  = "leave"
	cmdPause  = "pause"
	cmdResume = "resume"
	cmdNext   = "next"
	cmdStop   = "stop"
)

func (b *Bot) registerCommand(s *discordgo.Session) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        cmdPlay,
			Description: "Play queued songs",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "args",
					Description: "Song title or YouTube URLs",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        cmdQueue,
			Description: "Add songs to queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "args",
					Description: "Song title or YouTube URLs",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        cmdPurge,
			Description: "Purge queue",
		},
		{
			Name:        cmdJoin,
			Description: "Tell bot to join your current voice channel",
		},
		{
			Name:        cmdLeave,
			Description: "Tell bot to leave current voice channel",
		},
		{
			Name:        cmdPause,
			Description: "Pause the song",
		},
		{
			Name:        cmdResume,
			Description: "Resume the song",
		},
		{
			Name:        cmdNext,
			Description: "Go to next song",
		},
		{
			Name:        cmdStop,
			Description: "Stop the song",
		},
	}

	for _, cmd := range commands {
		if _, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd); err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) commandHandler(nrApp *newrelic.Application) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionApplicationCommand {
			return
		}

		ctx := stack.Init(context.Background())
		defer b.log(ctx)

		tx := nrApp.StartTransaction("Command " + i.ApplicationCommandData().Name)
		defer tx.End()

		ctx = newrelic.NewContext(ctx, tx)

		b.service.InitPlayer(i.GuildID)

		switch i.ApplicationCommandData().Name {
		case cmdPlay:
			stack.Wrap(ctx, b.service.HandleCommandPlay(ctx, i.Interaction))
		case cmdQueue:
			stack.Wrap(ctx, b.service.HandleCommandQueue(ctx, i.Interaction))
		case cmdPurge:
			stack.Wrap(ctx, b.service.HandleCommandPurge(ctx, i.Interaction))
		case cmdJoin:
			stack.Wrap(ctx, b.service.HandleCommandJoin(ctx, i.Interaction))
		case cmdLeave:
			stack.Wrap(ctx, b.service.HandleCommandLeave(ctx, i.Interaction))
		case cmdPause:
			stack.Wrap(ctx, b.service.HandleCommandPause(ctx, i.Interaction))
		case cmdResume:
			stack.Wrap(ctx, b.service.HandleCommandResume(ctx, i.Interaction))
		case cmdNext:
			stack.Wrap(ctx, b.service.HandleCommandNext(ctx, i.Interaction))
		case cmdStop:
			stack.Wrap(ctx, b.service.HandleCommandStop(ctx, i.Interaction))
		}
	}
}
