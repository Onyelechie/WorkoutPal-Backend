package handler

import "net/http"

type ExerciseHandler interface {
	ReadExerciseByID(w http.ResponseWriter, r *http.Request)
	ReadExercises(w http.ResponseWriter, r *http.Request)
	CreateExercise(w http.ResponseWriter, r *http.Request)
}
