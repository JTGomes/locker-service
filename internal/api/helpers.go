package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParsePagination(c *gin.Context) (Pagination, error) {
	limit := DefaultLimit
	offset := 0

	if inputLimit := c.Query("limit"); inputLimit != "" {
		n, err := strconv.Atoi(inputLimit)
		if err != nil || n <= 0 || n > MaxLimit {
			return Pagination{}, fmt.Errorf("limit must be between 1 and %d: %w", MaxLimit, ErrValidation)
		}
		limit = n
	}

	if inputOffset := c.Query("offset"); inputOffset != "" {
		n, err := strconv.Atoi(inputOffset)
		if err != nil || n < 0 {
			return Pagination{}, fmt.Errorf("Offset must be a non-negative integer: %w", ErrValidation)
		}
		offset = n
	}

	return Pagination{Limit: limit, Offset: offset}, nil
}
