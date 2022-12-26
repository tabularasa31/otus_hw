package v1

import (
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type calendarRoutes struct {
	u usecase.EventUseCase
	l logger.Interface
}

func newCalendarRoutes(handler *gin.RouterGroup, u usecase.EventUseCase, l logger.Interface) {
	r := &calendarRoutes{u, l}

	h := handler.Group("/cal")
	{
		h.GET("/", r.home)
		h.GET("/hello", r.hello)
	}
}

// @Summary     Welcome message
// @Description Show Hello World
// @ID          home
// @Tags  	    home
func (r *calendarRoutes) home(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

// @Summary     Welcome message
// @Description Show Hello World
// @ID          hello
// @Tags  	    hello
// @Accept      json
// @Produce     json
// @Success     200 {object}
// @Failure     500 {object}
// @Router      /cal/hello [get]
func (r *calendarRoutes) hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello, world!",
	})
}
