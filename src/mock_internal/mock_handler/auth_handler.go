package mock_handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type authHandler struct{}

func NewMockAuthHandler() handler.AuthHandler {
	return &authHandler{}
}

func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}

func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.BasicResponse{Message: "success"})
}

func (a *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, user)
}
