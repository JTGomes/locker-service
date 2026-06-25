package rent

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
	createErr   error
	allocateErr error
}

func (m *mockService) Create(ctx context.Context, in createInput) (*Rent, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}
	return &Rent{Weight: in.Weight, Size: in.Size, Status: StatusCreated}, nil
}

func (m *mockService) Get(ctx context.Context, id string) (*Rent, error) {
	return &Rent{ID: uuid.New(), Weight: 1.0, Size: SizeM, Status: StatusCreated}, nil
}

func (m *mockService) AllocateLocker(ctx context.Context, id string, in allocateLockerInput) (*Rent, error) {
	if m.allocateErr != nil {
		return nil, m.allocateErr
	}
	return &Rent{ID: uuid.New(), LockerID: ptrUUID(uuid.New()), Status: StatusWaitingDropoff}, nil
}

func (m *mockService) Dropoff(ctx context.Context, id string) (*Rent, error) {
	return &Rent{ID: uuid.New(), Status: StatusWaitingPickup}, nil
}

func (m *mockService) Pickup(ctx context.Context, id string) (*Rent, error) {
	return &Rent{ID: uuid.New(), Status: StatusDelivered}, nil
}

func TestCreateHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("creating rent", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		payload := bytes.NewBufferString(`{"size":"M","weight":1.2}`)
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", payload)
		ctx.Request.Header.Set("Content-Type", "application/json")

		h.Create(ctx)

		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected %d got %d", http.StatusCreated, recorder.Code)
		}
	})

	t.Run("creating rent with invalid size", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		payload := bytes.NewBufferString(`{"size":"INVALID","weight":1.2}`)
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

func TestAllocateLockerHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allocate locker", func(t *testing.T) {
		service := &mockService{}
		h := NewHandler(service)

		payload := bytes.NewBufferString(fmt.Sprintf(`{"bloq_id":"%s"}`, uuid.New().String()))
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", payload)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: "rent-id"}}

		h.AllocateLocker(ctx)

		if recorder.Code != http.StatusOK {
			t.Fatalf("expected %d got %d", http.StatusOK, recorder.Code)
		}
	})

	t.Run("allocate locker service error", func(t *testing.T) {
		service := &mockService{allocateErr: api.ErrConflict}
		h := NewHandler(service)

		payload := bytes.NewBufferString(fmt.Sprintf(`{"bloq_id":"%s"}`, uuid.New().String()))
		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = httptest.NewRequest(http.MethodPost, "/", payload)
		ctx.Request.Header.Set("Content-Type", "application/json")
		ctx.Params = gin.Params{{Key: "id", Value: "rent-id"}}

		h.AllocateLocker(ctx)

		if recorder.Code != http.StatusConflict {
			t.Fatalf("expected %d got %d", http.StatusConflict, recorder.Code)
		}
	})
}

func ptrUUID(u uuid.UUID) *uuid.UUID { return &u }
