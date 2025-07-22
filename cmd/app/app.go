package app

import (
	"errors"
	"log"
	"subscription_service/config"
	"subscription_service/infrastructure/postgres"
	"subscription_service/pkg/logger"
)

var (
	l logger.Logger
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initPackages(cfg)
}

func initPackages(cfg *config.Config) {
	var err error

	l = logger.NewConsoleLogger(logger.LevelSwitch(cfg.LogLevel))

	l.Info().Msgf("starting postgres client")
	postgresClient, err := postgres.New(cfg.PG, l)
	if err != nil {
		l.Fatal().Msgf("couldn't start postgres: %s", err.Error())
		return
	}

	err = postgresClient.MigrateUp()
	if err != nil {
		if errors.Is(err, postgres.ErrNoChanges) {
			l.Info().Msgf("postgres has the latest version. nothing to migrate")
			return
		}
		l.Fatal().Msgf("couldn't migrate postgres: %s", err.Error())
	}

	l.Info().Msgf("postgres client successfully migrated")
}
