package rent

import (
	"context"
	"fmt"
	"locker-service/internal/api"
	"locker-service/internal/locker"
	"locker-service/internal/platform/storage"
)

type Service interface {
	Create(ctx context.Context, in createInput) (*Rent, error)
	Get(ctx context.Context, id string) (*Rent, error)
	AllocateLocker(ctx context.Context, id string, in allocateLockerInput) (*Rent, error)
	Dropoff(ctx context.Context, id string) (*Rent, error)
	Pickup(ctx context.Context, id string) (*Rent, error)
}

type service struct {
	repo       Repository
	lockerRepo locker.Repository
	txManager  storage.TxManager
}

func NewService(repo Repository, lockerRepo locker.Repository, txManager storage.TxManager) *service {
	return &service{repo: repo, lockerRepo: lockerRepo, txManager: txManager}
}

func (s *service) Create(ctx context.Context, in createInput) (*Rent, error) {
	rent := &Rent{
		Weight: in.Weight,
		Size:   in.Size,
		Status: StatusCreated,
	}

	if err := s.repo.Create(ctx, rent); err != nil {
		return nil, err
	}

	return rent, nil
}

func (s *service) Get(ctx context.Context, id string) (*Rent, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *service) AllocateLocker(ctx context.Context, id string, in allocateLockerInput) (*Rent, error) {
	var rent *Rent
	err := s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		r, err := s.repo.GetByIDForUpdate(ctx, id)
		if err != nil {
			return err
		}

		if !r.Status.CanTransitionTo(StatusWaitingDropoff) {
			return fmt.Errorf("%w: cannot allocate from status %q", api.ErrConflict, r.Status)
		}

		lockerID, err := s.lockerRepo.GetFreeByBloqID(ctx, in.BloqID)
		if err != nil {
			return err
		}
		if err := s.lockerRepo.SetOccupied(ctx, lockerID, true); err != nil {
			return err
		}

		r.Status = StatusWaitingDropoff
		r.LockerID = &lockerID

		if err := s.repo.UpdateAllocation(ctx, r); err != nil {
			return err
		}

		rent = r
		return nil

	})

	if err != nil {
		return nil, err
	}
	return rent, nil
}

func (s *service) Dropoff(ctx context.Context, id string) (*Rent, error) {
	return s.repo.UpdateDropoffStatus(ctx, id, StatusWaitingDropoff, StatusWaitingPickup)
}

func (s *service) Pickup(ctx context.Context, id string) (*Rent, error) {
	var rent *Rent
	err := s.txManager.WithinTx(ctx, func(ctx context.Context) error {
		r, err := s.repo.GetByIDForUpdate(ctx, id)
		if err != nil {
			return err
		}
		if !r.Status.CanTransitionTo(StatusDelivered) {
			return fmt.Errorf("%w: cannot move from %q to %q", api.ErrConflict, r.Status, StatusDelivered)
		}

		if r.LockerID == nil {
			return fmt.Errorf("%w: rent has no locker", api.ErrConflict)
		}
		r.Status = StatusDelivered
		if err := s.repo.UpdatePickup(ctx, r); err != nil {
			return err
		}
		rent = r
		return s.lockerRepo.SetOccupied(ctx, *r.LockerID, false)

	})

	if err != nil {
		return nil, err
	}
	return rent, nil
}
