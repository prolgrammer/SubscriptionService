package usecases

import (
	"context"
	"subscription_service/internal/entities"
	"time"
)

//go:generate mockgen -source=contracts.go --destination=mock_test.go -package=usecases

type CreateSubRepository interface {
	Insert(ctx context.Context, sub *entities.Subscription) error
}

type DeleteSubRepository interface {
	Delete(ctx context.Context, subID string) error
}

type UpdateSubRepository interface {
	Update(ctx context.Context, sub *entities.Subscription) error
}

type GetSubRepository interface {
	SelectByID(ctx context.Context, subID string) (entities.Subscription, error)
}

type GetAllSubsRepository interface {
	SelectAll(ctx context.Context, limit, offset int) ([]entities.Subscription, error)
}

type CalculateTotalCostRepository interface {
	CalculateTotalCost(ctx context.Context, startPeriod, endPeriod time.Time, userID *string, serviceName *string) (int, error)
}
