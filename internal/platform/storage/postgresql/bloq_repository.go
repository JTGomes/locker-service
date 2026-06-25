package postgresql

import (
	"context"
	"fmt"
	"locker-service/internal/api"
	"locker-service/internal/bloq"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BloqRepository struct {
	pool *pgxpool.Pool
}

func NewBloqRepository(pool *pgxpool.Pool) *BloqRepository {
	return &BloqRepository{pool: pool}
}

func (bloqRepo *BloqRepository) Create(ctx context.Context, b *bloq.Bloq) error {
	const query = `
		INSERT INTO bloqs (title, address)
		VALUES ($1, $2)
		RETURNING id`
	err := bloqRepo.pool.QueryRow(ctx, query, b.Title, b.Address).Scan(&b.ID)

	if err != nil {
		return fmt.Errorf("inserting bloq: %w", mapError(err))
	}
	return nil
}

func (bloqRepo *BloqRepository) GetByID(ctx context.Context, id string) (*bloq.Bloq, error) {
	const query = `SELECT id, title, address FROM bloqs WHERE id = $1`
	var b bloq.Bloq

	err := bloqRepo.pool.QueryRow(ctx, query, id).Scan(&b.ID, &b.Title, &b.Address)
	if err != nil {
		return nil, fmt.Errorf("getting bloq %s: %w", id, mapError(err))
	}

	return &b, nil
}

func (bloqRepo *BloqRepository) List(ctx context.Context, p api.Pagination) (api.Page[bloq.Bloq], error) {
	const query = `
		SELECT id, title, address, count(*) OVER() AS total 
		FROM bloqs 
		ORDER BY id
		LIMIT $1 OFFSET $2`

	rows, err := bloqRepo.pool.Query(ctx, query, p.Limit, p.Offset)
	if err != nil {
		return api.Page[bloq.Bloq]{}, fmt.Errorf("listing bloqs: %w", mapError(err))
	}
	defer rows.Close()

	bloqs := []bloq.Bloq{}
	var total int64

	for rows.Next() {
		var bloqObj bloq.Bloq
		if err := rows.Scan(&bloqObj.ID, &bloqObj.Title, &bloqObj.Address, &total); err != nil {
			return api.Page[bloq.Bloq]{}, fmt.Errorf("listing bloqs: %w", mapError(err))
		}
		bloqs = append(bloqs, bloqObj)
	}

	if err := rows.Err(); err != nil {
		return api.Page[bloq.Bloq]{}, fmt.Errorf("scanning bloqs: %w", mapError(err))
	}

	return api.Page[bloq.Bloq]{Data: bloqs, Total: total}, nil
}

func (bloqRepo *BloqRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM bloqs WHERE id = $1`
	cmdTag, err := bloqRepo.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deleting bloq %s: %w", id, mapError(err))
	}

	return checkRowsAffected(cmdTag)
}
