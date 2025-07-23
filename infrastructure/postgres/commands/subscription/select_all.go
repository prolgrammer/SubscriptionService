package subscription

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"subscription_service/internal/entities"
)

func (r *subRepo) SelectAll(ctx context.Context, limit, offset int) ([]entities.Subscription, error) {
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
		Limit(uint64(limit)).
		Offset(uint64(offset)).
		ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build select all query")
		return nil, errors.Wrap(err, "failed to build query")
	}

	rows, err := r.client.Pool.Query(ctx, sql, args...)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute select all query")
		return nil, errors.Wrap(err, "failed to get subscriptions")
	}
	defer rows.Close()

	var subscriptions []entities.Subscription
	for rows.Next() {
		var sub entities.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate, // *time.Time — работает корректно с NULL
		)
		if err != nil {
			r.logger.Error().Err(err).Msg("Failed to scan subscription row")
			return nil, errors.Wrap(err, "failed to scan subscription")
		}
		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error().Err(err).Msg("Error iterating subscription rows")
		return nil, errors.Wrap(err, "failed to get subscriptions")
	}

	return subscriptions, nil
}
