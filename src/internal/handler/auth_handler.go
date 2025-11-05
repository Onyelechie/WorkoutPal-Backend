package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/middleware"
	"workoutpal/src/internal/model"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type authHandler struct {
	userService service.UserService
	authService service.AuthService
	secret      []byte
}

func NewAuthHandler(us service.UserService, as service.AuthService, secret []byte) handler.AuthHandler {
	return &authHandler{
		userService: us,
		authService: as,
		secret:      secret,
	}
}

// Login godoc
// @Summary Logs in a user
// @Description Authenticates a user and sets access_token as cookie
// @Tags auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "comment"
// @Success 200 {object} model.User
// @Router /login [post]
func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	req.Email = strings.ToLower(req.Email)

	user, err := h.authService.Authenticate(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenWithClaims.SignedString(h.secret)
	if err != nil {
		http.Error(w, "failed to create token", http.StatusInternalServerError)
		return
	}

	isSecure := true
	if v := os.Getenv("COOKIE_SECURE"); strings.ToLower(v) == "false" {
		isSecure = false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   int(time.Hour.Seconds()),
	})

	render.JSON(w, r, user)
}

// Me godoc
// @Summary Get current authenticated user
// @Tags auth
// @Produce json
// @Success 200 {object} model.User
// @Router /me [get]
func (h *authHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		http.Error(w, "invalid token data", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.ReadUserByEmail(userEmail)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, user)
}

// Logout godoc
// @Summary Logs out user by clearing access_token
// @Tags auth
// @Success 200 {object} model.BasicResponse "successful logout"
// @Router /logout [post]
func (h *authHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, model.BasicResponse{Message: "success"})
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
			ID:         1,
			Email:      "user@example.com",
			Name:       "Google User",
			Provider:   "google",
			IsVerified: true,
		},
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, response)
}
