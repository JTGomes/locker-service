package bloq

import (
	"bytes"
	"context"
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

func (m *mockService) Create(ctx context.Context, in createInput) (*Bloq, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &Bloq{Title: in.Title, Address: in.Address}, nil
}

func (m *mockService) Get(ctx context.Context, id string) (*Bloq, error) {
	return &Bloq{ID: uuid.New(), Title: "test", Address: "addr"}, nil
}

func (m *mockService) List(ctx context.Context, pagination api.Pagination) (api.Page[Bloq], error) {
	if m.listErr != nil {
		return api.Page[Bloq]{}, m.listErr
	}
	return api.Page[Bloq]{Data: []Bloq{{Title: "one"}}, Total: 1}, nil
}

func (m *mockService) Delete(ctx context.Context, id string) error {
	return nil
}

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creating bloq", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		body := bytes.NewBufferString(`{"title":"Test Bloq","address":"123 Street"}`)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Create(ctx)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("creating bloq with invalid JSON", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		body := bytes.NewBufferString(`{"title":"missing_quote}`)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", body)
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Create(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}

func TestListHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("listing bloqs with invalid pagination", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodGet, "/?limit=-1", nil)

		h.List(ctx)

		if recorder.Code != http.StatusBadRequest {
			t.Fatalf("expected %d got %d", http.StatusBadRequest, recorder.Code)
		}
	})
}
