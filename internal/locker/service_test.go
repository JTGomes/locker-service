package locker

import (
	"context"
	"errors"
	"locker-service/internal/api"
	"testing"

	"github.com/google/uuid"
)

type mockRepo struct {
	created   *Locker
	createErr error
	getErr    error
	listErr   error
	deleteErr error
}

func (f *mockRepo) Create(ctx context.Context, l *Locker) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = l
	return nil
}
func (f *mockRepo) GetByID(ctx context.Context, id string) (*Locker, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &Locker{Status: StatusClosed}, nil
}
func (f *mockRepo) List(ctx context.Context, filter LockerFilterQuery) (api.Page[Locker], error) {
	if f.listErr != nil {
		return api.Page[Locker]{}, f.listErr
	}
	return api.Page[Locker]{Data: []Locker{{Status: StatusClosed}}, Total: 1}, nil
}
func (f *mockRepo) Delete(ctx context.Context, id string) error { return f.deleteErr }
func (f *mockRepo) GetFreeByBloqID(ctx context.Context, bloqID uuid.UUID) (uuid.UUID, error) {
	return uuid.New(), nil
}
func (f *mockRepo) SetOccupied(ctx context.Context, id uuid.UUID, occupied bool) error { return nil }

func TestService(t *testing.T) {
	t.Run("creating locker", func(t *testing.T) {
		repo := &mockRepo{}
		svc := NewService(repo)
		in := createInput{BloqId: uuid.New(), Status: StatusClosed}
		l, err := svc.Create(context.Background(), in)
		if err != nil {
			t.Fatalf("create error: %v", err)
		}
		if l.Status != in.Status {
			t.Fatalf("status mismatch")
		}
	})

	t.Run("creating locker with invalid status", func(t *testing.T) {
		in := createInput{BloqId: uuid.New(), Status: "invalid"}
		if err := in.Validate(); err == nil {
			t.Fatalf("expected validation error")
		}
	})

	t.Run("list lockers error", func(t *testing.T) {
		repo := &mockRepo{listErr: errors.New("list fail")}
		svc := NewService(repo)
		if _, err := svc.List(context.Background(), LockerFilterQuery{}); err == nil {
			t.Fatalf("expected list error")
		}
	})
}
