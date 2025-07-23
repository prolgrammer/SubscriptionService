package subscription

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"subscription_service/internal/usecases"
)

func (r *subRepo) Delete(ctx context.Context, subID string) error {
	sql, args, err := r.client.Builder.
		Delete(commands.SubscriptionTable).
		Where("id = ?", subID).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build delete query")
		return errors.Wrap(err, "failed to build query")
	}

	result, err := r.client.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute delete query")
		return errors.Wrap(err, "failed to delete subscription")
	}

	if result.RowsAffected() == 0 {
		r.logger.Error().Msg("Subscription not found")
		return usecases.ErrEntityNotFound
	}

	return nil
}
