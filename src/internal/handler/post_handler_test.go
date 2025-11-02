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

func TestPostHandler_CreatePost_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSONString(t, "{"))
	h.CreatePost(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestPostHandler_CreatePost_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CreatePostRequest{Title: "Test", Body: "Body"}
	svc.EXPECT().CreatePost(gomock.AssignableToTypeOf(model.CreatePostRequest{})).
		Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSON(t, req))
	h.CreatePost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_CreatePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CreatePostRequest{Title: "New Post", Body: "Hello", PostedBy: 1}
	want := &model.Post{ID: 10, Title: req.Title, Body: req.Body}

	svc.EXPECT().CreatePost(gomock.AssignableToTypeOf(model.CreatePostRequest{})).
		Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSON(t, req))
	h.CreatePost(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}

	var got model.Post
	_ = json.NewDecoder(w.Body).Decode(&got)
	if got.ID != want.ID || got.Title != want.Title {
		t.Fatalf("unexpected post: %+v", got)
	}
}

func TestPostHandler_ReadPosts_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	svc.EXPECT().ReadPosts().Return(nil, errors.New("read failed"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/posts", nil)
	h.ReadPosts(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_ReadPosts_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	want := []*model.Post{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}}
	svc.EXPECT().ReadPosts().Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/posts", nil)
	h.ReadPosts(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var got []model.Post
	_ = json.NewDecoder(w.Body).Decode(&got)
	if len(got) != len(want) {
		t.Fatalf("unexpected posts length: %+v", got)
	}
}

func TestPostHandler_CommentOnPost_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment", mustJSONString(t, "{"))
	h.CommentOnPost(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestPostHandler_CommentOnPost_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CommentOnPostRequest{PostID: 1, UserID: 1, Comment: "Nice"}
	svc.EXPECT().CommentOnPost(req).Return(errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment", mustJSON(t, req))
	h.CommentOnPost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_CommentOnPost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CommentOnPostRequest{PostID: 1, UserID: 1, Comment: "Nice"}
	svc.EXPECT().CommentOnPost(req).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment", mustJSON(t, req))
	h.CommentOnPost(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestPostHandler_CommentOnComment_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment/reply", mustJSONString(t, "{"))
	h.CommentOnComment(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestPostHandler_CommentOnComment_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CommentOnCommentRequest{CommentID: 1, PostID: 1, UserID: 1, Comment: "Reply"}
	svc.EXPECT().CommentOnComment(req).Return(errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment/reply", mustJSON(t, req))
	h.CommentOnComment(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_CommentOnComment_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	req := model.CommentOnCommentRequest{CommentID: 1, PostID: 1, UserID: 1, Comment: "Reply"}
	svc.EXPECT().CommentOnComment(req).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment/reply", mustJSON(t, req))
	h.CommentOnComment(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}
