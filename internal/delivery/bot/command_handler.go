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
	cmdRemove = "remove"
	cmdPurge  = "purge"
	cmdJoin   = "join"
	cmdLeave  = "leave"
	cmdPause  = "pause"
	cmdResume = "resume"
	cmdPrev   = "previous"
	cmdNext   = "next"
	cmdJump   = "jump"
	cmdStop   = "stop"
)

func (b *Bot) registerCommand(s *discordgo.Session) error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        cmdPlay,
			Description: "Play queued songs",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "song_title_or_youtube_urls",
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
					Name:        "song_title_or_youtube_urls",
					Description: "Song title or YouTube URLs",
					Type:        discordgo.ApplicationCommandOptionString,
				},
			},
		},
		{
			Name:        cmdRemove,
			Description: "Remove songs to queue",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "song_numbers",
					Description: "Queued song numbers",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
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
			Name:        cmdPrev,
			Description: "Go to previous song",
		},
		{
			Name:        cmdNext,
			Description: "Go to next song",
		},
		{
			Name:        cmdJump,
			Description: "Jump to desired song number",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "song_number",
					Description: "Queued song number",
					Type:        discordgo.ApplicationCommandOptionString,
					Required:    true,
				},
			},
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
		case cmdRemove:
			stack.Wrap(ctx, b.service.HandleCommandRemove(ctx, i.Interaction))
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
		case cmdPrev:
			stack.Wrap(ctx, b.service.HandleCommandPrevious(ctx, i.Interaction))
		case cmdNext:
			stack.Wrap(ctx, b.service.HandleCommandNext(ctx, i.Interaction))
		case cmdJump:
			stack.Wrap(ctx, b.service.HandleCommandJump(ctx, i.Interaction))
		case cmdStop:
			stack.Wrap(ctx, b.service.HandleCommandStop(ctx, i.Interaction))
		}
	}
}
