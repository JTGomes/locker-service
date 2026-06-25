package bloq

import "github.com/google/uuid"

type Bloq struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Address string    `json:"address"`
}

type createInput struct {
	Title   string `json:"title" binding:"required"`
	Address string `json:"address" binding:"required"`
}
