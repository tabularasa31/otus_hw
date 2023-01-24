package v1

import (
	"go.uber.org/zap"
	"net/http"

	"github.com/swaggo/swag/example/basic/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
)

// NewRouter -.
// Swagger spec:
//
//	@title			Calendar API
//	@description	Homework project
//	@version		1.0
//	@host			localhost:8080
//	@BasePath		/api/v1
func NewRouter(handler *gin.Engine, logg *zap.SugaredLogger, u usecase.EventUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Health check
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	v1 := handler.Group("/api/v1")
	{
		newCalendarRoutes(v1, u, *logg)
	}

	// Ping test
	handler.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
}
