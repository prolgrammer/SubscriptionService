package usecases

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"subscription_service/internal/controllers/requests"
	"subscription_service/internal/controllers/responses"
	"subscription_service/pkg/logger"
	"time"
)

type CalculateTotalCostUseCase interface {
	CalculateTotalCost(c *gin.Context, req requests.CalculateTotalCost) (responses.CalculateTotalCost, error)
}

type calculateTotalCostUseCase struct {
	subRepo CalculateTotalCostRepository
	logger  logger.Logger
}

func NewCalculateTotalCostUseCase(subRepo CalculateTotalCostRepository, logger logger.Logger) CalculateTotalCostUseCase {
	return &calculateTotalCostUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (c *calculateTotalCostUseCase) CalculateTotalCost(ginCtx *gin.Context, req requests.CalculateTotalCost) (responses.CalculateTotalCost, error) {
	startPeriod, err := time.Parse("01-2006", req.StartPeriod)
	if err != nil {
		c.logger.Error().Err(err).Msg("Invalid start_period format")
		return responses.CalculateTotalCost{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse start_period")
	}

	endPeriod, err := time.Parse("01-2006", req.EndPeriod)
	if err != nil {
		c.logger.Error().Err(err).Msg("Invalid end_period format")
		return responses.CalculateTotalCost{}, errors.Wrap(ErrInvalidDateFormat, "failed to parse end_period")
	}

	var userID *string
	if req.UserID != "" {
		if _, err := uuid.Parse(req.UserID); err != nil {
			c.logger.Error().Err(err).Msg("Invalid user_id format")
			return responses.CalculateTotalCost{}, errors.Wrap(ErrInvalidUUID, "failed to parse user_id")
		}
		userID = &req.UserID
	}

	var serviceName *string
	if req.ServiceName != "" {
		serviceName = &req.ServiceName
	}

	ctx := ginCtx.Request.Context()
	total, err := c.subRepo.CalculateTotalCost(ctx, startPeriod, endPeriod, userID, serviceName)
	if err != nil {
		c.logger.Error().Err(err).Msg("Failed to calculate total cost")
		return responses.CalculateTotalCost{}, errors.Wrap(err, "failed to calculate total cost")
	}

	return responses.CalculateTotalCost{Total: total}, nil
}
