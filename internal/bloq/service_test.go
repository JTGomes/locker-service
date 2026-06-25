package bloq

import (
	"context"
	"errors"
	"testing"

	"locker-service/internal/api"
)

type mockRepo struct {
	created   *Bloq
	createErr error
	getErr    error
	listErr   error
	deleteErr error
}

func (f *mockRepo) Create(ctx context.Context, b *Bloq) error {
	if f.createErr != nil {
		return f.createErr
	}
	f.created = b
	return nil
}
func (f *mockRepo) GetByID(ctx context.Context, id string) (*Bloq, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return &Bloq{Title: "t", Address: "a"}, nil
}
func (f *mockRepo) List(ctx context.Context, pagination api.Pagination) (api.Page[Bloq], error) {
	if f.listErr != nil {
		return api.Page[Bloq]{}, f.listErr
	}
	return api.Page[Bloq]{Data: []Bloq{{Title: "t1"}, {Title: "t2"}}, Total: 2}, nil
}
func (f *mockRepo) Delete(ctx context.Context, id string) error {
	return f.deleteErr
}

func TestService(t *testing.T) {
	t.Run("creating bloq", func(t *testing.T) {
		repo := &mockRepo{}
		svc := NewService(repo)

		in := createInput{Title: "My Bloq", Address: "Addr"}
		b, err := svc.Create(context.Background(), in)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if b.Title != in.Title || b.Address != in.Address {
			t.Fatalf("mismatch created bloq")
		}
		if repo.created == nil {
			t.Fatalf("repo.Create not called")
		}
	})

	t.Run("creating bloq with repo error", func(t *testing.T) {
		repo := &mockRepo{createErr: errors.New("db error")}
		svc := NewService(repo)
		in := createInput{Title: "x", Address: "y"}
		if _, err := svc.Create(context.Background(), in); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("getting bloq", func(t *testing.T) {
		repo := &mockRepo{}
		svc := NewService(repo)
		if _, err := svc.Get(context.Background(), "id"); err != nil {
			t.Fatalf("get error: %v", err)
		}
	})

	t.Run("getting bloq not found", func(t *testing.T) {
		repo := &mockRepo{getErr: errors.New("not found")}
		svc := NewService(repo)
		if _, err := svc.Get(context.Background(), "id"); err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("listing bloqs", func(t *testing.T) {
		repo := &mockRepo{}
		svc := NewService(repo)
		page, err := svc.List(context.Background(), api.Pagination{Limit: 10, Offset: 0})
		if err != nil {
			t.Fatalf("list error: %v", err)
		}
		if len(page.Data) != 2 || page.Total != 2 {
			t.Fatalf("unexpected page result")
		}
	})

	t.Run("deleting bloq with repo error", func(t *testing.T) {
		repo := &mockRepo{deleteErr: errors.New("delete failed")}
		svc := NewService(repo)
		if err := svc.Delete(context.Background(), "id"); err == nil {
			t.Fatalf("expected delete error")
		}
	})
}
