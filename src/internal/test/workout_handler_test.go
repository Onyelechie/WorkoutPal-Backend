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

// Workout handler test helpers
func setupWorkoutHandler() (handler.WorkoutHandler, service.UserService) {
	repo := repository.NewInMemoryUserRepository()
	userService := service_impl.NewUserService(repo)
	workoutHandler := handler_impl.NewWorkoutHandler(userService)
	return workoutHandler, userService
}

func createTestRoutine() model.CreateRoutineRequest {
	return model.CreateRoutineRequest{
		Name:        "Morning Workout",
		Description: "Daily morning routine",
	}
}

func TestWorkoutHandler_CreateUserRoutine(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	body, _ := json.Marshal(routineReq)
	req := createRequestWithContext("POST", "/users/1/routines", "1", body)
	w := httptest.NewRecorder()

	workoutHandler.CreateUserRoutine(w, req)

	assertStatusCode(t, http.StatusCreated, w.Code)

	var response model.ExerciseRoutine
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, routineReq.Name, response.Name, "routine name")
}

func TestWorkoutHandler_GetUserRoutines(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	userService.CreateRoutine(1, routineReq)

	req := createRequestWithContext("GET", "/users/1/routines", "1", nil)
	w := httptest.NewRecorder()

	workoutHandler.GetUserRoutines(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []model.ExerciseRoutine
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) != 1 {
		t.Errorf("Expected 1 routine, got %d", len(response))
	}
}