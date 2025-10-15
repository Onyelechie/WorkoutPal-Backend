package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"workoutpal/src/internal/model"
	"workoutpal/src/util/constants"

	mock_service "workoutpal/src/mock_internal/domain/service"

	"github.com/golang/mock/gomock"
)

func mustJSON[T any](t *testing.T, v T) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func TestExerciseHandler_ReadExercises_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockExerciseService(ctrl)
	h := &exerciseHandler{exerciseService: mockSvc}

	want := []*model.Exercise{
		{ID: 1, Name: "Squat"},
		{ID: 2, Name: "Bench Press"},
	}

	mockSvc.EXPECT().
		ReadAllExercises().
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercises", nil)

	h.ReadExercises(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var got []model.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0].Name != want[0].Name {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestExerciseHandler_ReadExercises_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockExerciseService(ctrl)
	h := &exerciseHandler{exerciseService: mockSvc}

	mockSvc.EXPECT().
		ReadAllExercises().
		Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercises", nil)

	h.ReadExercises(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "boom" {
		t.Fatalf("message = %q, want %q", br.Message, "boom")
	}
}

func TestExerciseHandler_ReadExerciseByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockExerciseService(ctrl)
	h := &exerciseHandler{exerciseService: mockSvc}

	const id int64 = 42
	want := &model.Exercise{ID: id, Name: "Deadlift"}

	mockSvc.EXPECT().
		ReadExerciseByID(id).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercises/42", nil)
	ctx := context.WithValue(r.Context(), constants.ID_KEY, id)
	r = r.WithContext(ctx)

	h.ReadExerciseByID(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got model.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.Name != want.Name {
		t.Fatalf("unexpected payload: %+v", got)
	}
}

func TestExerciseHandler_ReadExerciseByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockExerciseService(ctrl)
	h := &exerciseHandler{exerciseService: mockSvc}

	const id int64 = 7
	mockSvc.EXPECT().
		ReadExerciseByID(id).
		Return(&model.Exercise{}, errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/exercises/7", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, id))

	h.ReadExerciseByID(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "not found" {
		t.Fatalf("message = %q, want %q", br.Message, "not found")
	}
}
