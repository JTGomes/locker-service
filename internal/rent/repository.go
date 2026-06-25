package rent

import "context"

type Repository interface {
	Create(ctx context.Context, r *Rent) error
	GetByID(ctx context.Context, id string) (*Rent, error)
	GetByIDForUpdate(ctx context.Context, id string) (*Rent, error)
	UpdateAllocation(ctx context.Context, r *Rent) error
	UpdateDropoffStatus(ctx context.Context, id string, old, new Status) (*Rent, error)
	UpdatePickup(ctx context.Context, r *Rent) error
}
