package app

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"subscription_service/config"
	"subscription_service/infrastructure/postgres"
	"subscription_service/infrastructure/postgres/commands/subscription"
	http2 "subscription_service/internal/controllers/http"
	"subscription_service/internal/controllers/http/middleware"
	"subscription_service/internal/usecases"
	"subscription_service/pkg/logger"
)

var (
	l              logger.Logger
	postgresClient *postgres.Client

	createSubscriptionUseCase usecases.CreateSubUseCase
	updateSubscriptionUseCase usecases.UpdateSubUseCase
	getSubscriptionUseCase    usecases.GetSubUseCase
	getSubscriptionsUseCase   usecases.GetListSubUseCase
	DeleteSubscriptionUseCase usecases.DeleteSubUseCase
	CalculateTotalCostUseCase usecases.CalculateTotalCostUseCase

	subRepo subscription.SubRepository
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	initPackages(cfg)
	initRepository()
	initUseCases()

	defer postgresClient.Close()
	runHTTP(cfg)
}

func initUseCases() {
	createSubscriptionUseCase = usecases.NewCreateSubUseCase(subRepo, l)
	updateSubscriptionUseCase = usecases.NewUpdateSubUseCase(subRepo, l)
	getSubscriptionUseCase = usecases.NewGetSubUseCase(subRepo, l)
	getSubscriptionsUseCase = usecases.NewGetListSubUseCase(subRepo, l)
	DeleteSubscriptionUseCase = usecases.NewDeleteSubUseCase(subRepo, l)
	CalculateTotalCostUseCase = usecases.NewCalculateTotalCostUseCase(subRepo, l)
}

func initRepository() {
	subRepo = subscription.NewSubRepository(postgresClient, l)
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

	mw := middleware.NewMiddleware(l)

	http2.InitServiceMiddleware(router)
	http2.NewCreateSubController(router, createSubscriptionUseCase, mw, l)
	http2.NewUpdateSubController(router, updateSubscriptionUseCase, mw, l)
	http2.NewGetSubController(router, getSubscriptionUseCase, mw, l)
	http2.NewGetListSubController(router, getSubscriptionsUseCase, mw, l)
	http2.NewDeleteSubController(router, DeleteSubscriptionUseCase, mw, l)
	http2.NewCalculateTotalCostController(router, CalculateTotalCostUseCase, mw, l)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	l.Info().Msgf("starting HTTP server on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
