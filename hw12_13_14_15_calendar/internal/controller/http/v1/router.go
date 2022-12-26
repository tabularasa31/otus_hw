package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
)

// NewRouter -.
// Swagger spec:
// @title       Calendar API
// @description Homework project
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, logg logger.Interface, u usecase.EventUseCase) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// Routers
	h := handler.Group("/v1")
	{
		newCalendarRoutes(h, u, logg)
	}
}
