package test

import (
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
)

// Goal handler test helpers
func setupGoalHandler() (handler.GoalHandler, service.UserService) {
	repo := repository.NewInMemoryUserRepository()
	userService := service_impl.NewUserService(repo)
	goalHandler := handler_impl.NewGoalHandler(userService)
	return goalHandler, userService
}

func createTestGoal() model.CreateGoalRequest {
	return model.CreateGoalRequest{
		Name:        "Weight Loss Goal",
		Description: "Lose weight to 65kg",
		Deadline:    "2024-12-31",
	}
}

func TestGoalHandler_CreateUserGoal(t *testing.T) {
	goalHandler, userService := setupGoalHandler()
	createTestUser(userService)

	goalReq := createTestGoal()
	body, _ := json.Marshal(goalReq)
	req := createRequestWithContext("POST", "/users/1/goals", "1", body)
	w := httptest.NewRecorder()

	goalHandler.CreateUserGoal(w, req)

	assertStatusCode(t, http.StatusCreated, w.Code)

	var response model.Goal
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, goalReq.Name, response.Name, "goal name")
}

func TestGoalHandler_GetUserGoals(t *testing.T) {
	goalHandler, userService := setupGoalHandler()
	createTestUser(userService)

	goalReq := createTestGoal()
	userService.CreateGoal(1, goalReq)

	req := createRequestWithContext("GET", "/users/1/goals", "1", nil)
	w := httptest.NewRecorder()

	goalHandler.GetUserGoals(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []model.Goal
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 1 {
		t.Errorf("Expected 1 goal, got %d", len(response))
	}
}