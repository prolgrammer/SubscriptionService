package subscription

import (
	"context"
	"github.com/pkg/errors"
	"subscription_service/infrastructure/postgres/commands"
	"time"
)

func (r *subRepo) CalculateTotalCost(ctx context.Context, startPeriod, endPeriod time.Time, userID, serviceName *string) (int, error) {
	builder := r.client.Builder.
		Select("COALESCE(SUM(price), 0)").
		From(commands.SubscriptionTable).
		Where("start_date <= ?", endPeriod).
		Where("(end_date >= ? OR end_date IS NULL)", startPeriod)

	if userID != nil {
		builder = builder.Where("user_id = ?", *userID)
	}
	if serviceName != nil {
		builder = builder.Where("service_name = ?", *serviceName)
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to build calculate total cost query")
		return 0, errors.Wrap(err, "failed to build query")
	}

	var total int
	err = r.client.Pool.QueryRow(ctx, sql, args...).Scan(&total)
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to execute calculate total cost query")
		return 0, errors.Wrap(err, "failed to calculate total cost")
	}

	return total, nil
}
