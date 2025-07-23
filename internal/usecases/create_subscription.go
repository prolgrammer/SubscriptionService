package usecases

import (
	"context"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/controllers/responses"
	"subscription_service/internal/entities"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"subscription_service/pkg/logger"
)

type createSubUseCase struct {
	subRepo CreateSubRepository
	logger  logger.Logger
}

type CreateSubUseCase interface {
	CreateSubscription(ctx context.Context, req requests.SubRequest) (responses.SubResponse, error)
}

func NewCreateSubUseCase(subRepo CreateSubRepository, logger logger.Logger) CreateSubUseCase {
	return &createSubUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (c *createSubUseCase) CreateSubscription(ctx context.Context, req requests.SubRequest) (responses.SubResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.logger.Error().Err(err).Msg("Invalid user_id format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidUUID, "failed to parse user_id")
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		c.logger.Error().Err(err).Msg("Invalid start_date format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse start_date")
	}

	var endDate *time.Time
	if req.EndDate != "" {
		ed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			c.logger.Error().Err(err).Msg("Invalid end_date format")
			return responses.SubResponse{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse end_date")
		}
		endDate = &ed
	}

	sub := &entities.Subscription{
		ID:          uuid.New(),
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := c.subRepo.Insert(ctx, sub); err != nil {
		c.logger.Error().Err(err).Msg("Failed to insert subscription")
		return responses.SubResponse{}, errors.Wrap(err, "failed to create subscription")
	}

	response := responses.SubResponse{
		ID:          sub.ID.String(),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate.Format("01-2006"),
	}
	if sub.EndDate != nil {
		endDateStr := sub.EndDate.Format("01-2006")
		response.EndDate = endDateStr
	}

	return response, nil
}
