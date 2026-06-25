package rent

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

	r, err := h.svc.Create(ctx, req)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, r)
}

func (h *Handler) Get(ctx *gin.Context) {
	r, err := h.svc.Get(ctx, ctx.Param("id"))
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, r)
}

func (h *Handler) AllocateLocker(ctx *gin.Context) {
	var req allocateLockerInput
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.ErrorResponse(ctx, api.MapBindingError(err))
		return
	}

	r, err := h.svc.AllocateLocker(ctx, ctx.Param("id"), req)
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, r)
}

func (h *Handler) Dropoff(ctx *gin.Context) {
	r, err := h.svc.Dropoff(ctx, ctx.Param("id"))
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, r)
}

func (h *Handler) Pickup(ctx *gin.Context) {
	r, err := h.svc.Pickup(ctx, ctx.Param("id"))
	if err != nil {
		api.ErrorResponse(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, r)
}
