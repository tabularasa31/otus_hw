package v1

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	errapp "github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/controller/repo"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/utils/utils"

	"github.com/gin-gonic/gin"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/entity"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/internal/usecase"
	"github.com/tabularasa31/hw_otus/hw12_13_14_15_calendar/pkg/logger"
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
		h.DELETE("/delete/:id", r.delete)
		h.GET("/daily", r.daily)
		h.GET("/weekly", r.weekly)
		h.GET("/monthly", r.monthly)
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
// @Failure     422 {object} response
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
			UserID:       req.UserID,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
			Notification: req.Notification,
		},
	)
	if err != nil {
		if err == errapp.ErrEventTimeBusy {
			r.l.Error("http - v1 - create - ErrEventTimeBusy")
			errorResponse(c, http.StatusUnprocessableEntity, "this event time is already busy")
		} else {
			r.l.Error(err, "http - v1 - create")
			errorResponse(c, http.StatusInternalServerError, "event creating problems")
		}
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
// @Failure     422 {object} response
// @Failure     500 {object} response
// @Router      /event/update [post]
func (r *calendarRoutes) update(c *gin.Context) {
	var req entity.Event
	fmt.Printf("--------REQ--------, %v", c.JSON)
	if err := c.ShouldBindJSON(&req); err != nil {
		r.l.Error(err, "http - v1 - update")
		errorResponse(c, http.StatusBadRequest, "invalid request body")

		return
	}

	result, err := r.u.Update(
		c.Request.Context(),
		entity.Event{
			ID:           req.ID,
			Title:        req.Title,
			Desc:         req.Desc,
			UserID:       req.UserID,
			StartTime:    req.StartTime,
			EndTime:      req.EndTime,
			Notification: req.Notification,
		},
	)
	if err != nil {
		if err == errapp.ErrEventTimeBusy {
			r.l.Error("http - v1 - update - ErrEventTimeBusy")
			errorResponse(c, http.StatusUnprocessableEntity, "this event time is already busy")
		} else {
			r.l.Error(err, "http - v1 - update")
			errorResponse(c, http.StatusInternalServerError, "event updating problems")
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary     Delete event
// @Description Delete event by event_id
// @ID          delete
// @Tags  	    event
// @Accept      int
// @Produce     plain
// @Success     200 {string} string "Deleted Success"
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/delete [delete]
func (r *calendarRoutes) delete(c *gin.Context) {
	param := c.Param("id")
	id, e := strconv.Atoi(param)
	if e != nil {
		errorResponse(c, http.StatusBadRequest, "param ID not int")
	}

	err := r.u.Delete(c.Request.Context(), id)
	if err != nil {
		r.l.Error(err, "http - v1 - delete")
		errorResponse(c, http.StatusInternalServerError, "event deleting problems")

		return
	}

	c.String(http.StatusOK, "Deleted Success")
}

type eventsResponse struct {
	Events []entity.Event `json:"events"`
}

// @Summary     Get daily events
// @Description Get daily events by userId
// @ID          daily
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200 {object} eventsResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/daily [get]
func (r *calendarRoutes) daily(c *gin.Context) {
	uid, ok := c.GetQuery("uid")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "user id missed")
	}

	userID, err := strconv.Atoi(uid)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad user id")
	}

	date, ok := c.GetQuery("date")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "event date missed")
	}

	start, err := utils.StringToDay(date)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad event date")
	}

	end := start.Add(24 * time.Hour)

	result, err := r.u.EventsByDates(c.Request.Context(), userID, start, end)
	if err != nil {
		r.l.Error("http - v1 - daily - EventsByDates: %w", err)
		errorResponse(c, http.StatusInternalServerError, "getting daily events by date problems")
	}

	c.JSON(http.StatusOK, eventsResponse{result})
}

// @Summary     Get weekly events
// @Description Get weekly events by userId
// @ID          weekly
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200 {object} eventsResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/weekly [get]
func (r *calendarRoutes) weekly(c *gin.Context) {
	uid, ok := c.GetQuery("uid")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "user id missed")
	}

	userID, err := strconv.Atoi(uid)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad user id")
	}

	date, ok := c.GetQuery("date")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "event date missed")
	}

	start, err := utils.StringToDay(date)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad event date")
	}

	end := start.Add(7 * 24 * time.Hour)

	result, err := r.u.EventsByDates(c.Request.Context(), userID, start, end)
	if err != nil {
		r.l.Error("http - v1 - weekly - EventsByDates: %w", err)
		errorResponse(c, http.StatusInternalServerError, "getting weekly events by date problems")
	}

	c.JSON(http.StatusOK, eventsResponse{result})
}

// @Summary     Get monthly events
// @Description Get monthly events by userId
// @ID          monthly
// @Tags  	    event
// @Accept      json
// @Produce     json
// @Success     200 {object} eventsResponse
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /event/monthly [get]
func (r *calendarRoutes) monthly(c *gin.Context) {
	uid, ok := c.GetQuery("uid")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "user id missed")
	}

	userID, err := strconv.Atoi(uid)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad user id")
	}

	date, ok := c.GetQuery("date")
	if !ok {
		errorResponse(c, http.StatusBadRequest, "event date missed")
	}

	start, err := utils.StringToDay(date)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "bad event date")
	}

	end := start.Add(30 * 24 * time.Hour)

	result, err := r.u.EventsByDates(c.Request.Context(), userID, start, end)
	if err != nil {
		r.l.Error("http - v1 - monthly - EventsByDates: %w", err)
		errorResponse(c, http.StatusInternalServerError, "getting monthly events by date problems")
	}

	c.JSON(http.StatusOK, eventsResponse{result})
}
