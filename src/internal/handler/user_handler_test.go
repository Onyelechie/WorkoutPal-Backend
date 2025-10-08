package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	"workoutpal/src/internal/service"

	"github.com/go-chi/chi/v5"
)

func TestCreateUser(t *testing.T) {
	// Setup
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// Test data
	user := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Execute
	userHandler.CreateNewUser(w, req)

	// Assert
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestGetAllUsers(t *testing.T) {
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	userHandler.ReadAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}