package subscription

import (
	"context"
	"subscription_service/infrastructure/postgres"
	"subscription_service/internal/entities"
	"subscription_service/pkg/logger"
	"time"
)

type subRepo struct {
	client *postgres.Client
	logger logger.Logger
}

type SubRepository interface {
	Insert(ctx context.Context, sub *entities.Subscription) error
	Delete(ctx context.Context, subID string) error
	Update(ctx context.Context, sub *entities.Subscription) error
	SelectByID(ctx context.Context, subID string) (entities.Subscription, error)
	SelectAll(ctx context.Context, limit, offset int) ([]entities.Subscription, error)
	CalculateTotalCost(ctx context.Context, startPeriod, endPeriod time.Time, userID *string, serviceName *string) (int, error)
}

func NewSubRepository(client *postgres.Client, logger logger.Logger) SubRepository {
	return &subRepo{
		client: client,
		logger: logger,
	}
}
