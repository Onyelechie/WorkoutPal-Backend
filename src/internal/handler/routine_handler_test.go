package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

func mustJSONString(t *testing.T, s string) *strings.Reader {
	t.Helper()
	return strings.NewReader(s)
}

func TestRoutineHandler_CreateUserRoutine_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/1/routines", mustJSONString(t, "{"))
	r = withIDCtx(r, 1)

	h.CreateUserRoutine(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRoutineHandler_CreateUserRoutine_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const userID int64 = 1
	req := model.CreateRoutineRequest{Name: "Push Day"}
	svc.EXPECT().
		CreateRoutine(userID, gomock.AssignableToTypeOf(model.CreateRoutineRequest{})).
		Return((*model.ExerciseRoutine)(nil), errors.New("validation failed"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/1/routines", mustJSON(t, req))
	r = withIDCtx(r, userID)

	h.CreateUserRoutine(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRoutineHandler_CreateUserRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const userID int64 = 2
	req := model.CreateRoutineRequest{Name: "Leg Day"}
	want := &model.ExerciseRoutine{ID: 10, UserID: userID, Name: req.Name}

	svc.EXPECT().
		CreateRoutine(userID, gomock.AssignableToTypeOf(model.CreateRoutineRequest{})).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/2/routines", mustJSON(t, req))
	r = withIDCtx(r, userID)

	h.CreateUserRoutine(w, r)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}
	var got model.ExerciseRoutine
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.UserID != want.UserID || got.Name != want.Name {
		t.Fatalf("unexpected routine: %+v", got)
	}
}

func TestRoutineHandler_ReadUserRoutines_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const userID int64 = 3
	svc.EXPECT().ReadUserRoutines(userID).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/3/routines", nil)
	r = withIDCtx(r, userID)

	h.ReadUserRoutines(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRoutineHandler_ReadUserRoutines_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const userID int64 = 4
	want := []*model.ExerciseRoutine{
		{ID: 1, UserID: userID, Name: "A"},
		{ID: 2, UserID: userID, Name: "B"},
	}
	svc.EXPECT().ReadUserRoutines(userID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/4/routines", nil)
	r = withIDCtx(r, userID)

	h.ReadUserRoutines(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got []model.ExerciseRoutine
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Name != want[1].Name {
		t.Fatalf("unexpected routines: %+v", got)
	}
}

func TestRoutineHandler_DeleteRoutine_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 5
	svc.EXPECT().DeleteRoutine(routineID).Return(errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/routines/5", nil)
	r = withIDCtx(r, routineID)

	h.DeleteRoutine(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRoutineHandler_DeleteRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 6
	svc.EXPECT().DeleteRoutine(routineID).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/routines/6", nil)
	r = withIDCtx(r, routineID)

	h.DeleteRoutine(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "Routine deleted successfully" {
		t.Fatalf("message = %q, want %q", br.Message, "Routine deleted successfully")
	}
}

func TestRoutineHandler_ReadRoutineWithExercises_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 7
	svc.EXPECT().ReadRoutineWithExercises(routineID).Return((*model.ExerciseRoutine)(nil), errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/routines/7", nil)
	r = withIDCtx(r, routineID)

	h.ReadRoutineWithExercises(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
}

func TestRoutineHandler_ReadRoutineWithExercises_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 8
	want := &model.ExerciseRoutine{ID: routineID, Name: "Pull Day"}
	svc.EXPECT().ReadRoutineWithExercises(routineID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/routines/8", nil)
	r = withIDCtx(r, routineID)

	h.ReadRoutineWithExercises(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got model.ExerciseRoutine
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != want.ID || got.Name != want.Name {
		t.Fatalf("unexpected routine: %+v", got)
	}
}

func TestRoutineHandler_AddExerciseToRoutine_BadExerciseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/routines/1/exercises?exercise_id=bad", nil)
	r = withIDCtx(r, 1)

	h.AddExerciseToRoutine(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRoutineHandler_AddExerciseToRoutine_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 2
	const exerciseID int64 = 9

	svc.EXPECT().AddExerciseToRoutine(routineID, exerciseID).Return(errors.New("already added"))

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("exercise_id", strconv.FormatInt(exerciseID, 10))
	r := httptest.NewRequest(http.MethodPost, "/routines/2/exercises?"+q.Encode(), nil)
	r = withIDCtx(r, routineID)

	h.AddExerciseToRoutine(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "already added" {
		t.Fatalf("message = %q, want %q", br.Message, "already added")
	}
}

func TestRoutineHandler_AddExerciseToRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 3
	const exerciseID int64 = 10

	svc.EXPECT().AddExerciseToRoutine(routineID, exerciseID).Return(nil)

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("exercise_id", strconv.FormatInt(exerciseID, 10))
	r := httptest.NewRequest(http.MethodPost, "/routines/3/exercises?"+q.Encode(), nil)
	r = withIDCtx(r, routineID)

	h.AddExerciseToRoutine(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "Exercise added to routine successfully" {
		t.Fatalf("message = %q, want %q", br.Message, "Exercise added to routine successfully")
	}
}

func TestRoutineHandler_RemoveExerciseFromRoutine_BadExerciseID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/routines/1/exercises/bad", nil)
	r = withIDCtx(r, 1)
	r = withChiURLParam(r, "exercise_id", "bad")

	h.RemoveExerciseFromRoutine(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRoutineHandler_RemoveExerciseFromRoutine_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 4
	const exerciseID int64 = 12

	svc.EXPECT().RemoveExerciseFromRoutine(routineID, exerciseID).Return(errors.New("not in routine"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/routines/4/exercises/12", nil)
	r = withIDCtx(r, routineID)
	r = withChiURLParam(r, "exercise_id", "12")

	h.RemoveExerciseFromRoutine(w, r)
	if w.Code != http.StatusNotFound { // your handler maps svc error to 404 here
		t.Fatalf("status = %d, want 404", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "not in routine" {
		t.Fatalf("message = %q, want %q", br.Message, "not in routine")
	}
}

func TestRoutineHandler_RemoveExerciseFromRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const routineID int64 = 5
	const exerciseID int64 = 13

	svc.EXPECT().RemoveExerciseFromRoutine(routineID, exerciseID).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/routines/5/exercises/13", nil)
	r = withIDCtx(r, routineID)
	r = withChiURLParam(r, "exercise_id", "13")

	h.RemoveExerciseFromRoutine(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "Exercise removed from routine successfully" {
		t.Fatalf("message = %q, want %q", br.Message, "Exercise removed from routine successfully")
	}
}

func TestRoutineHandler_DeleteUserRoutine_UsesContextID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockRoutineService(ctrl)
	h := &workoutHandler{routineService: svc}

	const ctxID int64 = 77
	svc.EXPECT().DeleteRoutine(ctxID).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/users/1/routines/77", nil)
	r = withIDCtx(r, ctxID)

	h.DeleteUserRoutine(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}
