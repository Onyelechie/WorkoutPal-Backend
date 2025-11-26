package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
	"workoutpal/src/util/constants"

	"github.com/golang/mock/gomock"
)

func TestExerciseSettingHandler_UpdateExerciseSetting_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/exercise-settings", mustJSONString(t, "{"))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.UpdateExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_UpdateExerciseSetting_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	body := model.UpdateExerciseSettingRequest{
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}

	svc.EXPECT().
		UpdateExerciseSetting(gomock.AssignableToTypeOf(model.UpdateExerciseSettingRequest{})).
		Return((*model.ExerciseSetting)(nil), errors.New("service fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/exercise-settings", mustJSON(t, body))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.UpdateExerciseSetting(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestExerciseSettingHandler_UpdateExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	body := model.UpdateExerciseSettingRequest{
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}
	want := &model.ExerciseSetting{
		UserID:           userID,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}

	svc.EXPECT().
		UpdateExerciseSetting(gomock.AssignableToTypeOf(model.UpdateExerciseSettingRequest{})).
		DoAndReturn(func(req model.UpdateExerciseSettingRequest) (*model.ExerciseSetting, error) {
			if req.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, req.UserID)
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/exercise-settings", mustJSON(t, body))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.UpdateExerciseSetting(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var got model.ExerciseSetting
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.UserID != want.UserID ||
		got.ExerciseID != want.ExerciseID ||
		got.WorkoutRoutineID != want.WorkoutRoutineID {
		t.Fatalf("unexpected setting: %#v", got)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_MissingExerciseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?workout_routine_id=3", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_InvalidExerciseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?exercise_id=abc&workout_routine_id=3", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_MissingWorkoutRoutineID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?exercise_id=2", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_InvalidWorkoutRoutineID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?exercise_id=2&workout_routine_id=xyz", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	req := model.ReadExerciseSettingRequest{
		UserID:           userID,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	svc.EXPECT().
		ReadExerciseSetting(req).
		Return((*model.ExerciseSetting)(nil), errors.New("service fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?exercise_id=2&workout_routine_id=3", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestExerciseSettingHandler_ReadExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	req := model.ReadExerciseSettingRequest{
		UserID:           userID,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}
	want := &model.ExerciseSetting{
		UserID:           userID,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           60,
		Reps:             10,
		Sets:             4,
		BreakInterval:    90,
	}

	svc.EXPECT().
		ReadExerciseSetting(req).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercise-settings?exercise_id=2&workout_routine_id=3", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.ReadExerciseSetting(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var got model.ExerciseSetting
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.UserID != want.UserID || got.ExerciseID != want.ExerciseID || got.WorkoutRoutineID != want.WorkoutRoutineID {
		t.Fatalf("unexpected setting: %#v", got)
	}
}

// ==== CreateExerciseSetting tests ====

func TestExerciseSettingHandler_CreateExerciseSetting_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/exercise-settings", mustJSONString(t, "{"))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.CreateExerciseSetting(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestExerciseSettingHandler_CreateExerciseSetting_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	body := model.CreateExerciseSettingRequest{
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           50,
		Reps:             8,
		Sets:             3,
		BreakInterval:    60,
	}

	svc.EXPECT().
		CreateExerciseSetting(gomock.AssignableToTypeOf(model.CreateExerciseSettingRequest{})).
		Return((*model.ExerciseSetting)(nil), errors.New("service fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/exercise-settings", mustJSON(t, body))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.CreateExerciseSetting(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestExerciseSettingHandler_CreateExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockExerciseSettingService(ctrl)
	h := NewExerciseSettingHandler(svc)

	userID := int64(42)
	body := model.CreateExerciseSettingRequest{
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           50,
		Reps:             8,
		Sets:             3,
		BreakInterval:    60,
	}
	want := &model.ExerciseSetting{
		UserID:           userID,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           50,
		Reps:             8,
		Sets:             3,
		BreakInterval:    60,
	}

	svc.EXPECT().
		CreateExerciseSetting(gomock.AssignableToTypeOf(model.CreateExerciseSettingRequest{})).
		DoAndReturn(func(req model.CreateExerciseSettingRequest) (*model.ExerciseSetting, error) {
			if req.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, req.UserID)
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/exercise-settings", mustJSON(t, body))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.CreateExerciseSetting(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}

	var got model.ExerciseSetting
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.ExerciseID != want.ExerciseID || got.WorkoutRoutineID != want.WorkoutRoutineID || got.UserID != want.UserID {
		t.Fatalf("unexpected setting: %#v", got)
	}
}
