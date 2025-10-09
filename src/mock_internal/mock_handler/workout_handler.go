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

func (h *workoutHandler) CreateUserRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, routine)
}

func (h *workoutHandler) GetUserRoutines(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.ExerciseRoutine{routine})
}

func (h *workoutHandler) DeleteRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Routine deleted successfully"})
}

func (h *workoutHandler) GetRoutineWithExercises(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, routine)
}

func (h *workoutHandler) AddExerciseToRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Exercise added to routine successfully"})
}

func (h *workoutHandler) RemoveExerciseFromRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Exercise removed from routine successfully"})
}

func (h *workoutHandler) DeleteUserRoutine(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "Routine deleted successfully"})
}
