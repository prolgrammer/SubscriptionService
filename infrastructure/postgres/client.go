package postgres

import (
	"context"
	"errors"
	"fmt"
	"subscription_service/config/pg"
	"subscription_service/pkg/logger"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultConnAttempts = 10
	_defaultConnTimeout  = 10 * time.Second
)

var ErrNoChanges = errors.New("no changes applied")

type Client struct {
	Builder squirrel.StatementBuilderType
	Pool    *pgxpool.Pool
	cfg     *pg.Config
	logger  logger.Logger
}

func New(config pg.Config, logger logger.Logger) (*Client, error) {
	client := &Client{
		logger: logger,
		cfg:    &config,
	}

	connAttempts := config.RetryConnectionAttempts
	if connAttempts == 0 {
		connAttempts = _defaultConnAttempts
	}

	connTimeout := config.RetryConnectionTimeout
	if connTimeout == 0 {
		connTimeout = _defaultConnTimeout
	}

	client.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	connectionString := client.connectionString()

	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		logger.Err(err).Msg("couldn't parse postgres connection string")
		return nil, err
	}

	for connAttempts > 0 {
		client.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			err = client.Pool.Ping(context.TODO())
			if err == nil {
				break
			}
		}

		logger.Err(err).Msgf("failed to connect to postgres")
		logger.Info().Msgf("Postgres client is trying to connect, attempts left: %d", connAttempts)
		<-time.After(connTimeout)

		connAttempts--
	}

	if err != nil {
		logger.Err(err).Msgf("couldn't connect to postgres")
		return nil, err
	}

	return client, nil
}

func (c *Client) connectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.cfg.User,
		c.cfg.Password,
		c.cfg.Host,
		c.cfg.Port,
		c.cfg.Database)
}

func (c *Client) MigrateUp() error {
	m, err := migrate.New(c.cfg.MigrationsPath, c.connectionString())
	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChanges
		}

		return fmt.Errorf("failed to migrate up: %w", err)
	}

	return nil
}

func (c *Client) MigrateDown() error {
	m, err := migrate.New(c.cfg.MigrationsPath, c.connectionString())
	if err != nil {
		return fmt.Errorf("failed to create migration handler: %w", err)
	}

	if err = m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return ErrNoChanges
		}

		return fmt.Errorf("failed to migrate down: %w", err)
	}

	return nil
}

func (c *Client) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}
