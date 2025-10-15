package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"golang.org/x/net/context"

	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

func setChiURLParam(r *http.Request, key, val string) *http.Request {
	rc := chi.NewRouteContext()
	rc.URLParams.Add(key, val)
	return r.WithContext(contextWithRouteCtx(r.Context(), rc))
}

func contextWithRouteCtx(ctx context.Context, rc *chi.Context) context.Context {
	type ctxKey struct{}
	return context.WithValue(ctx, chi.RouteCtxKey, rc)
}

func TestGoalHandler_CreateUserGoal_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/abc/goals", bytes.NewBufferString(`{}`))
	r = setChiURLParam(r, "id", "abc")

	h.CreateUserGoal(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGoalHandler_CreateUserGoal_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/1/goals", bytes.NewBufferString("{"))
	r.Header.Set("Content-Type", "application/json")
	r = setChiURLParam(r, "id", "1")

	h.CreateUserGoal(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGoalHandler_CreateUserGoal_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	const userID int64 = 1
	req := model.CreateGoalRequest{Name: "PR Deadlift", Description: "Description"}

	mockSvc.
		EXPECT().
		CreateGoal(userID, gomock.AssignableToTypeOf(model.CreateGoalRequest{})).
		Return(&model.Goal{}, errors.New("validation failed"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/1/goals", mustJSONBody(t, req))
	r.Header.Set("Content-Type", "application/json")
	r = setChiURLParam(r, "id", "1")

	h.CreateUserGoal(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	var br struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "validation failed" {
		t.Fatalf("message = %q, want %q", br.Message, "validation failed")
	}
}

func TestGoalHandler_CreateUserGoal_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	const userID int64 = 2
	req := model.CreateGoalRequest{Name: "PR Deadlift", Description: "Description"}
	want := &model.Goal{ID: 1, Name: "PR Deadlift", Description: "Description", UserID: 2}

	mockSvc.
		EXPECT().
		CreateGoal(userID, gomock.AssignableToTypeOf(model.CreateGoalRequest{})).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/2/goals", mustJSONBody(t, req))
	r.Header.Set("Content-Type", "application/json")
	r = setChiURLParam(r, "id", "2")

	h.CreateUserGoal(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}
	var got model.Goal
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.UserID != userID || got.Name != want.Name {
		t.Fatalf("unexpected goal: %+v", got)
	}
}

func TestGoalHandler_GetUserGoals_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/xyz/goals", nil)
	r = setChiURLParam(r, "id", "xyz")

	h.GetUserGoals(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestGoalHandler_GetUserGoals_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	const userID int64 = 3
	mockSvc.
		EXPECT().
		ReadUserGoals(userID).
		Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/3/goals", nil)
	r = setChiURLParam(r, "id", "3")

	h.GetUserGoals(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	var br struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "user not found" {
		t.Fatalf("message = %q, want %q", br.Message, "user not found")
	}
}

func TestGoalHandler_GetUserGoals_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockGoalService(ctrl)
	h := &goalHandler{goalService: mockSvc}

	const userID int64 = 4
	want := []*model.Goal{
		{Name: "PR Deadlift", Description: "Description"},
		{Name: "PR Deadlift2", Description: "Description2", UserID: 4},
	}

	mockSvc.
		EXPECT().
		ReadUserGoals(userID).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/4/goals", nil)
	r = setChiURLParam(r, "id", "4")

	h.GetUserGoals(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got []*model.Goal
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Name != want[1].Name {
		t.Fatalf("unexpected goals: %+v", got)
	}
}
