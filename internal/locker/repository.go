package locker

import (
	"context"
	"locker-service/internal/api"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, l *Locker) error
	GetByID(ctx context.Context, id string) (*Locker, error)
	List(ctx context.Context, filter LockerFilterQuery) (api.Page[Locker], error)
	Delete(ctx context.Context, id string) error
	GetFreeByBloqID(ctx context.Context, bloqID uuid.UUID) (uuid.UUID, error)
	SetOccupied(ctx context.Context, id uuid.UUID, occupied bool) error
}
