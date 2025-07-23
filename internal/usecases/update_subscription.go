package usecases

import (
	"context"
	"subscription_service/internal/controllers/responses"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/entities"
	"subscription_service/pkg/logger"
)

type updateSubUseCase struct {
	subRepo UpdateSubRepository
	logger  logger.Logger
}

type UpdateSubUseCase interface {
	UpdateSubscription(ctx context.Context, subID string, req requests.SubRequest) (responses.SubResponse, error)
}

func NewUpdateSubUseCase(subRepo UpdateSubRepository, logger logger.Logger) UpdateSubUseCase {
	return &updateSubUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (u *updateSubUseCase) UpdateSubscription(ctx context.Context, subID string, req requests.SubRequest) (responses.SubResponse, error) {
	subUUID, err := uuid.Parse(subID)
	if err != nil {
		u.logger.Error().Err(err).Msg("Invalid sub_id format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidUUID, "failed to parse sub_id")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		u.logger.Error().Err(err).Msg("Invalid user_id format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidUUID, "failed to parse user_id")
	}

	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		u.logger.Error().Err(err).Msg("Invalid start_date format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse start_date")
	}

	var endDate *time.Time
	if req.EndDate != "" {
		ed, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			u.logger.Error().Err(err).Msg("Invalid end_date format")
			return responses.SubResponse{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse end_date")
		}
		endDate = &ed
	}

	sub := &entities.Subscription{
		ID:          subUUID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := u.subRepo.Update(ctx, sub); err != nil {
		u.logger.Error().Err(err).Msg("Failed to update subscription")
		return responses.SubResponse{}, errors.Wrap(err, "failed to update subscription")
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
