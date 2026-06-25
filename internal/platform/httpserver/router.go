package httpserver

import (
	"locker-service/internal/bloq"
	"locker-service/internal/locker"
	"locker-service/internal/rent"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	Bloq   *bloq.Handler
	Locker *locker.Handler
	Rent   *rent.Handler
}

func NewRouter(dep *RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.HandleMethodNotAllowed = true

	router.GET("/health", healthHandler)

	api := router.Group("/api")

	v1 := api.Group("/v1")
	bloq := v1.Group("bloq")
	bloq.POST("", dep.Bloq.Create)
	bloq.GET("", dep.Bloq.List)
	bloq.GET("/:id", dep.Bloq.Get)
	bloq.DELETE("/:id", dep.Bloq.Delete)

	locker := v1.Group("locker")
	locker.POST("", dep.Locker.Create)
	locker.GET("", dep.Locker.List)
	locker.GET("/:id", dep.Locker.Get)
	locker.DELETE("/:id", dep.Locker.Delete)

	rent := v1.Group("rent")
	rent.POST("", dep.Rent.Create)
	rent.GET("/:id", dep.Rent.Get)
	rent.POST("/:id/allocate", dep.Rent.AllocateLocker)
	rent.POST("/:id/dropoff", dep.Rent.Dropoff)
	rent.POST("/:id/pickup", dep.Rent.Pickup)

	return router

}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
