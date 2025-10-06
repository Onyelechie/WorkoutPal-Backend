package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type exerciseHandler struct{}

func NewExerciseHandler() handler.ExerciseHandler {
	return &exerciseHandler{}
}

// ReadExercises godoc
// @Summary List exercises
// @Tags Exercises
// @Accept json
// @Produce json
// @Param target query string false "Primary target muscle"
// @Param intensity query string false "Training intensity"
// @Param expertise query string false "Recommended expertise level"
// @Success 200 {array} model.Exercise "Exercises retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercises [get]
func (h *exerciseHandler) ReadExercises(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

// CreateExercise godoc
// @Summary Create a new exercise
// @Description Creates a custom exercise with metadata such as targets, intensity, expertise, and media references.
// @Tags Exercises
// @Accept json
// @Produce json
// @Param request body model.CreateExerciseRequest true "New exercise payload"
// @Success 201 {object} model.Exercise "Exercise created successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercises [post]
func (h *exerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
