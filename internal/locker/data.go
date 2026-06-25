package locker

import (
	"locker-service/internal/api"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusOpen   Status = "open"
	StatusClosed Status = "closed"
)

func (s Status) Valid() bool {
	switch s {
	case StatusOpen, StatusClosed:
		return true
	default:
		return false
	}
}

// var allowedTransitions = map[Status][]Status{
// 	StatusOpen:   {StatusClosed},
// 	StatusClosed: {StatusOpen},
// }

// func (s Status) CanTransitionTo(next Status) bool {
// 	return slices.Contains(allowedTransitions[s], next)
// }

type Locker struct {
	ID         uuid.UUID `json:"id"`
	BloqId     uuid.UUID `json:"bloq_id"`
	Status     Status    `json:"status"`
	IsOccupied bool      `json:"is_occupied"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type createInput struct {
	BloqId uuid.UUID `json:"bloq_id" binding:"required"`
	Status Status    `json:"status" binding:"required"`
}

func (c createInput) Validate() error {
	var errs []api.FieldError

	if !c.Status.Valid() {
		errs = append(errs, api.FieldError{
			Field:   "status",
			Message: "invalid status",
		})
	}

	if len(errs) > 0 {
		return api.ValidationError{
			Fields: errs,
		}
	}

	return nil
}

type LockerFilterQuery struct {
	BloqID string `form:"bloq_id" binding:"omitempty,uuid4"`
	api.Pagination
	IsOccupied *bool `form:"is_occupied" binding:"omitempty"`
}
