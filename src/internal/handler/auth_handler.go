package handler

import (
	"encoding/json"
	"net/http"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
)

type authHandler struct {
	userService service.UserService
}

func NewAuthHandler(us service.UserService) *authHandler {
	return &authHandler{
		userService: us,
	}
}

// GoogleAuth godoc
// @Summary Authenticate with Google OAuth
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body model.GoogleAuthRequest true "Google ID token"
// @Success 200 {object} model.AuthResponse "Authentication successful"
// @Failure 400 {object} model.BasicResponse "Invalid request or token"
// @Failure 401 {object} model.BasicResponse "Authentication failed"
// @Router /auth/google [post]
func (a *authHandler) GoogleAuth(w http.ResponseWriter, r *http.Request) {
	var req model.GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "Invalid request body"})
		return
	}

	if req.IDToken == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, model.BasicResponse{Message: "ID token is required"})
		return
	}

	// TODO: Verify Google ID token and extract user info
	// For now, return a placeholder response
	response := model.AuthResponse{
		Token: "jwt-token-placeholder",
		User: model.User{
			ID:       1,
			Email:    "user@example.com",
			Name:     "Google User",
			Provider: "google",
			IsVerified: true,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}