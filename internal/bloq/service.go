package bloq

import (
	"context"
	"locker-service/internal/api"
)

type Service interface {
	Create(ctx context.Context, in createInput) (*Bloq, error)
	Get(ctx context.Context, id string) (*Bloq, error)
	List(ctx context.Context, pagination api.Pagination) (api.Page[Bloq], error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dataInput createInput) (*Bloq, error) {
	bloq := &Bloq{
		Title:   dataInput.Title,
		Address: dataInput.Address,
	}

	if err := s.repo.Create(ctx, bloq); err != nil {
		return nil, err
	}

	return bloq, nil
}

func (s *service) Get(ctx context.Context, id string) (*Bloq, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *service) List(ctx context.Context, pagination api.Pagination) (api.Page[Bloq], error) {
	return s.repo.List(ctx, pagination)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
