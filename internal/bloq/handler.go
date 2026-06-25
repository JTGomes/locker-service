package bloq

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

	bloqs, err := h.svc.List(ctx, pagination)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, bloqs)
}

func (h *Handler) Get(ctx *gin.Context) {
	b, err := h.svc.Get(ctx, ctx.Param("id"))
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, b)
}

func (h *Handler) Delete(ctx *gin.Context) {
	if err := h.svc.Delete(ctx, ctx.Param("id")); err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
