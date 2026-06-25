package api

import "errors"

var (
	ErrNotFound           = errors.New("record not found")
	ErrNoLockersAvailable = errors.New("no lockers available")
	ErrValidation         = errors.New("invalid input")
	ErrConflict           = errors.New("conflict")
)
