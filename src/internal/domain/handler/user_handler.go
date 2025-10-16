package handler

import "net/http"

type UserHandler interface {
	CreateNewUser(w http.ResponseWriter, r *http.Request)
	ReadAllUsers(w http.ResponseWriter, r *http.Request)
	ReadUserByID(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}
