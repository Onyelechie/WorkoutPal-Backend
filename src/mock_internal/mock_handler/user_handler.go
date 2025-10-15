package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type userHandler struct{}

func NewMockUserHandler() handler.UserHandler {
	return &userHandler{}
}

func (u *userHandler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) ReadAllUsers(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, []model.User{user})
}

func (u *userHandler) ReadUserByID(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (u *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}
