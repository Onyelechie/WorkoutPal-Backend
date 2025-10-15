package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/domain/handler"
	"workoutpal/src/internal/domain/service"
	handler_impl "workoutpal/src/internal/handler"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	service_impl "workoutpal/src/internal/service"

	"github.com/go-chi/chi/v5"
)

// Test helpers
func setupUserHandler() (handler.UserHandler, service.UserService) {
	repo := repository.NewInMemoryUserRepository()
	userService := service_impl.NewUserService(repo)
	userHandler := handler_impl.NewUserHandler(userService)
	return userHandler, userService
}

func createTestUser(userService service.UserService) model.User {
	userReq := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	user, _ := userService.CreateUser(userReq)
	return user
}

func createRequestWithContext(method, url, id string, body []byte) *http.Request {
	var req *http.Request
	if body != nil {
		req = httptest.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}

	if id != "" {
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("id", id)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	}
	return req
}

func assertStatusCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Errorf("Expected status %d, got %d", expected, actual)
	}
}

func assertResponseField(t *testing.T, expected, actual, fieldName string) {
	if actual != expected {
		t.Errorf("Expected %s %s, got %s", fieldName, expected, actual)
	}
}

func TestUserHandler_CreateNewUser(t *testing.T) {
	userHandler, _ := setupUserHandler()

	user := model.CreateUserRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	body, _ := json.Marshal(user)
	req := createRequestWithContext("POST", "/users", "", body)
	w := httptest.NewRecorder()

	userHandler.CreateNewUser(w, req)

	assertStatusCode(t, http.StatusCreated, w.Code)

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, user.Username, response.Username, "username")
}

func TestUserHandler_ReadAllUsers(t *testing.T) {
	userHandler, userService := setupUserHandler()

	// Create test users
	users := []model.CreateUserRequest{
		{Username: "user1", Email: "user1@example.com", Name: "User 1"},
		{Username: "user2", Email: "user2@example.com", Name: "User 2"},
	}

	for _, user := range users {
		userService.CreateUser(user)
	}

	req := createRequestWithContext("GET", "/users", "", nil)
	w := httptest.NewRecorder()

	userHandler.ReadAllUsers(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 2 {
		t.Errorf("Expected 2 users, got %d", len(response))
	}
}

func TestUserHandler_GetUserByID(t *testing.T) {
	userHandler, userService := setupUserHandler()
	user := createTestUser(userService)

	req := createRequestWithContext("GET", "/users/1", "1", nil)
	w := httptest.NewRecorder()

	userHandler.ReadUserByID(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, user.Username, response.Username, "username")
}

func TestUserHandler_UpdateUser(t *testing.T) {
	userHandler, userService := setupUserHandler()
	createTestUser(userService)

	updateReq := model.UpdateUserRequest{
		Username: "updateduser",
		Email:    "updated@example.com",
		Name:     "Updated User",
	}

	body, _ := json.Marshal(updateReq)
	req := createRequestWithContext("PATCH", "/users/1", "1", body)
	w := httptest.NewRecorder()

	userHandler.UpdateUser(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.User
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, updateReq.Username, response.Username, "username")
}

func TestUserHandler_DeleteUser(t *testing.T) {
	userHandler, userService := setupUserHandler()
	createTestUser(userService)

	req := createRequestWithContext("DELETE", "/users/1", "1", nil)
	w := httptest.NewRecorder()

	userHandler.DeleteUser(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, "User deleted successfully", response.Message, "message")
}
