package rent

import (
	"context"
	"errors"
	"testing"
	"time"

	"locker-service/internal/api"
	"locker-service/internal/locker"

	"github.com/google/uuid"
)

type mockRentRepo struct {
	stored    *Rent
	createErr error
	getErr    error
}

func (f *mockRentRepo) Create(ctx context.Context, r *Rent) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.stored = r
	return nil
}
func (f *mockRentRepo) GetByID(ctx context.Context, id string) (*Rent, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.stored != nil {
		return f.stored, nil
	}
	return &Rent{Weight: 1.0, Size: SizeM, Status: StatusCreated}, nil
}
func (f *mockRentRepo) GetByIDForUpdate(ctx context.Context, id string) (*Rent, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	if f.stored != nil {
		return f.stored, nil
	}
	return &Rent{ID: uuid.New(), Weight: 1.0, Size: SizeM, Status: StatusCreated}, nil
}
func (f *mockRentRepo) UpdateAllocation(ctx context.Context, r *Rent) error { f.stored = r; return nil }
func (f *mockRentRepo) UpdateDropoffStatus(ctx context.Context, id string, old, new Status) (*Rent, error) {
	now := time.Now()
	if f.stored != nil {
		f.stored.Status = new
		f.stored.DroppedOffAt = &now
		return f.stored, nil
	}
	return &Rent{ID: uuid.New(), Weight: 1.0, Size: SizeM, Status: new, DroppedOffAt: &now}, nil
}
func (f *mockRentRepo) UpdatePickup(ctx context.Context, r *Rent) error { f.stored = r; return nil }

type mockLockerRepo struct {
	failGetFree bool
}

func (f *mockLockerRepo) Create(ctx context.Context, l *locker.Locker) error { return nil }
func (f *mockLockerRepo) GetByID(ctx context.Context, id string) (*locker.Locker, error) {
	return &locker.Locker{}, nil
}
func (f *mockLockerRepo) List(ctx context.Context, filter locker.LockerFilterQuery) (api.Page[locker.Locker], error) {
	return api.Page[locker.Locker]{}, nil
}
func (f *mockLockerRepo) Delete(ctx context.Context, id string) error { return nil }
func (f *mockLockerRepo) GetFreeByBloqID(ctx context.Context, bloqID uuid.UUID) (uuid.UUID, error) {
	if f.failGetFree {
		return uuid.Nil, errors.New("no free locker")
	}
	id := uuid.New()
	return id, nil
}
func (f *mockLockerRepo) SetOccupied(ctx context.Context, id uuid.UUID, occupied bool) error {
	return nil
}

type mockTxManager struct{}

func (f *mockTxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestService(t *testing.T) {
	t.Run("creating rent", func(t *testing.T) {
		rentRepo := &mockRentRepo{}
		lockerRepo := &mockLockerRepo{}
		tx := &mockTxManager{}
		svc := NewService(rentRepo, lockerRepo, tx)

		in := createInput{Size: SizeM, Weight: 2.0}
		r, err := svc.Create(context.Background(), in)
		if err != nil {
			t.Fatalf("create rent error: %v", err)
		}
		if r.Weight != in.Weight {
			t.Fatalf("weight mismatch")
		}
	})

	t.Run("create rent with repo error", func(t *testing.T) {
		rentRepo := &mockRentRepo{createErr: errors.New("create fail")}
		lockerRepo := &mockLockerRepo{}
		tx := &mockTxManager{}
		svc := NewService(rentRepo, lockerRepo, tx)

		in := createInput{Size: SizeM, Weight: 2.0}
		if _, err := svc.Create(context.Background(), in); err == nil {
			t.Fatalf("expected create error")
		}
	})

	t.Run("allocate locker when none free", func(t *testing.T) {
		rentRepo := &mockRentRepo{}
		lockerRepo := &mockLockerRepo{failGetFree: true}
		tx := &mockTxManager{}
		svc := NewService(rentRepo, lockerRepo, tx)

		if _, err := svc.AllocateLocker(context.Background(), "id", allocateLockerInput{BloqID: uuid.New()}); err == nil {
			t.Fatalf("expected allocate error")
		}
	})

	t.Run("allocate invalid transition", func(t *testing.T) {
		rentRepo := &mockRentRepo{stored: &Rent{ID: uuid.New(), Status: StatusWaitingPickup}}
		lockerRepo := &mockLockerRepo{}
		tx := &mockTxManager{}
		svc := NewService(rentRepo, lockerRepo, tx)

		if _, err := svc.AllocateLocker(context.Background(), "id", allocateLockerInput{BloqID: uuid.New()}); err == nil {
			t.Fatalf("expected transition error")
		}
	})

	t.Run("dropoff and pickup flow", func(t *testing.T) {
		rentRepo := &mockRentRepo{}
		lockerRepo := &mockLockerRepo{}
		tx := &mockTxManager{}
		svc := NewService(rentRepo, lockerRepo, tx)

		// allocate
		if _, err := svc.AllocateLocker(context.Background(), "id", allocateLockerInput{BloqID: uuid.New()}); err != nil {
			t.Fatalf("allocate error: %v", err)
		}

		// dropoff
		if _, err := svc.Dropoff(context.Background(), "id"); err != nil {
			t.Fatalf("dropoff error: %v", err)
		}

		// pickup without locker should fail
		rentRepo.stored = &Rent{ID: uuid.New(), LockerID: nil, Status: StatusWaitingPickup}
		if _, err := svc.Pickup(context.Background(), "id"); err == nil {
			t.Fatalf("expected pickup error due to missing locker")
		}
	})
}
