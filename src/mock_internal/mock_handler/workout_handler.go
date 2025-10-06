package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type workoutHandler struct{}

func NewMockWorkoutHandler() handler.WorkoutHandler {
	return &workoutHandler{}
}

func (h *workoutHandler) ReadWorkouts(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.Workout{workout})
}

func (h *workoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, workout)
}

func (h *workoutHandler) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, workout)
}
