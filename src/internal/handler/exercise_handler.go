package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/util"
	"workoutpal/src/util/constants"

	"github.com/go-chi/render"
)

type exerciseHandler struct {
	exerciseService service.ExerciseService
}

func NewExerciseHandler(es service.ExerciseService) handler.ExerciseHandler {
	return &exerciseHandler{
		exerciseService: es,
	}
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
	exercises, err := h.exerciseService.ReadAllExercises()
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	render.JSON(w, r, exercises)
}

// ReadExerciseByID godoc
// @Summary Returns the exercise with the corresponding ID
// @Tags Exercises
// @Accept json
// @Produce json
// @Success 200 {object} model.Exercise "Exercises retrieved successfully"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Security BearerAuth
// @Router /exercises/{id} [get]
func (h *exerciseHandler) ReadExerciseByID(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(constants.ID_KEY).(int64)
	exercise, err := h.exerciseService.ReadExerciseByID(id)
	if err != nil {
		responseErr := util.Error(err, r.URL.Path)
		util.ErrorResponse(w, r, responseErr)
		return
	}
	render.JSON(w, r, exercise)
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
