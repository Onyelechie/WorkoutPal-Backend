package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type exerciseHandler struct{}

func NewMockExerciseHandler() handler.ExerciseHandler {
	return &exerciseHandler{}
}

func (h *exerciseHandler) ReadExercises(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.Exercise{exercise})
}

func (h *exerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, exercise)
}

func (h *exerciseHandler) ReadExerciseByID(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, exercise)
}
