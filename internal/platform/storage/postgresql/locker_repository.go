package postgresql

import (
	"context"
	"errors"
	"fmt"
	"locker-service/internal/api"
	"locker-service/internal/locker"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LockerRepository struct {
	pool *pgxpool.Pool
}

func NewLockerRepository(pool *pgxpool.Pool) *LockerRepository {
	return &LockerRepository{pool: pool}
}

func (lockerRepo *LockerRepository) Create(ctx context.Context, l *locker.Locker) error {
	const query = `
		INSERT INTO lockers (bloq_id, status)
		VALUES ($1, $2)
		RETURNING id, is_occupied, created_at, updated_at`
	err := lockerRepo.pool.QueryRow(ctx, query, l.BloqId, l.Status).Scan(&l.ID, &l.IsOccupied, &l.CreatedAt, &l.UpdatedAt)

	if err != nil {
		return fmt.Errorf("inserting locker: %w", mapError(err))
	}
	return nil
}

func (lockerRepo *LockerRepository) GetByID(ctx context.Context, id string) (*locker.Locker, error) {
	const query = `SELECT id, bloq_id, status, is_occupied, created_at, updated_at FROM lockers WHERE id = $1`
	var l locker.Locker

	err := lockerRepo.pool.QueryRow(ctx, query, id).Scan(&l.ID, &l.BloqId, &l.Status, &l.IsOccupied, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("getting locker %s: %w", id, mapError(err))
	}

	return &l, nil
}

func (lockerRepo *LockerRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM lockers WHERE id = $1`
	cmdTag, err := lockerRepo.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting locker %s: %w", id, mapError(err))
	}

	return checkRowsAffected(cmdTag)
}

func (lockerRepo *LockerRepository) List(ctx context.Context, filter locker.LockerFilterQuery) (api.Page[locker.Locker], error) {
	query := `SELECT id, bloq_id, status, is_occupied, created_at, updated_at,  COUNT(*) OVER() as total FROM lockers WHERE 1=1`

	var args []any
	i := 1

	if filter.BloqID != "" {
		query += fmt.Sprintf(" AND bloq_id = $%d", i)
		args = append(args, filter.BloqID)
		i++
	}

	if filter.IsOccupied != nil {
		query += fmt.Sprintf(" AND is_occupied = $%d", i)
		args = append(args, *filter.IsOccupied)
		i++
	}

	//pagination
	query += fmt.Sprintf(
		" ORDER BY created_at DESC LIMIT $%d OFFSET $%d",
		i, i+1,
	)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := lockerRepo.pool.Query(ctx, query, args...)
	if err != nil {
		return api.Page[locker.Locker]{}, fmt.Errorf("listing lockers: %w", mapError(err))
	}
	defer rows.Close()

	lockers := []locker.Locker{}
	var total int64

	for rows.Next() {
		var lockerObj locker.Locker
		if err := rows.Scan(
			&lockerObj.ID,
			&lockerObj.BloqId,
			&lockerObj.Status,
			&lockerObj.IsOccupied,
			&lockerObj.CreatedAt,
			&lockerObj.UpdatedAt,
			&total); err != nil {
			return api.Page[locker.Locker]{}, fmt.Errorf("listing lockers: %w", mapError(err))
		}
		lockers = append(lockers, lockerObj)
	}

	if err := rows.Err(); err != nil {
		return api.Page[locker.Locker]{}, fmt.Errorf("scanning lockers: %w", mapError(err))
	}

	return api.Page[locker.Locker]{Data: lockers, Total: total}, nil
}

func (lockerRepo *LockerRepository) GetFreeByBloqID(ctx context.Context, bloqID uuid.UUID) (uuid.UUID, error) {
	const query = `SELECT id FROM lockers where bloq_id = $1 AND is_occupied = false limit 1 for update skip locked`
	var lockerID uuid.UUID
	err := querier(ctx, lockerRepo.pool).QueryRow(ctx, query, bloqID).Scan(&lockerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.UUID{}, api.ErrNoLockersAvailable
		}
		return uuid.UUID{}, mapError(err)
	}

	return lockerID, nil
}

func (lockerRepo *LockerRepository) SetOccupied(ctx context.Context, id uuid.UUID, occupied bool) error {
	const query = `UPDATE lockers SET is_occupied = $2 WHERE id = $1`
	cmdTag, err := querier(ctx, lockerRepo.pool).Exec(ctx, query, id, occupied)
	if err != nil {
		return fmt.Errorf("setting locker %s occupancy: %w", id, mapError(err))
	}
	return checkRowsAffected(cmdTag)
}
