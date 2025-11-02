package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

func TestAchievementHandler_Create_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/achievements", mustJSONString(t, "{"))
	h.CreateAchievement(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status=%d want=400", w.Code)
	}
}

func TestAchievementHandler_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	req := model.CreateAchievementRequest{UserID: 1, Title: "T"}
	svc.EXPECT().CreateAchievement(gomock.AssignableToTypeOf(model.CreateAchievementRequest{})).
		Return(nil, errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/achievements", mustJSON(t, req))
	h.CreateAchievement(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want=500", w.Code)
	}
}

func TestAchievementHandler_Create_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	req := model.CreateAchievementRequest{UserID: 1, Title: "T"}
	want := &model.Achievement{ID: 9, UserID: 1, Title: "T"}

	svc.EXPECT().CreateAchievement(gomock.AssignableToTypeOf(model.CreateAchievementRequest{})).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/achievements", mustJSON(t, req))
	h.CreateAchievement(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d want=201", w.Code)
	}
	var got model.Achievement
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.ID != want.ID || got.Title != want.Title {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestAchievementHandler_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	svc.EXPECT().ReadAchievements().Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	h.ReadAchievements(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want=500", w.Code)
	}
}

func TestAchievementHandler_Read_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	want := []*model.Achievement{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}}
	svc.EXPECT().ReadAchievements().Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	h.ReadAchievements(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=200", w.Code)
	}
	var got []model.Achievement
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != len(want) {
		t.Fatalf("len=%d want=%d", len(got), len(want))
	}
}

func TestAchievementHandler_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	const id int64 = 7
	svc.EXPECT().DeleteAchievement(id).Return(errors.New("nope"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/achievements/7", nil)
	r = withIDCtx(r, id)

	h.DeleteAchievement(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want=500", w.Code)
	}
}

func TestAchievementHandler_Delete_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	const id int64 = 8
	svc.EXPECT().DeleteAchievement(id).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/achievements/8", nil)
	r = withIDCtx(r, id)

	h.DeleteAchievement(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=200", w.Code)
	}
	var br model.BasicResponse
	_ = json.NewDecoder(w.Body).Decode(&br)
	if br.Message != "Success" {
		t.Fatalf("msg=%q", br.Message)
	}
}
