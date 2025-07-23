package subscription

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"subscription_service/internal/entities"
	"subscription_service/internal/usecases"
)

func (r *subRepo) SelectByID(ctx context.Context, subID string) (entities.Subscription, error) {
	sql, args, err := r.client.Builder.
		Select(
			commands.SubscriptionIDField,
			commands.SubscriptionServiceNameField,
			commands.SubscriptionPriceField,
			commands.SubscriptionUserIDField,
			commands.SubscriptionStartDateField,
			commands.SubscriptionEndDateField,
		).
		From(commands.SubscriptionTable).
		Where("id = ?", subID).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build select query")
		return entities.Subscription{}, errors.Wrap(err, "failed to build query")
	}

	var sub entities.Subscription
	err = r.client.Pool.QueryRow(ctx, sql, args...).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate, // сразу в *time.Time
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error().Msg("Subscription not found")
			return entities.Subscription{}, usecases.ErrEntityNotFound
		}
		r.logger.Error().Err(err).Msg("Failed to execute select query")
		return entities.Subscription{}, errors.Wrap(err, "failed to get subscription")
	}

	return sub, nil
}
