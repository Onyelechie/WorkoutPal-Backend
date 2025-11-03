package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
	"workoutpal/src/util/constants"
)

/* ---------- CreateAchievement ---------- */

func TestAchievementHandler_Create_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/achievements", mustJSONString(t, "{")) // invalid JSON
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

	req := model.CreateAchievementRequest{UserID: 1, AchievementID: 55}
	svc.EXPECT().CreateAchievement(req).Return(nil, errors.New("fail"))

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

	req := model.CreateAchievementRequest{UserID: 1, AchievementID: 55}
	want := &model.UserAchievement{ID: 9, UserID: 1, Title: "First Workout", EarnedAt: "2025-01-01T00:00:00Z"}

	svc.EXPECT().CreateAchievement(req).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/achievements", mustJSON(t, req))
	h.CreateAchievement(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status=%d want=201", w.Code)
	}
	var got model.UserAchievement
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.ID != want.ID || got.UserID != want.UserID || got.Title != want.Title || got.EarnedAt != want.EarnedAt {
		t.Fatalf("unexpected body: %#v", got)
	}
}

/* ---------- ReadAllAchievements (catalog) ---------- */

func TestAchievementHandler_ReadAllAchievements_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	want := []*model.Achievement{
		{ID: 1, Title: "First Workout", BadgeIcon: "first.png", Description: "Finish your first workout"},
		{ID: 2, Title: "7-Day Streak", BadgeIcon: "streak7.png", Description: "Train 7 days in a row"},
	}

	svc.EXPECT().ReadAllAchievements().Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	h.ReadAllAchievements(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=200", w.Code)
	}
	var got []*model.Achievement
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 || got[0].ID != 1 || got[1].ID != 2 {
		t.Fatalf("unexpected body: %#v", got)
	}
}

func TestAchievementHandler_ReadAllAchievements_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	svc.EXPECT().ReadAllAchievements().Return(nil, errors.New("boom"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements", nil)
	h.ReadAllAchievements(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want=500", w.Code)
	}
}

/* ---------- ReadUnlockedAchievements (current user) ---------- */

func TestAchievementHandler_ReadUnlockedAchievements_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	const userID int64 = 7
	want := []*model.UserAchievement{
		{ID: 10, UserID: userID, Title: "First Workout", EarnedAt: "2025-01-01T00:00:00Z"},
		{ID: 11, UserID: userID, Title: "7-Day Streak", EarnedAt: "2025-01-05T00:00:00Z"},
	}

	svc.EXPECT().ReadUnlockedAchievements(userID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements/unlocked", nil)
	r = withUserCtx(r, userID)

	h.ReadUnlockedAchievements(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status=%d want=200", w.Code)
	}
	var got []*model.UserAchievement
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != 2 || got[0].UserID != userID {
		t.Fatalf("unexpected body: %#v", got)
	}
}

func TestAchievementHandler_ReadUnlockedAchievements_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	svc := mock_service.NewMockAchievementService(ctrl)
	h := &AchievementHandler{svc: svc}

	const userID int64 = 7
	svc.EXPECT().ReadUnlockedAchievements(userID).Return(nil, errors.New("nope"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/achievements/unlocked", nil)
	r = withUserCtx(r, userID)

	h.ReadUnlockedAchievements(w, r)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status=%d want=500", w.Code)
	}
}

func withUserCtx(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), constants.USER_ID_KEY, userID)
	return r.WithContext(ctx)
}
