package test

import (
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

func TestWorkoutHandler_DeleteRoutine(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	_, _ = userService.CreateRoutine(1, routineReq)

	req := createRequestWithContext("DELETE", "/routines/1", "1", nil)
	w := httptest.NewRecorder()

	workoutHandler.DeleteRoutine(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, "Routine deleted successfully", response.Message, "message")

	// Verify routine was deleted
	routines, _ := userService.GetUserRoutines(1)
	if len(routines) != 0 {
		t.Errorf("Expected 0 routines after deletion, got %d", len(routines))
	}
}

func TestWorkoutHandler_GetRoutineWithExercises(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	routine, _ := userService.CreateRoutine(1, routineReq)

	req := createRequestWithContext("GET", "/routines/1", "1", nil)
	w := httptest.NewRecorder()

	workoutHandler.GetRoutineWithExercises(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.ExerciseRoutine
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, routine.Name, response.Name, "routine name")
}

func TestWorkoutHandler_AddExerciseToRoutine(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	userService.CreateRoutine(1, routineReq)

	req := createRequestWithContext("POST", "/routines/1/exercises?exercise_id=1", "1", nil)
	w := httptest.NewRecorder()

	workoutHandler.AddExerciseToRoutine(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, "Exercise added to routine successfully", response.Message, "message")
}

func TestWorkoutHandler_RemoveExerciseFromRoutine(t *testing.T) {
	workoutHandler, userService := setupWorkoutHandler()
	createTestUser(userService)

	routineReq := createTestRoutine()
	userService.CreateRoutine(1, routineReq)
	userService.AddExerciseToRoutine(1, 1)

	req := createRequestWithContext("DELETE", "/routines/1/exercises/1", "1", nil)
	// Add exercise_id to URL params
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	rctx.URLParams.Add("exercise_id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	workoutHandler.RemoveExerciseFromRoutine(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response model.BasicResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assertResponseField(t, "Exercise removed from routine successfully", response.Message, "message")
}