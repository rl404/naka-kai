package main

import (
	"github.com/rl404/naka-kai/internal/domain/queue/repository/sql"
	"github.com/rl404/naka-kai/internal/utils"
)

func migrate() error {
	// Get config.
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	tmp, _ := db.DB()
	defer tmp.Close()

	// Migrate.
	utils.Info("migrating...")
	if err := db.AutoMigrate(
		sql.Queue{},
	); err != nil {
		return err
	}

	utils.Info("done")
	return nil
}
