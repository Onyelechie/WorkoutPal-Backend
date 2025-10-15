package handler

import "net/http"

type RoutineHandler interface {
	CreateUserRoutine(w http.ResponseWriter, r *http.Request)
	ReadUserRoutines(w http.ResponseWriter, r *http.Request)
	DeleteRoutine(w http.ResponseWriter, r *http.Request)
	DeleteUserRoutine(w http.ResponseWriter, r *http.Request)
	ReadRoutineWithExercises(w http.ResponseWriter, r *http.Request)
	AddExerciseToRoutine(w http.ResponseWriter, r *http.Request)
	RemoveExerciseFromRoutine(w http.ResponseWriter, r *http.Request)
}
