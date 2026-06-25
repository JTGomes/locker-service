package postgresql

import (
	"context"
	"errors"
	"fmt"
	"locker-service/internal/api"
	"locker-service/internal/rent"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RentRepository struct {
	pool *pgxpool.Pool
}

func NewRentRepository(pool *pgxpool.Pool) *RentRepository {
	return &RentRepository{pool: pool}
}

func (rentRepo *RentRepository) Create(ctx context.Context, r *rent.Rent) error {
	const query = `INSERT INTO rents (weight, size, status) VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	err := rentRepo.pool.QueryRow(ctx, query, r.Weight, r.Size, r.Status).Scan(&r.ID, &r.CreatedAt, &r.UpdatedAt)
	if err != nil {
		return fmt.Errorf("inserting rent: %w", mapError(err))
	}

	return nil
}

func (rentRepo *RentRepository) GetByID(ctx context.Context, id string) (*rent.Rent, error) {
	const query = `SELECT id, locker_id, weight, size, status, created_at, updated_at, dropped_off_at, picked_up_at 
		FROM rents WHERE id = $1`

	var r rent.Rent

	err := querier(ctx, rentRepo.pool).QueryRow(ctx, query, id).Scan(
		&r.ID,
		&r.LockerID,
		&r.Weight,
		&r.Size,
		&r.Status,
		&r.CreatedAt,
		&r.UpdatedAt,
		&r.DroppedOffAt,
		&r.PickedUpAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getting rent %s: %w", id, mapError(err))
	}

	return &r, nil
}

func (rentRepo *RentRepository) GetByIDForUpdate(ctx context.Context, id string) (*rent.Rent, error) {
	const query = `SELECT id, locker_id, weight, size, status, created_at, updated_at, dropped_off_at, picked_up_at 
		FROM rents WHERE id = $1 FOR UPDATE`

	var r rent.Rent
	err := querier(ctx, rentRepo.pool).QueryRow(ctx, query, id).Scan(
		&r.ID, &r.LockerID, &r.Weight, &r.Size, &r.Status,
		&r.CreatedAt, &r.UpdatedAt, &r.DroppedOffAt, &r.PickedUpAt,
	)
	if err != nil {
		return nil, fmt.Errorf("getting rent %s for update: %w", id, mapError(err))
	}
	return &r, nil
}

func (rentRepo *RentRepository) UpdateAllocation(ctx context.Context, r *rent.Rent) error {
	const query = `UPDATE rents SET locker_id = $2, status = $3 WHERE id = $1
		RETURNING updated_at`
	err := querier(ctx, rentRepo.pool).QueryRow(ctx, query, r.ID, r.LockerID, r.Status).Scan(&r.UpdatedAt)

	if err != nil {
		return fmt.Errorf("updating rent allocation %w", mapError(err))
	}

	return nil
}

func (rentRepo *RentRepository) UpdateDropoffStatus(ctx context.Context, id string, from, to rent.Status) (*rent.Rent, error) {
	const query = `
		UPDATE rents
		SET status = $3, dropped_off_at = now()
		WHERE id = $1 AND status = $2
		RETURNING id, locker_id, weight, size, status, created_at, updated_at, dropped_off_at, picked_up_at`

	var r rent.Rent
	err := querier(ctx, rentRepo.pool).QueryRow(ctx, query, id, from, to).Scan(
		&r.ID, &r.LockerID, &r.Weight, &r.Size, &r.Status,
		&r.CreatedAt, &r.UpdatedAt, &r.DroppedOffAt, &r.PickedUpAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, rentRepo.statusConflict(ctx, id, from)
		}
		return nil, fmt.Errorf("dropoff rent %s: %w", id, mapError(err))
	}
	return &r, nil
}

func (rentRepo *RentRepository) statusConflict(ctx context.Context, id string, expected rent.Status) error {
	const q = `SELECT EXISTS(SELECT 1 FROM rents WHERE id = $1)`
	var exists bool
	if err := querier(ctx, rentRepo.pool).QueryRow(ctx, q, id).Scan(&exists); err != nil {
		return fmt.Errorf("checking rent %s: %w", id, mapError(err))
	}
	if !exists {
		return api.ErrNotFound
	}
	return fmt.Errorf("%w: rent %s is not in %q state", api.ErrConflict, id, expected)
}

func (rentRepo *RentRepository) UpdatePickup(ctx context.Context, r *rent.Rent) error {
	const query = `UPDATE rents SET status = $2, picked_up_at = now() WHERE id = $1
		RETURNING picked_up_at, updated_at`
	err := querier(ctx, rentRepo.pool).QueryRow(ctx, query, r.ID, r.Status).Scan(&r.PickedUpAt, &r.UpdatedAt)

	if err != nil {
		return fmt.Errorf("updating rent pickup %w", mapError(err))
	}
	return nil
}
