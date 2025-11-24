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

func TestPostHandler_CreatePost_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSONString(t, "{"))

	ctx := context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42))
	r = r.WithContext(ctx)

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

	userID := int64(42)
	req := model.CreatePostRequest{Title: "Test", Body: "Body"}

	svc.EXPECT().
		CreatePost(gomock.AssignableToTypeOf(model.CreatePostRequest{})).
		Return(nil, errors.New("service error"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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

	userID := int64(42)
	req := model.CreatePostRequest{Title: "New Post", Body: "Hello"} // client PostedBy ignored
	want := &model.Post{ID: 10, Title: req.Title, Body: req.Body}

	svc.EXPECT().
		CreatePost(gomock.AssignableToTypeOf(model.CreatePostRequest{})).
		DoAndReturn(func(r model.CreatePostRequest) (*model.Post, error) {
			// ensure handler overwrote PostedBy with context user
			if r.PostedBy != userID {
				t.Fatalf("expected PostedBy=%d, got %d", userID, r.PostedBy)
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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

	userID := int64(42)

	svc.EXPECT().ReadPosts(userID).Return(nil, errors.New("read failed"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/posts", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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

	userID := int64(42)
	want := []*model.Post{{ID: 1, Title: "A"}, {ID: 2, Title: "B"}}

	svc.EXPECT().ReadPosts(userID).Return(want, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/posts", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

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

	userID := int64(42)
	req := model.CommentOnPostRequest{PostID: 1, Comment: "Nice"} // client UserID ignored

	svc.EXPECT().
		CommentOnPost(gomock.AssignableToTypeOf(model.CommentOnPostRequest{})).
		Return(errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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

	userID := int64(42)
	req := model.CommentOnPostRequest{PostID: 1, Comment: "Nice"}

	svc.EXPECT().
		CommentOnPost(gomock.AssignableToTypeOf(model.CommentOnPostRequest{})).
		DoAndReturn(func(r model.CommentOnPostRequest) error {
			if r.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, r.UserID)
			}
			return nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

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

	userID := int64(42)
	req := model.CommentOnCommentRequest{CommentID: 1, PostID: 1, Comment: "Reply"}

	svc.EXPECT().
		CommentOnComment(gomock.AssignableToTypeOf(model.CommentOnCommentRequest{})).
		Return(errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment/reply", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

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

	userID := int64(42)
	req := model.CommentOnCommentRequest{CommentID: 1, PostID: 1, Comment: "Reply"}

	svc.EXPECT().
		CommentOnComment(gomock.AssignableToTypeOf(model.CommentOnCommentRequest{})).
		DoAndReturn(func(r model.CommentOnCommentRequest) error {
			if r.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, r.UserID)
			}
			return nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/comment/reply", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.CommentOnComment(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

// ==== Like / Unlike handlers ====

func TestPostHandler_LikePost_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/like", mustJSONString(t, "{"))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.LikePost(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestPostHandler_LikePost_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	userID := int64(42)
	req := model.LikePostRequest{PostID: 1}

	svc.EXPECT().
		LikePost(gomock.AssignableToTypeOf(model.LikePostRequest{})).
		Return((*model.Post)(nil), errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/like", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.LikePost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_LikePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	userID := int64(42)
	req := model.LikePostRequest{PostID: 1}
	want := &model.Post{ID: 1, IsLiked: true}

	svc.EXPECT().
		LikePost(gomock.AssignableToTypeOf(model.LikePostRequest{})).
		DoAndReturn(func(r model.LikePostRequest) (*model.Post, error) {
			if r.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, r.UserID)
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/like", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.LikePost(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestPostHandler_UnlikePost_BadJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/unlike", mustJSONString(t, "{"))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, int64(42)))

	h.UnlikePost(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestPostHandler_UnlikePost_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	userID := int64(42)
	req := model.UnikePostRequest{PostID: 1}

	svc.EXPECT().
		UnlikePost(gomock.AssignableToTypeOf(model.UnikePostRequest{})).
		Return((*model.Post)(nil), errors.New("fail"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/unlike", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.UnlikePost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}

func TestPostHandler_UnlikePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	userID := int64(42)
	req := model.UnikePostRequest{PostID: 1}
	want := &model.Post{ID: 1, IsLiked: false}

	svc.EXPECT().
		UnlikePost(gomock.AssignableToTypeOf(model.UnikePostRequest{})).
		DoAndReturn(func(r model.UnikePostRequest) (*model.Post, error) {
			if r.UserID != userID {
				t.Fatalf("expected UserID=%d, got %d", userID, r.UserID)
			}
			return want, nil
		})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/posts/unlike", mustJSON(t, req))
	r = r.WithContext(context.WithValue(r.Context(), constants.USER_ID_KEY, userID))

	h.UnlikePost(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

// ==== DeletePost handler ====

func TestPostHandler_DeletePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	postID := int64(10)
	svc.EXPECT().DeletePost(postID).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/posts/10", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, postID))

	h.DeletePost(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
}

func TestPostHandler_DeletePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	svc := mock_service.NewMockPostService(ctrl)
	h := &PostHandler{svc: svc}

	postID := int64(10)
	svc.EXPECT().DeletePost(postID).Return(errors.New("not found"))

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/posts/10", nil)
	r = r.WithContext(context.WithValue(r.Context(), constants.ID_KEY, postID))

	h.DeletePost(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}
