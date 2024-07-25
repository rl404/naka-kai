package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
)

const (
	compAdd    = "add"
	compPause  = "pause"
	compResume = "resume"
	compStop   = "stop"
	compPrev   = "previous"
	compNext   = "next"
)

func (b *Bot) componentHandler(nrApp *newrelic.Application) func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type != discordgo.InteractionMessageComponent {
			return
		}

		ctx := stack.Init(context.Background())
		defer b.log(ctx)

		tx := nrApp.StartTransaction("Song Button")
		defer tx.End()

		ctx = newrelic.NewContext(ctx, tx)

		b.service.InitPlayer(i.GuildID)

		customIDs := strings.Split(i.MessageComponentData().CustomID, "-")

		switch customIDs[0] {
		case compAdd:
			stack.Wrap(ctx, b.service.HandleComponentAdd(ctx, i.Interaction, customIDs[1:]))
		case compPause:
			stack.Wrap(ctx, b.service.HandleComponentPause(ctx, i.Interaction))
		case compResume:
			stack.Wrap(ctx, b.service.HandleComponentResume(ctx, i.Interaction))
		case compPrev:
			stack.Wrap(ctx, b.service.HandleComponentPrevious(ctx, i.Interaction))
		case compNext:
			stack.Wrap(ctx, b.service.HandleComponentNext(ctx, i.Interaction))
		case compStop:
			stack.Wrap(ctx, b.service.HandleComponentStop(ctx, i.Interaction))
		}
	}
}
