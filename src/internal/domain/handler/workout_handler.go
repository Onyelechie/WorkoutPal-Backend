package handler

import "net/http"

type WorkoutHandler interface {
	ReadWorkouts(w http.ResponseWriter, r *http.Request)
	CreateWorkout(w http.ResponseWriter, r *http.Request)
	UpdateWorkout(w http.ResponseWriter, r *http.Request)
	CreateUserRoutine(w http.ResponseWriter, r *http.Request)
	GetUserRoutines(w http.ResponseWriter, r *http.Request)
}
