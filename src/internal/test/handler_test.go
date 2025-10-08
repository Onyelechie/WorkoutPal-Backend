package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/handler"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	"workoutpal/src/internal/service"

	"github.com/go-chi/chi/v5"
)

func TestUserHandler_CreateNewUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

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

	userHandler.CreateNewUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, response.Username)
	}
}

func TestUserHandler_ReadAllUsers(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// Create test users
	users := []model.CreateUserRequest{
		{Username: "user1", Email: "user1@example.com", Name: "User 1"},
		{Username: "user2", Email: "user2@example.com", Name: "User 2"},
	}

	for _, user := range users {
		userService.CreateUser(user)
	}

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	userHandler.ReadAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response []model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 2 {
		t.Errorf("Expected 2 users, got %d", len(response))
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// Create user
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	userService.CreateUser(userReq)

	req := httptest.NewRequest("GET", "/users/1", nil)
	w := httptest.NewRecorder()

	// Add chi context for URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userHandler.GetUserByID(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Username != userReq.Username {
		t.Errorf("Expected username %s, got %s", userReq.Username, response.Username)
	}
}

func TestUserHandler_UpdateUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// Create user
	createReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	userService.CreateUser(createReq)

	// Update request
	updateReq := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		Name:     "Updated User",
	}

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi context for URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userHandler.UpdateUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Username != updateReq.Username {
		t.Errorf("Expected username %s, got %s", updateReq.Username, response.Username)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// Create user
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	userService.CreateUser(userReq)

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	w := httptest.NewRecorder()

	// Add chi context for URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userHandler.DeleteUser(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Message != "User deleted successfully" {
		t.Errorf("Expected success message, got %s", response.Message)
	}
}

func TestUserHandler_CreateUserGoal(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	userService := service.NewUserService(repo)
	userHandler := handler.NewUserHandler(userService)

	// Create user
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
	}
	userService.CreateUser(userReq)

	// Create goal
	goalReq := model.CreateGoalRequest{
		Name:        "Weight Loss Goal",
		Description: "Lose weight to 65kg",
		Deadline:    "2024-12-31",
	}

	body, _ := json.Marshal(goalReq)
	req := httptest.NewRequest("POST", "/users/1/goals", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add chi context for URL parameters
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	userHandler.CreateUserGoal(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response model.Goal
	json.Unmarshal(w.Body.Bytes(), &response)
	if response.Name != goalReq.Name {
		t.Errorf("Expected goal name %s, got %s", goalReq.Name, response.Name)
	}
}