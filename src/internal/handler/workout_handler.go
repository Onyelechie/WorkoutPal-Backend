package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type workoutHandler struct{}

func NewWorkoutHandler() handler.WorkoutHandler {
	return &workoutHandler{}
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
