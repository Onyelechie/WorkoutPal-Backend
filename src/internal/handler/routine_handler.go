package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util"
	"workoutpal/src/util/constants"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type workoutHandler struct {
	routineService service.RoutineService
}

func NewRoutineHandler(rs service.RoutineService) handler.RoutineHandler {
	return &workoutHandler{
		routineService: rs,
	}
}

// CreateUserRoutine godoc
// @Summary Create a workout routine for user
// @Tags Routines
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body model.CreateRoutineRequest true "Routine payload"
// @Success 201 {object} model.ExerciseRoutine "Routine created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/routines [post]
func (h *workoutHandler) CreateUserRoutine(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	var req model.CreateRoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	routine, err := h.routineService.CreateRoutine(id, req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, routine)
}

// ReadUserRoutines godoc
// @Summary Get all routines for a user
// @Tags Routines
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.ExerciseRoutine "Routines retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/routines [get]
func (h *workoutHandler) ReadUserRoutines(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	routines, err := h.routineService.ReadUserRoutines(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, routines)
}

// DeleteRoutine godoc
// @Summary Delete a routine
// @Tags Routines
// @Produce json
// @Param id path int true "Routine ID"
// @Success 200 {object} model.BasicResponse "Routine deleted successfully"
// @Failure 400 {object} model.BasicResponse "Invalid routine ID"
// @Failure 404 {object} model.BasicResponse "Routine not found"
// @Router /routines/{id} [delete]
func (h *workoutHandler) DeleteRoutine(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	err := h.routineService.DeleteRoutine(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Routine deleted successfully"})
}

// ReadRoutineWithExercises godoc
// @Summary Get routine with exercises
// @Tags Routines
// @Produce json
// @Param id path int true "Routine ID"
// @Success 200 {object} model.ExerciseRoutine "Routine with exercises retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid routine ID"
// @Failure 404 {object} model.BasicResponse "Routine not found"
// @Router /routines/{id} [get]
func (h *workoutHandler) ReadRoutineWithExercises(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	routine, err := h.routineService.ReadRoutineWithExercises(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, routine)
}

// AddExerciseToRoutine godoc
// @Summary Add exercise to routine
// @Tags Routines
// @Accept json
// @Produce json
// @Param id path int true "Routine ID"
// @Param exercise_id query int true "Exercise ID"
// @Success 200 {object} model.BasicResponse "Exercise added to routine successfully"
// @Failure 400 {object} model.BasicResponse "Invalid ID"
// @Router /routines/{id}/exercises [post]
func (h *workoutHandler) AddExerciseToRoutine(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	exerciseIDStr := r.URL.Query().Get("exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.routineService.AddExerciseToRoutine(id, exerciseID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Exercise added to routine successfully"})
}

// RemoveExerciseFromRoutine godoc
// @Summary Remove exercise from routine
// @Tags Routines
// @Produce json
// @Param id path int true "Routine ID"
// @Param exercise_id path int true "Exercise ID"
// @Success 200 {object} model.BasicResponse "Exercise removed from routine successfully"
// @Failure 400 {object} model.BasicResponse "Invalid ID"
// @Router /routines/{id}/exercises/{exercise_id} [delete]
func (h *workoutHandler) RemoveExerciseFromRoutine(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	exerciseIDStr := chi.URLParam(r, "exercise_id")
	exerciseID, err := strconv.ParseInt(exerciseIDStr, 10, 64)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	err = h.routineService.RemoveExerciseFromRoutine(id, exerciseID)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Exercise removed from routine successfully"})
}

// DeleteUserRoutine godoc
// @Summary Delete user's routine
// @Tags Routines
// @Produce json
// @Param id path int true "User ID"
// @Param routine_id path int true "Routine ID"
// @Success 200 {object} model.BasicResponse "Routine deleted successfully"
// @Failure 400 {object} model.BasicResponse "Invalid ID"
// @Failure 404 {object} model.BasicResponse "Routine not found"
// @Router /users/{id}/routines/{routine_id} [delete]
func (h *workoutHandler) DeleteUserRoutine(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)

	err := h.routineService.DeleteRoutine(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.JSON(w, r, model.BasicResponse{Message: "Routine deleted successfully"})
}
