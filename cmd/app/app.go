package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"subscription_service/config"
	"subscription_service/infrastructure/postgres"
	"subscription_service/pkg/logger"
)

var (
	l              logger.Logger
	postgresClient *postgres.Client
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initPackages(cfg)

	defer postgresClient.Close()
	runHTTP(cfg)
}

func initPackages(cfg *config.Config) {
	var err error

	l = logger.NewConsoleLogger(logger.LevelSwitch(cfg.LogLevel))

	l.Info().Msgf("starting postgres client")
	postgresClient, err = postgres.New(cfg.PG, l)
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

func runHTTP(cfg *config.Config) {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	l.Info().Msgf("starting HTTP server on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
