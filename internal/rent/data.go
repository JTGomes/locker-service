package rent

import (
	"locker-service/internal/api"
	"slices"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusCreated        Status = "created"
	StatusWaitingDropoff Status = "waiting_dropoff"
	StatusWaitingPickup  Status = "waiting_pickup"
	StatusDelivered      Status = "delivered"
)

var allowedTransitions = map[Status][]Status{
	StatusCreated:        {StatusWaitingDropoff},
	StatusWaitingDropoff: {StatusWaitingPickup},
	StatusWaitingPickup:  {StatusDelivered},
	StatusDelivered:      {},
}

func (s Status) CanTransitionTo(next Status) bool {
	return slices.Contains(allowedTransitions[s], next)
}

type Size string

const (
	SizeXS Size = "XS"
	SizeS  Size = "S"
	SizeM  Size = "M"
	SizeL  Size = "L"
	SizeXL Size = "XL"
)

func (s Size) Valid() bool {
	switch s {
	case SizeXS, SizeS, SizeM, SizeL, SizeXL:
		return true
	default:
		return false
	}
}

type Rent struct {
	ID           uuid.UUID  `json:"id"`
	LockerID     *uuid.UUID `json:"locker_id"`
	Weight       float64    `json:"weight"`
	Size         Size       `json:"size"`
	Status       Status     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DroppedOffAt *time.Time `json:"dropped_off_at"`
	PickedUpAt   *time.Time `json:"picked_up_at"`
}

type createInput struct {
	Size   Size    `json:"size" binding:"required"`
	Weight float64 `json:"weight" binding:"required"`
}

func (c createInput) Validate() error {
	var errs []api.FieldError

	if !c.Size.Valid() {
		errs = append(errs, api.FieldError{
			Field:   "size",
			Message: "invalid size",
		})
	}

	if c.Weight <= 0 {
		errs = append(errs, api.FieldError{
			Field:   "weight",
			Message: "weight must be greater than 0",
		})
	}

	if len(errs) > 0 {
		return api.ValidationError{
			Fields: errs,
		}
	}

	return nil
}

type allocateLockerInput struct {
	BloqID uuid.UUID `json:"bloq_id" binding:"required"`
}
