package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorBody struct {
	Error string `json:"error"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationError struct {
	Fields []FieldError
}

func (e ValidationError) Error() string {
	return "validation error"
}

func MapBindingError(err error) error {

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return ValidationError{
			Fields: formatValidationErrors(ve),
		}
	}

	return fmt.Errorf("invalid request: %w", ErrValidation)
}

func formatValidationErrors(ve validator.ValidationErrors) []FieldError {
	out := make([]FieldError, 0, len(ve))

	for _, fe := range ve {
		out = append(out, FieldError{
			Field:   strings.ToLower(fe.Field()),
			Message: fe.Tag(),
		})
	}

	return out
}

func ErrorResponse(ctx *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		ctx.JSON(http.StatusBadRequest, ValidationError{
			Fields: formatValidationErrors(ve),
		})
		return
	}

	var de ValidationError
	if errors.As(err, &de) {
		ctx.JSON(http.StatusBadRequest, de)
		return
	}

	switch {
	case errors.Is(err, ErrNotFound) || errors.Is(err, ErrNoLockersAvailable):
		ctx.JSON(http.StatusNotFound, errorBody{Error: err.Error()})
	case errors.Is(err, ErrConflict):
		ctx.JSON(http.StatusConflict, errorBody{Error: err.Error()})
	case errors.Is(err, ErrValidation):
		ctx.JSON(http.StatusBadRequest, errorBody{Error: err.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, errorBody{Error: "internal server error"})
		log.Println(err.Error())
	}

}
