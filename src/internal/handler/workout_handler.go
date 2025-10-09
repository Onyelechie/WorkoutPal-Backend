package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type workoutHandler struct {
	userService service.UserService
}

func NewWorkoutHandler(us service.UserService) handler.WorkoutHandler {
	return &workoutHandler{
		userService: us,
	}
}

// ReadWorkouts godoc
// @Summary List workouts
// @Tags Workouts
// @Accept json
// @Produce json
// @Success 200 {array} model.Workout "Workouts retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /workouts [get]
func (h *workoutHandler) ReadWorkouts(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// CreateWorkout godoc
// @Summary Create a new workout
// @Tags Workouts
// @Accept json
// @Produce json
// @Param request body model.CreateWorkoutRequest true "New workout payload"
// @Success 201 {object} model.Workout "Workout created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /workouts [post]
func (h *workoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// UpdateWorkout godoc
// @Summary Update an existing workout
// @Tags Workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout ID"
// @Param request body model.UpdateWorkoutRequest true "Updated workout payload"
// @Success 200 {object} model.Workout "Workout updated successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /workouts/{id} [put]
func (h *workoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
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
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	var req model.CreateRoutineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	routine, err := h.userService.CreateRoutine(id, req)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, routine)
}

// GetUserRoutines godoc
// @Summary Get all routines for a user
// @Tags Routines
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} model.ExerciseRoutine "Routines retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Invalid user ID"
// @Failure 404 {object} model.BasicResponse "User not found"
// @Router /users/{id}/routines [get]
func (h *workoutHandler) GetUserRoutines(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid user ID"})
		return
	}

	routines, err := h.userService.GetUserRoutines(id)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, model.BasicResponse{Message: err.Error()})
		return
	}

	render.JSON(w, r, routines)
}
