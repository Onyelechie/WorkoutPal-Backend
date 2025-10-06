package handler

import "net/http"

type ExerciseHandler interface {
	ReadExercises(w http.ResponseWriter, r *http.Request)
	CreateExercise(w http.ResponseWriter, r *http.Request)
}
