package postgresql

import (
	"context"
	"errors"
	"fmt"
	"locker-service/internal/api"
	"locker-service/internal/config"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pgForeignKeyViolation = "23503"
	pgCheckViolation      = "23514"
	pgUniqueViolation     = "23505"
)

func NewPool(ctx context.Context, dbCfg *config.PostgresConfig) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dbCfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("database dsn problem: %w", err)
	}

	// Pool settings
	cfg.MaxConns = dbCfg.MaxConns
	cfg.MinConns = dbCfg.MinConns
	cfg.MaxConnLifetime = dbCfg.MaxConnLifetime
	cfg.MaxConnIdleTime = dbCfg.MaxConnIdleTime
	cfg.HealthCheckPeriod = dbCfg.HealthCheckPeriod

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return pool, nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return api.ErrNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgForeignKeyViolation:
			return fmt.Errorf("%w: record not found", api.ErrValidation)
		case pgCheckViolation:
			return fmt.Errorf("%w: a value violates a constraint of the database", api.ErrValidation)
		case pgUniqueViolation:
			return fmt.Errorf("%w: duplicate record", api.ErrValidation)
		}
	}

	return err
}

func checkRowsAffected(cmdTag pgconn.CommandTag) error {
	if cmdTag.RowsAffected() == 0 {
		return api.ErrNotFound
	}
	return nil
}
