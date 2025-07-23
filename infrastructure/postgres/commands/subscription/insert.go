package subscription

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"subscription_service/internal/entities"
)

func (r *subRepo) Insert(ctx context.Context, sub *entities.Subscription) error {
	sql, args, err := r.client.Builder.
		Insert(commands.SubscriptionTable).
		Columns(
			commands.SubscriptionIDField,
			commands.SubscriptionServiceNameField,
			commands.SubscriptionPriceField,
			commands.SubscriptionUserIDField,
			commands.SubscriptionStartDateField,
			commands.SubscriptionEndDateField,
		).
		Values(
			sub.ID,
			sub.ServiceName,
			sub.Price,
			sub.UserID,
			sub.StartDate,
			sub.EndDate,
		).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build insert query")
		return errors.Wrap(err, "failed to build query")
	}

	_, err = r.client.Pool.Exec(ctx, sql, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute insert query")
		return errors.Wrap(err, "failed to insert subscription")
	}

	return nil
}
