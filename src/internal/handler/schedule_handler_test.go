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
	mock_service "workoutpal/src/mock_internal/domain/service"
	"workoutpal/src/util/constants"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
)

func mustJSONBuf[T any](t *testing.T, v T) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func TestScheduleHandler_ReadUserSchedules_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(99)

	want := []*model.Schedule{
		{ID: 1, Name: "Push Day", UserID: userID},
		{ID: 2, Name: "Leg Day", UserID: userID},
	}

	mockSvc.EXPECT().
		ReadUserSchedules(userID).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/me", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, userID))

	h.ReadUserSchedules(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var got []model.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0].Name != want[0].Name {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestScheduleHandler_ReadUserSchedules_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(99)

	mockSvc.EXPECT().
		ReadUserSchedules(userID).
		Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/me", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, userID))

	h.ReadUserSchedules(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br map[string]string
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br["error"] != "boom" {
		t.Fatalf("message = %q, want %q", br["error"], "boom")
	}
}

func TestScheduleHandler_ReadUserSchedulesByDay_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(33)
	day := int64(2)

	want := []*model.Schedule{
		{ID: 10, Name: "Legs", UserID: userID, DayOfWeek: day},
	}

	mockSvc.EXPECT().
		ReadUserSchedulesByDay(userID, day).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/me/2", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))
	r = muxWithParam(r, constants.DAY_OF_WEEK_KEY, "2")

	h.ReadUserSchedulesByDay(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var got []model.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Legs" {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestScheduleHandler_ReadUserSchedulesByDay_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(33)
	day := int64(4)

	mockSvc.EXPECT().
		ReadUserSchedulesByDay(userID, day).
		Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/me/4", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))
	r = muxWithParam(r, constants.DAY_OF_WEEK_KEY, "4")

	h.ReadUserSchedulesByDay(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br map[string]string
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br["error"] != "boom" {
		t.Fatalf("message = %q, want %q", br["error"], "boom")
	}
}

func TestScheduleHandler_ReadScheduleByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	id := int64(55)
	want := &model.Schedule{ID: id, Name: "Pull Day"}

	mockSvc.EXPECT().
		ReadScheduleByID(id).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/55", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, id))

	h.ReadScheduleByID(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got model.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != id || got.Name != "Pull Day" {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestScheduleHandler_ReadScheduleByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	id := int64(123)

	mockSvc.EXPECT().
		ReadScheduleByID(id).
		Return(nil, errors.New("nope"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/schedules/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, id))

	h.ReadScheduleByID(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br map[string]string
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br["error"] != "nope" {
		t.Fatalf("message = %q, want %q", br["error"], "nope")
	}
}

func TestScheduleHandler_CreateSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(77)

	inReq := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               0,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}
	expectedReq := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               userID,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}
	outSchedule := &model.Schedule{
		ID:                   999,
		Name:                 "AM Upper",
		UserID:               userID,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}

	mockSvc.EXPECT().
		CreateSchedule(expectedReq).
		Return(outSchedule, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/schedules", mustJSONBuf(t, inReq))
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, userID))

	h.CreateSchedule(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", resp.StatusCode)
	}

	var got model.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != outSchedule.ID || got.Name != outSchedule.Name {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestScheduleHandler_CreateSchedule_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(77)

	inReq := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               0,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}
	expectedReq := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               userID,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}

	mockSvc.EXPECT().
		CreateSchedule(expectedReq).
		Return(nil, errors.New("explode"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/schedules", mustJSONBuf(t, inReq))
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, userID))

	h.CreateSchedule(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br map[string]string
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br["error"] != "explode" {
		t.Fatalf("message = %q, want %q", br["error"], "explode")
	}
}

func TestScheduleHandler_UpdateSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(77)
	id := int64(123)

	inReq := model.UpdateScheduleRequest{
		ID:                   0,
		UserID:               0,
		Name:                 "PM Lower",
		DayOfWeek:            4,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}
	expectedReq := model.UpdateScheduleRequest{
		ID:                   id,
		UserID:               userID,
		Name:                 "PM Lower",
		DayOfWeek:            4,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}
	outSchedule := &model.Schedule{
		ID:                   id,
		UserID:               userID,
		Name:                 "PM Lower",
		DayOfWeek:            4,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}

	mockSvc.EXPECT().
		UpdateSchedule(expectedReq).
		Return(outSchedule, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/schedules/123", mustJSONBuf(t, inReq))
	ctx := context.WithValue(r.Context(), constants.ID_KEY, id)
	ctx = context.WithValue(ctx, constants.USER_ID_KEY, userID)
	r = r.WithContext(ctx)

	h.UpdateSchedule(w, r)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var got model.Schedule
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ID != id || got.Name != "PM Lower" {
		t.Fatalf("unexpected payload: %#v", got)
	}
}

func TestScheduleHandler_UpdateSchedule_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	userID := int64(77)
	id := int64(123)

	inReq := model.UpdateScheduleRequest{
		ID:                   0,
		UserID:               0,
		Name:                 "PM Lower",
		DayOfWeek:            4,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}
	expectedReq := model.UpdateScheduleRequest{
		ID:                   id,
		UserID:               userID,
		Name:                 "PM Lower",
		DayOfWeek:            4,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}

	mockSvc.EXPECT().
		UpdateSchedule(expectedReq).
		Return(nil, errors.New("bad update"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/schedules/123", mustJSONBuf(t, inReq))
	ctx := context.WithValue(r.Context(), constants.ID_KEY, id)
	ctx = context.WithValue(ctx, constants.USER_ID_KEY, userID)
	r = r.WithContext(ctx)

	h.UpdateSchedule(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
	var br map[string]string
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br["error"] != "bad update" {
		t.Fatalf("message = %q, want %q", br["error"], "bad update")
	}
}

func TestScheduleHandler_DeleteSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockScheduleService(ctrl)
	h := &scheduleHandler{service: mockSvc}

	id := int64(123)

	mockSvc.EXPECT().
		DeleteSchedule(model.DeleteScheduleRequest{ID: id}).
		Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/schedules/123", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, id))

	h.DeleteSchedule(w, r)
}

func muxWithParam(r *http.Request, key, val string) *http.Request {
	rc := chi.NewRouteContext()
	rc.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rc))
}
