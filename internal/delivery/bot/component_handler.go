package bot

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/errors/stack"
)

const (
	compAdd   = "add"
	compNext  = "next"
	compPurge = "purge"
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
			stack.Wrap(ctx, b.service.HandleComponentAddSong(ctx, i.Interaction, customIDs[1:]))
		case compNext:
			stack.Wrap(ctx, b.service.HandleComponentNext(ctx, i.Interaction))
		case compPurge:
			stack.Wrap(ctx, b.service.HandleComponentPurge(ctx, i.Interaction))
		}
	}
}
