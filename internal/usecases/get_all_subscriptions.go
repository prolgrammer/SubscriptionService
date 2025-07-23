package usecases

import (
	"context"
	"subscription_service/internal/controllers/responses"

	"github.com/pkg/errors"
	"subscription_service/pkg/logger"
)

type GetListSubUseCase interface {
	GetListSubscriptions(ctx context.Context, limit, offset int) ([]responses.SubResponse, error)
}

type getListSubUseCase struct {
	subRepo GetAllSubsRepository
	logger  logger.Logger
}

func NewGetListSubUseCase(subRepo GetAllSubsRepository, logger logger.Logger) GetListSubUseCase {
	return &getListSubUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (g *getListSubUseCase) GetListSubscriptions(ctx context.Context, limit, offset int) ([]responses.SubResponse, error) {
	subs, err := g.subRepo.SelectAll(ctx, limit, offset)
	if err != nil {
		g.logger.Error().Err(err).Msg("Failed to get subscriptions")
		return nil, errors.Wrap(err, "failed to get subscriptions")
	}

	response := make([]responses.SubResponse, 0, len(subs))
	for _, sub := range subs {
		resp := responses.SubResponse{
			ID:          sub.ID.String(),
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID.String(),
			StartDate:   sub.StartDate.Format("01-2006"),
		}
		if sub.EndDate != nil {
			endDateStr := sub.EndDate.Format("01-2006")
			resp.EndDate = endDateStr
		}
		response = append(response, resp)
	}

	return response, nil
}
