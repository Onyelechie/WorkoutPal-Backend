package handler

import (
	"net/http"
	"workoutpal/src/internal/domain/handler"
)

type authHandler struct{}

func NewAuthHandler() handler.AuthHandler {
	return &authHandler{}
}

// Login godoc
// @Summary Authenticates a user and issues an access token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login credentials including email and password"
// @Success 200 {object} model.BasicResponse "Successful login"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Router /login [post]
func (a *authHandler) Login(w http.ResponseWriter, r *http.Request) {}

// Logout godoc
// @Summary Logs out an authenticated user and invalidates their session
// @Tags Authentication
// @Produce json
// @Success 200 {object} model.BasicResponse "Logout successful"
// @Failure 400 {object} model.BasicResponse "Validation error"
// @Failure 401 {object} model.BasicResponse "Unauthorized"
// @Failure 500 {object} model.BasicResponse "Internal server error"
// @Router /logout [post]
func (a *authHandler) Logout(w http.ResponseWriter, r *http.Request) {}
