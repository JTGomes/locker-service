package locker

import (
	"locker-service/internal/api"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Create(ctx *gin.Context) {
	var req createInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.ErrorResponse(ctx, api.MapBindingError(err))
		return
	}

	if err := req.Validate(); err != nil {
		api.ErrorResponse(ctx, err)
		return
	}

	b, err := h.svc.Create(ctx, req)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, b)
}

func (h *Handler) List(ctx *gin.Context) {
	pagination, err := api.ParsePagination(ctx)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}

	var queryFilter LockerFilterQuery

	if err := ctx.ShouldBindQuery(&queryFilter); err != nil {
		api.ErrorResponse(ctx, api.MapBindingError(err))
		return
	}

	queryFilter.Limit = pagination.Limit
	queryFilter.Offset = pagination.Offset

	lockers, err := h.svc.List(ctx, queryFilter)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, lockers)
}

func (h *Handler) Get(ctx *gin.Context) {
	l, err := h.svc.Get(ctx, ctx.Param("id"))
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, l)
}

func (h *Handler) Delete(ctx *gin.Context) {
	if err := h.svc.Delete(ctx, ctx.Param("id")); err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
