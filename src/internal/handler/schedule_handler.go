package handler

import (
	"net/http"
	"strconv"
	"workoutpal/src/util/constants"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type scheduleHandler struct {
	service service.ScheduleService
}

func NewScheduleHandler(s service.ScheduleService) handler.ScheduleHandler {
	return &scheduleHandler{service: s}
}

// ReadUserSchedules godoc
// @Summary Read all schedules for the authenticated user
// @Tags schedules
// @Success 200 {array} model.Schedule
// @Router /me/schedules [get]
func (h *scheduleHandler) ReadUserSchedules(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)

	schedules, err := h.service.ReadUserSchedules(userID)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, schedules)
}

// ReadUserSchedulesByDay godoc
// @Summary Read schedules for the authenticated user on a specific day
// @Tags schedules
// @Param dayOfWeek path int true "Day of week (0-6)"
// @Success 200 {array} model.Schedule
// @Router /me/schedules/{dayOfWeek} [get]
func (h *scheduleHandler) ReadUserSchedulesByDay(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)
	dayOfWeek, _ := strconv.ParseInt(chi.URLParam(r, constants.DAY_OF_WEEK_KEY), 10, 64)

	schedules, err := h.service.ReadUserSchedulesByDay(userID, dayOfWeek)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, schedules)
}

// ReadScheduleByID godoc
// @Summary Read schedule by ID
// @Tags schedules
// @Param id path int true "Schedule ID"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id} [get]
func (h *scheduleHandler) ReadScheduleByID(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	schedule, err := h.service.ReadScheduleByID(id)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, schedule)
}

// CreateSchedule godoc
// @Summary Create a schedule
// @Tags schedules
// @Accept json
// @Produce json
// @Param request body model.CreateScheduleRequest true "Schedule create payload"
// @Success 201 {object} model.Schedule
// @Router /schedules [post]
func (h *scheduleHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)

	var req model.CreateScheduleRequest
	_ = render.DecodeJSON(r.Body, &req)
	req.UserID = userID

	schedule, err := h.service.CreateSchedule(req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, schedule)
}

// UpdateSchedule godoc
// @Summary Update a schedule
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param request body model.UpdateScheduleRequest true "Schedule update payload"
// @Success 200 {object} model.Schedule
// @Router /schedules/{id} [put]
func (h *scheduleHandler) UpdateSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)
	id := r.Context().Value(constants.ID_KEY).(int64)

	var req model.UpdateScheduleRequest
	_ = render.DecodeJSON(r.Body, &req)
	req.ID = id
	req.UserID = userID

	schedule, err := h.service.UpdateSchedule(req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": err.Error()})
		return
	}

	render.JSON(w, r, schedule)
}

// DeleteSchedule godoc
// @Summary Delete a schedule
// @Tags schedules
// @Param id path int true "Schedule ID"
// @Success 204 {string} string "No Content"
// @Router /schedules/{id} [delete]
func (h *scheduleHandler) DeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	_ = h.service.DeleteSchedule(model.DeleteScheduleRequest{ID: id})
	render.Status(r, http.StatusNoContent)
}
