package bloq

import (
	"context"
	"locker-service/internal/api"
)

type Repository interface {
	Create(ctx context.Context, b *Bloq) error
	GetByID(ctx context.Context, id string) (*Bloq, error)
	List(ctx context.Context, pagination api.Pagination) (api.Page[Bloq], error)
	Delete(ctx context.Context, id string) error
}
