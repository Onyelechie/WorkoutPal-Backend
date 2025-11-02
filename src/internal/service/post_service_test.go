package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestPostService_CreatePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CreatePostRequest{Title: "Test", Body: "Body"}
	want := &model.Post{ID: 1, Title: "Test", Body: "Body"}

	repo.EXPECT().CreatePost(req).Return(want, nil)
	repo.EXPECT().ReadCommentsByPost(int64(1)).Return(nil, nil)

	got, err := svc.CreatePost(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Title != want.Title {
		t.Fatalf("unexpected post: %#v", got)
	}
}

func TestPostService_CreatePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CreatePostRequest{Title: "Fail"}
	repo.EXPECT().CreatePost(req).Return((*model.Post)(nil), errors.New("db error"))

	got, err := svc.CreatePost(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestPostService_ReadPosts_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	want := []*model.Post{{ID: 1}, {ID: 2}}
	repo.EXPECT().ReadPosts().Return(want, nil)
	repo.EXPECT().ReadCommentsByPost(int64(1)).Return(nil, nil)
	repo.EXPECT().ReadCommentsByPost(int64(2)).Return(nil, nil)

	got, err := svc.ReadPosts()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != len(want) {
		t.Fatalf("unexpected posts: %#v", got)
	}
}

func TestPostService_ReadPosts_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	repo.EXPECT().ReadPosts().Return(nil, errors.New("failed"))
	got, err := svc.ReadPosts()
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "failed" {
		t.Fatalf("expected failed, got %v", err)
	}
}

func TestPostService_UpdatePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.UpdatePostRequest{ID: 1, Title: "Updated"}
	want := &model.Post{ID: 1, Title: "Updated"}

	repo.EXPECT().UpdatePost(req).Return(want, nil)
	repo.EXPECT().ReadCommentsByPost(int64(1)).Return(nil, nil)

	got, err := svc.UpdatePost(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != want.ID {
		t.Fatalf("unexpected post: %#v", got)
	}
}

func TestPostService_UpdatePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.UpdatePostRequest{ID: 1}
	repo.EXPECT().UpdatePost(req).Return((*model.Post)(nil), errors.New("no post"))

	got, err := svc.UpdatePost(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "no post" {
		t.Fatalf("expected no post, got %v", err)
	}
}

func TestPostService_DeletePost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	repo.EXPECT().DeletePost(int64(1)).Return(nil)

	if err := svc.DeletePost(1); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostService_DeletePost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	repo.EXPECT().DeletePost(int64(1)).Return(errors.New("not found"))

	if err := svc.DeletePost(1); err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestPostService_CommentOnPost_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CommentOnPostRequest{Comment: "Nice"}
	repo.EXPECT().CommentOnPost(req).Return(nil)

	if err := svc.CommentOnPost(req); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostService_CommentOnPost_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CommentOnPostRequest{Comment: "Bad"}
	repo.EXPECT().CommentOnPost(req).Return(errors.New("fail"))

	if err := svc.CommentOnPost(req); err == nil || err.Error() != "fail" {
		t.Fatalf("expected fail, got %v", err)
	}
}

func TestPostService_CommentOnComment_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CommentOnCommentRequest{Comment: "Reply"}
	repo.EXPECT().CommentOnComment(req).Return(nil)

	if err := svc.CommentOnComment(req); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestPostService_CommentOnComment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockPostRepository(ctrl)
	svc := NewPostService(repo)

	req := model.CommentOnCommentRequest{Comment: "Reply"}
	repo.EXPECT().CommentOnComment(req).Return(errors.New("bad"))

	if err := svc.CommentOnComment(req); err == nil || err.Error() != "bad" {
		t.Fatalf("expected bad, got %v", err)
	}
}
