package handler

import "net/http"

type UserHandler interface {
	CreateNewUser(w http.ResponseWriter, r *http.Request)
	ReadAllUsers(w http.ResponseWriter, r *http.Request)
	GetUserByID(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	// Goals
	CreateUserGoal(w http.ResponseWriter, r *http.Request)
	GetUserGoals(w http.ResponseWriter, r *http.Request)
	// Followers
	FollowUser(w http.ResponseWriter, r *http.Request)
	GetUserFollowers(w http.ResponseWriter, r *http.Request)
	GetUserFollowing(w http.ResponseWriter, r *http.Request)
	// Routines
	CreateUserRoutine(w http.ResponseWriter, r *http.Request)
	GetUserRoutines(w http.ResponseWriter, r *http.Request)
}
