package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	_bot "github.com/rl404/naka-kai/internal/delivery/bot"
	discordRepository "github.com/rl404/naka-kai/internal/domain/discord/repository"
	discordClient "github.com/rl404/naka-kai/internal/domain/discord/repository/client"
	queueRepository "github.com/rl404/naka-kai/internal/domain/queue/repository"
	queueSQL "github.com/rl404/naka-kai/internal/domain/queue/repository/sql"
	templateRepository "github.com/rl404/naka-kai/internal/domain/template/repository"
	templateClient "github.com/rl404/naka-kai/internal/domain/template/repository/client"
	youtubeRepository "github.com/rl404/naka-kai/internal/domain/youtube/repository"
	youtubeClient "github.com/rl404/naka-kai/internal/domain/youtube/repository/client"
	"github.com/rl404/naka-kai/internal/service"
	"github.com/rl404/naka-kai/internal/utils"
)

func bot() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.LogLevel(cfg.Log.Level)))
		utils.Info("newrelic initialized")
	}

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Init discord.
	var discord discordRepository.Repository
	discord, err = discordClient.New(cfg.Discord.Token, cfg.Discord.DeleteTime)
	if err != nil {
		return err
	}
	utils.Info("discord initialized")

	// Init youtube.
	var youtube youtubeRepository.Repository
	youtube, err = youtubeClient.New(cfg.Youtube.Key)
	if err != nil {
		return err
	}
	utils.Info("youtube initialized")

	// Init queue.
	var queue queueRepository.Repository = queueSQL.New(db, cfg.Discord.QueueLimit)
	utils.Info("queue initialized")

	// Init template.
	var template templateRepository.Repository = templateClient.New()
	utils.Info("template initialized")

	// Init service.
	service := service.New(discord, youtube, queue, template)
	utils.Info("service initialized")

	// Init bot.
	bot := _bot.New(service)
	bot.AddHandler(nrApp)
	utils.Info("bot initialized")

	// Run bot.
	if err := bot.Run(); err != nil {
		return err
	}
	utils.Info("naka-kai is running...")
	defer bot.Stop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	return nil
}
