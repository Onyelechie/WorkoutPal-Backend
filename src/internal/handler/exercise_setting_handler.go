package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
	"workoutpal/src/util"
	"workoutpal/src/util/constants"

	"github.com/go-chi/render"
)

type ExerciseSettingHandler struct {
	exerciseSettingService service.ExerciseSettingService
}

func NewExerciseSettingHandler(es service.ExerciseSettingService) *ExerciseSettingHandler {
	return &ExerciseSettingHandler{
		exerciseSettingService: es,
	}
}

// ReadExerciseSetting godoc
// @Tags Exercise Setting
// @Accept json
// @Produce json
// @Param exercise_id query string true "Exercise ID"
// @Param workout_routine_id query string true "Workout Routine ID"
// @Success 200 {object} model.ExerciseSetting "Exercise Setting"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercise-settings [get]
func (h *ExerciseSettingHandler) ReadExerciseSetting(w http.ResponseWriter, r *http.Request) {
	var req model.ReadExerciseSettingRequest
	req.UserID = r.Context().Value(constants.USER_ID_KEY).(int64)

	exerciseIDString, ok := r.URL.Query()["exercise_id"]
	if !ok {
		responseErr := util.Error(errors.New("missing exercise_id"), r.URL.Path)
		util.ErrorResponseWithStatus(w, r, responseErr, 400)
		return
	}
	exerciseID, err := strconv.ParseInt(exerciseIDString[0], 10, 64)
	if err != nil {
		responseErr := util.Error(errors.New("invalid exercise_id"), r.URL.Path)
		util.ErrorResponseWithStatus(w, r, responseErr, 400)
		return
	}
	req.ExerciseID = exerciseID

	workoutRoutineIDString, ok := r.URL.Query()["workout_routine_id"]
	if !ok {
		responseErr := util.Error(errors.New("missing workout_routine_id"), r.URL.Path)
		util.ErrorResponseWithStatus(w, r, responseErr, 400)
		return
	}
	workoutRoutineID, err := strconv.ParseInt(workoutRoutineIDString[0], 10, 64)
	if err != nil {
		responseErr := util.Error(errors.New("invalid workout_routine_id"), r.URL.Path)
		util.ErrorResponseWithStatus(w, r, responseErr, 400)
		return
	}
	req.WorkoutRoutineID = workoutRoutineID

	setting, err := h.exerciseSettingService.ReadExerciseSetting(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	render.JSON(w, r, setting)
}

// CreateExerciseSetting godoc
// @Tags Exercise Setting
// @Accept json
// @Produce json
// @Param request body model.CreateExerciseSettingRequest true "New exercise setting payload"
// @Success 201 {object} model.ExerciseSetting "Exercise Setting created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercise-settings [post]
func (h *ExerciseSettingHandler) CreateExerciseSetting(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)

	var req model.CreateExerciseSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	req.UserID = userID

	setting, err := h.exerciseSettingService.CreateExerciseSetting(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, setting)
}

// UpdateExerciseSetting godoc
// @Tags Exercise Setting
// @Accept json
// @Produce json
// @Param request body model.UpdateExerciseSettingRequest true "Update exercise setting payload"
// @Success 200 {object} model.ExerciseSetting "Exercise Setting updated successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercise-settings [put]
func (h *ExerciseSettingHandler) UpdateExerciseSetting(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(constants.USER_ID_KEY).(int64)

	var req model.UpdateExerciseSettingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	// never trust client for user identity
	req.UserID = userID

	setting, err := h.exerciseSettingService.UpdateExerciseSetting(req)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, setting)
}
