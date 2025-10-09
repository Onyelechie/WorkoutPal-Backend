package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/domain/handler"
	handler_impl "workoutpal/src/internal/handler"
	"workoutpal/src/internal/model"
	"workoutpal/src/internal/repository"
	service_impl "workoutpal/src/internal/service"
)

func setupExerciseHandler() handler.ExerciseHandler {
	repo := repository.NewInMemoryExerciseRepository()
	exerciseService := service_impl.NewExerciseService(repo)
	return handler_impl.NewExerciseHandler(exerciseService)
}

func TestExerciseHandler_ReadExercises(t *testing.T) {
	exerciseHandler := setupExerciseHandler()

	req := createRequestWithContext("GET", "/exercises", "", nil)
	w := httptest.NewRecorder()

	exerciseHandler.ReadExercises(w, req)

	assertStatusCode(t, http.StatusOK, w.Code)

	var response []model.Exercise
	json.Unmarshal(w.Body.Bytes(), &response)
	if len(response) == 0 {
		t.Error("Expected exercises to be returned")
	}
	if response[0].Name == "" {
		t.Error("Expected exercise to have a name")
	}
}