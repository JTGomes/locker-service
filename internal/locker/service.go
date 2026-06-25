package locker

import (
	"context"
	"locker-service/internal/api"
)

type Service interface {
	Create(ctx context.Context, in createInput) (*Locker, error)
	Get(ctx context.Context, id string) (*Locker, error)
	List(ctx context.Context, filter LockerFilterQuery) (api.Page[Locker], error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) *service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, dataInput createInput) (*Locker, error) {
	locker := &Locker{
		BloqId: dataInput.BloqId,
		Status: dataInput.Status,
	}

	if err := s.repo.Create(ctx, locker); err != nil {
		return nil, err
	}

	return locker, nil
}

func (s *service) Get(ctx context.Context, id string) (*Locker, error) {
	l, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (s *service) List(ctx context.Context, filter LockerFilterQuery) (api.Page[Locker], error) {
	return s.repo.List(ctx, filter)
}

func (s *service) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
