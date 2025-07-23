package usecases

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"subscription_service/pkg/logger"
)

type deleteSubUseCase struct {
	subRepo DeleteSubRepository
	logger  logger.Logger
}

type DeleteSubUseCase interface {
	DeleteSubscription(ctx context.Context, subID string) error
}

func NewDeleteSubUseCase(subRepo DeleteSubRepository, logger logger.Logger) DeleteSubUseCase {
	return &deleteSubUseCase{
		subRepo: subRepo,
		logger:  logger,
	}
}

func (d *deleteSubUseCase) DeleteSubscription(ctx context.Context, subID string) error {
	if _, err := uuid.Parse(subID); err != nil {
		d.logger.Error().Err(err).Msg("Invalid sub_id format")
		return errors.Wrap(ErrInvalidUUID, "failed to parse sub_id")
	}

	if err := d.subRepo.Delete(ctx, subID); err != nil {
		d.logger.Error().Err(err).Msg("Failed to delete subscription")
		return errors.Wrap(err, "failed to delete subscription")
	}

	return nil
}
