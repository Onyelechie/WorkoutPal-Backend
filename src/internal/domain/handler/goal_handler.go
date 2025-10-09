package handler

import "net/http"

type GoalHandler interface {
	CreateUserGoal(w http.ResponseWriter, r *http.Request)
	GetUserGoals(w http.ResponseWriter, r *http.Request)
}