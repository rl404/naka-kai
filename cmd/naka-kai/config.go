package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/rl404/fairy/monitoring/newrelic/database"
	"github.com/rl404/naka-kai/internal/errors"
	"github.com/rl404/naka-kai/internal/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type config struct {
	Discord  discordConfig  `envconfig:"DISCORD"`
	DB       dbConfig       `envconfig:"DB"`
	Youtube  youtubeConfig  `envconfig:"YOUTUBE"`
	Log      logConfig      `envconfig:"LOG"`
	Newrelic newrelicConfig `envconfig:"NEWRELIC"`
}

type discordConfig struct {
	Token      string `envconfig:"TOKEN" required:"true"`
	DeleteTime int    `envconfig:"DELETE_TIME"`
	QueueLimit int    `envconfig:"QUEUE_LIMIT" validate:"gte=0" mod:"default=20"`
}

type dbConfig struct {
	Dialect         string        `envconfig:"DIALECT" validate:"required,oneof=mysql postgresql" mod:"default=mysql,no_space,lcase"`
	Address         string        `envconfig:"ADDRESS" validate:"required" mod:"default=localhost:3306,no_space"`
	Name            string        `envconfig:"NAME" validate:"required" mod:"default=akatsuki"`
	User            string        `envconfig:"USER" validate:"required" mod:"default=root"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN" validate:"required,gt=0" mod:"default=10"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE" validate:"required,gt=0" mod:"default=10"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" validate:"required,gt=0" mod:"default=1m"`
}

type youtubeConfig struct {
	Key string `envconfig:"KEY" validate:"required"`
}

type logConfig struct {
	Level utils.LogLevel `envconfig:"LEVEL" default:"-1"`
	JSON  bool           `envconfig:"JSON" default:"false"`
	Color bool           `envconfig:"COLOR" default:"true"`
}

type newrelicConfig struct {
	Name       string `envconfig:"NAME" mod:"default=naka-kai"`
	LicenseKey string `envconfig:"LICENSE_KEY"`
}

const envPath = "../../.env"
const envPrefix = "NAKA_KAI"

func getConfig() (*config, error) {
	var cfg config

	// Load .env file.
	_ = godotenv.Load(envPath)

	// Convert env to struct.
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	// Init global log.
	utils.InitLog(cfg.Log.Level, cfg.Log.JSON, cfg.Log.Color)

	return &cfg, nil
}

func newDB(cfg dbConfig) (*gorm.DB, error) {
	// Split host and port.
	split := strings.Split(cfg.Address, ":")
	if len(split) != 2 {
		return nil, errors.ErrInvalidDBFormat
	}

	var dialector gorm.Dialector
	switch cfg.Dialect {
	case "mysql":
		dialector = mysql.Open(fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local", cfg.User, cfg.Password, cfg.Address, cfg.Name))
	case "postgresql":
		dialector = postgres.Open(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", split[0], split[1], cfg.User, cfg.Password, cfg.Name))
	default:
		panic("invalid db dialect")
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	tmp, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set basic config.
	tmp.SetMaxIdleConns(cfg.MaxConnIdle)
	tmp.SetMaxOpenConns(cfg.MaxConnOpen)
	tmp.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifetime) * time.Second)

	db.Use(database.NewGORM(cfg.Address, cfg.Name))

	return db, nil
}
