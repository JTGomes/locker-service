package locker

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"locker-service/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type mockService struct {
	createErr error
	listErr   error
}

func (m *mockService) Create(ctx context.Context, in createInput) (*Locker, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &Locker{BloqId: in.BloqId, Status: in.Status}, nil
}

func (m *mockService) Get(ctx context.Context, id string) (*Locker, error) {
	return &Locker{ID: uuid.New(), Status: StatusClosed}, nil
}

func (m *mockService) List(ctx context.Context, filter LockerFilterQuery) (api.Page[Locker], error) {
	if m.listErr != nil {
		return api.Page[Locker]{}, m.listErr
	}
	return api.Page[Locker]{Data: []Locker{{Status: StatusClosed}}, Total: 1}, nil
}

func (m *mockService) Delete(ctx context.Context, id string) error {
	return nil
}

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creating locker", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		payload := bytes.NewBufferString(fmt.Sprintf(`{"bloq_id":"%s","status":"closed"}`, uuid.New().String()))
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", payload)
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Create(ctx)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("creating locker with invalid status", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		payload := bytes.NewBufferString(fmt.Sprintf(`{"bloq_id":"%s","status":"invalid"}`, uuid.New().String()))
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", payload)
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Create(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}

func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("list lockers invalid limit", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/?limit=9999", nil)

		h.List(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})

	t.Run("list lockers invalid bloq_id filter", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/?bloq_id=not-a-uuid", nil)

		h.List(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}
