package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"

	"workoutpal/src/internal/model"
	mock_service "workoutpal/src/mock_internal/domain/service"
)

/* ----------------------------- helpers ----------------------------- */

func withChiURLParam(r *http.Request, key, val string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, val)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, routeCtx))
}

func decodeMessage(t *testing.T, rr *httptest.ResponseRecorder) string {
	t.Helper()
	var br struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(rr.Body).Decode(&br)
	return br.Message
}

/* ------------------------------- tests ------------------------------ */

func TestRelationshipHandler_ReadFollowers_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/abc/followers", nil)
	r = withChiURLParam(r, "id", "abc")

	h.ReadFollowers(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRelationshipHandler_ReadFollowers_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	const userID int64 = 1
	mockSvc.EXPECT().ReadUserFollowers(userID).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/1/followers", nil)
	r = withChiURLParam(r, "id", "1")

	h.ReadFollowers(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "user not found" {
		t.Fatalf("message = %q, want %q", msg, "user not found")
	}
}

func TestRelationshipHandler_ReadFollowers_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	const userID int64 = 2
	want := []int64{5, 7, 9}
	mockSvc.EXPECT().ReadUserFollowers(userID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/2/followers", nil)
	r = withChiURLParam(r, "id", "2")

	h.ReadFollowers(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got []int64
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[0] != want[0] || got[2] != want[2] {
		t.Fatalf("unexpected followers: %+v", got)
	}
}

func TestRelationshipHandler_ReadFollowings_BadID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/zzz/following", nil)
	r = withChiURLParam(r, "id", "zzz")

	h.ReadFollowings(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRelationshipHandler_ReadFollowings_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	const userID int64 = 3
	mockSvc.EXPECT().ReadUserFollowing(userID).Return(nil, errors.New("user not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/3/following", nil)
	r = withChiURLParam(r, "id", "3")

	h.ReadFollowings(w, r)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "user not found" {
		t.Fatalf("message = %q, want %q", msg, "user not found")
	}
}

func TestRelationshipHandler_ReadFollowings_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	const userID int64 = 4
	want := []int64{11, 12}
	mockSvc.EXPECT().ReadUserFollowing(userID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/users/4/following", nil)
	r = withChiURLParam(r, "id", "4")

	h.ReadFollowings(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var got []int64
	if err := json.NewDecoder(w.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(got) != len(want) || got[1] != want[1] {
		t.Fatalf("unexpected following: %+v", got)
	}
}

func TestRelationshipHandler_FollowUser_BadFolloweeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/abc/follow?follower_id=1", nil)
	r = withChiURLParam(r, "id", "abc")

	h.FollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRelationshipHandler_FollowUser_BadFollowerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/2/follow?follower_id=bad", nil)
	r = withChiURLParam(r, "id", "2")

	h.FollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "Invalid follower ID" {
		t.Fatalf("message = %q, want %q", msg, "Invalid follower ID")
	}
}

func TestRelationshipHandler_FollowUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("follower_id", "1")
	r := httptest.NewRequest(http.MethodPost, "/users/2/follow?"+q.Encode(), nil)
	r = withChiURLParam(r, "id", "2")

	mockSvc.EXPECT().FollowUser(int64(1), int64(2)).Return(errors.New("already following"))

	h.FollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "already following" {
		t.Fatalf("message = %q, want %q", msg, "already following")
	}
}

func TestRelationshipHandler_FollowUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("follower_id", "3")
	r := httptest.NewRequest(http.MethodPost, "/users/5/follow?"+q.Encode(), nil)
	r = withChiURLParam(r, "id", "5")

	mockSvc.EXPECT().FollowUser(int64(3), int64(5)).Return(nil)

	h.FollowUser(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "Successfully followed user" {
		t.Fatalf("message = %q, want %q", br.Message, "Successfully followed user")
	}
}

func TestRelationshipHandler_UnfollowUser_BadFolloweeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/abc/unfollow?follower_id=1", nil)
	r = withChiURLParam(r, "id", "abc")

	h.UnfollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestRelationshipHandler_UnfollowUser_BadFollowerID(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/users/2/unfollow?follower_id=x", nil)
	r = withChiURLParam(r, "id", "2")

	h.UnfollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "Invalid follower ID" {
		t.Fatalf("message = %q, want %q", msg, "Invalid follower ID")
	}
}

func TestRelationshipHandler_UnfollowUser_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("follower_id", "8")
	r := httptest.NewRequest(http.MethodPost, "/users/9/unfollow?"+q.Encode(), nil)
	r = withChiURLParam(r, "id", "9")

	mockSvc.EXPECT().UnfollowUser(int64(8), int64(9)).Return(errors.New("not following"))

	h.UnfollowUser(w, r)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if msg := decodeMessage(t, w); msg != "not following" {
		t.Fatalf("message = %q, want %q", msg, "not following")
	}
}

func TestRelationshipHandler_UnfollowUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockSvc := mock_service.NewMockRelationshipService(ctrl)
	h := &relationshipHandler{relationshipService: mockSvc}

	w := httptest.NewRecorder()
	q := url.Values{}
	q.Set("follower_id", "10")
	r := httptest.NewRequest(http.MethodPost, "/users/11/unfollow?"+q.Encode(), nil)
	r = withChiURLParam(r, "id", "11")

	mockSvc.EXPECT().UnfollowUser(int64(10), int64(11)).Return(nil)

	h.UnfollowUser(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var br model.BasicResponse
	if err := json.NewDecoder(w.Body).Decode(&br); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if br.Message != "Successfully unfollowed user" {
		t.Fatalf("message = %q, want %q", br.Message, "Successfully unfollowed user")
	}
}
