package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
	"net/http"
	"strconv"
)

type calendarRoutes struct {
	u usecase.EventUseCase
	l logger.Interface
}

func newCalendarRoutes(handler *gin.RouterGroup, u usecase.EventUseCase, l logger.Interface) {
	r := &calendarRoutes{u, l}

	h := handler.Group("/event")
	{
		h.POST("/create", r.create)
		h.POST("/update", r.update)
		h.POST("/delete", r.delete)
		h.GET("/daily/:user", r.daily)
		h.GET("/weekly/:user", r.weekly)
		h.GET("/monthly/:user", r.monthly)
	}
}

// @Summary     Create event
// @Description Create event
// @ID          create
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     201 {object} entity.Event
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/create [post]
func (r *calendarRoutes) create(c *gin.Context) {
	var req entity.Event
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := r.u.Create(
		c.Request.Context(),
		entity.Event{
			Title:        req.Title,
			Desc:         req.Desc,
			UserId:       req.UserId,
			EventTime:    req.EventTime,
			Duration:     req.Duration,
			Notification: req.Notification,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - create")
		errorResponse(c, http.StatusInternalServerError, "event creating problems")

		return
	}

	c.JSON(http.StatusCreated, &result)
}

// @Summary     Update event
// @Description Update event
// @ID          update
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.Event
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/update [post]
func (r *calendarRoutes) update(c *gin.Context) {
	var req entity.Event
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - update")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}
	result, err := r.u.Update(
		c.Request.Context(),
		entity.Event{
			Title:        req.Title,
			Desc:         req.Desc,
			UserId:       req.UserId,
			EventTime:    req.EventTime,
			Duration:     req.Duration,
			Notification: req.Notification,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "event updating problems")

		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     Delete event
// @Description Delete event by event_id
// @ID          delete
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200
// @Failure     400
// @Failure     500
// @Router      /event/delete [post]
func (r *calendarRoutes) delete(c *gin.Context) {
	var req int32
	if err := c.Bind(&req); err != nil {
		r.l.Error(err, "http - v1 - delete")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}
	err := r.u.Delete(c.Request.Context(), req)

	if err != nil {
		r.l.Error(err, "http - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "event deleting problems")

		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary     Get daily events
// @Description Get daily events by user_id
// @ID          daily
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200 {object} []entity.Event
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/update [get]
func (r *calendarRoutes) daily(c *gin.Context) {
	userId, err := strconv.Atoi(c.Params.ByName("user"))
	if err != nil {
		r.l.Error(err, "http - v1 - daily")
		errorResponse(c, http.StatusBadRequest, "invalid request")

		return
	}

	var date string
	if er := c.Bind(&date); er != nil {
		r.l.Error(er, "http - v1 - daily")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := r.u.DailyEvents(
		c.Request.Context(),
		entity.Event{
			UserId:    userId,
			EventTime: date,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - daily")
		errorResponse(c, http.StatusInternalServerError, "event updating problems")

		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *calendarRoutes) weekly(c *gin.Context) {
	userId, err := strconv.Atoi(c.Params.ByName("user"))
	if err != nil {
		r.l.Error(err, "http - v1 - weekly")
		errorResponse(c, http.StatusBadRequest, "invalid request")

		return
	}
	var date string
	if er := c.Bind(&date); er != nil {
		r.l.Error(er, "http - v1 - weekly")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := r.u.WeeklyEvents(
		c.Request.Context(),
		entity.Event{
			UserId:    userId,
			EventTime: date,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - weekly")
		errorResponse(c, http.StatusInternalServerError, "event updating problems")

		return
	}

	c.JSON(http.StatusOK, result)
}

func (r *calendarRoutes) monthly(c *gin.Context) {
	userId, err := strconv.Atoi(c.Params.ByName("user"))
	if err != nil {
		r.l.Error(err, "http - v1 - monthly")
		errorResponse(c, http.StatusBadRequest, "invalid request")

		return
	}
	var date string
	if er := c.Bind(&date); er != nil {
		r.l.Error(er, "http - v1 - monthly")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := r.u.MonthlyEvents(
		c.Request.Context(),
		entity.Event{
			UserId:    userId,
			EventTime: date,
		},
	)
	if err != nil {
		r.l.Error(err, "http - v1 - monthly")
		errorResponse(c, http.StatusInternalServerError, "event updating problems")

		return
	}

	c.JSON(http.StatusOK, result)
}
