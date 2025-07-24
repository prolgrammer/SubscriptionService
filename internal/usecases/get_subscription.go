package usecases

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"subscription_service/internal/controllers/responses"
	"subscription_service/pkg/logger"
)

type getSubUseCase struct {
	subRepo GetSubRepository
	logger  logger.Logger
}

type GetSubUseCase interface {
	GetSubscription(c context.Context, subID string) (responses.SubResponse, error)
}

func NewGetSubUseCase(subRepo GetSubRepository, logger logger.Logger) GetSubUseCase {
	return &getSubUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (g *getSubUseCase) GetSubscription(c context.Context, subID string) (responses.SubResponse, error) {
	if _, err := uuid.Parse(subID); err != nil {
		g.logger.Error().Err(err).Msg("Invalid sub_id format")
		return responses.SubResponse{}, errors.Wrap(ErrInvalidUUID, "failed to parse sub_id")
	}

	sub, err := g.subRepo.SelectByID(c, subID)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get subscription")
		return responses.SubResponse{}, errors.Wrap(err, "failed to get subscription")
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
