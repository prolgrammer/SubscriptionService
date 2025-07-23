package subscription

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"subscription_service/internal/entities"
	"subscription_service/internal/usecases"
)

func (r *subRepo) Update(ctx context.Context, sub *entities.Subscription) error {
	sql, args, err := r.client.Builder.
		Update(commands.SubscriptionTable).
		Set(commands.SubscriptionServiceNameField, sub.ServiceName).
		Set(commands.SubscriptionPriceField, sub.Price).
		Set(commands.SubscriptionUserIDField, sub.UserID).
		Set(commands.SubscriptionStartDateField, sub.StartDate).
		Set(commands.SubscriptionEndDateField, sub.EndDate).
		Where("id = ?", sub.ID).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build update query")
		return errors.Wrap(err, "failed to build query")
	}

	result, err := r.client.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute update query")
		return errors.Wrap(err, "failed to update subscription")
	}

	if result.RowsAffected() == 0 {
		r.logger.Error().Msg("Subscription not found")
		return usecases.ErrEntityNotFound
	}

	return nil
}
