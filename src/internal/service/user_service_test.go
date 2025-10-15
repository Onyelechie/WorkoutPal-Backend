package service

import (
	"errors"
	"testing"
	"workoutpal/src/internal/model"

	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestUserService_ReadUserByEmail_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	want := &model.User{ID: 1, Email: "a@b.com"}
	repo.EXPECT().ReadUserByEmail("a@b.com").Return(want, nil)

	got, err := svc.ReadUserByEmail("a@b.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Email != want.Email {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserService_ReadUserByEmail_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	repo.EXPECT().ReadUserByEmail("x@y.com").Return((*model.User)(nil), errors.New("not found"))

	got, err := svc.ReadUserByEmail("x@y.com")
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUserService_ReadUsers_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	want := []*model.User{
		{ID: 1, Username: "max"},
		{ID: 2, Username: "sam"},
	}
	repo.EXPECT().ReadUsers().Return(want, nil)

	got, err := svc.ReadUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Username != want[1].Username {
		t.Fatalf("unexpected users: %#v", got)
	}
}

func TestUserService_ReadUsers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	repo.EXPECT().ReadUsers().Return(nil, errors.New("db down"))

	got, err := svc.ReadUsers()
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down, got %v", err)
	}
}

func TestUserService_ReadUserByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	want := &model.User{ID: 42, Username: "max"}
	repo.EXPECT().ReadUserByID(int64(42)).Return(want, nil)

	got, err := svc.ReadUserByID(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserService_ReadUserByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	repo.EXPECT().ReadUserByID(int64(7)).Return((*model.User)(nil), errors.New("not found"))

	got, err := svc.ReadUserByID(7)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestUserService_CreateUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.CreateUserRequest{Username: "max", Email: "a@b.com", Name: "Max", Password: "Str0ng!Pass"}
	want := &model.User{ID: 1, Username: req.Username, Email: req.Email, Name: req.Name}

	repo.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(model.CreateUserRequest{})).
		Return(want, nil)

	got, err := svc.CreateUser(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserService_CreateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.CreateUserRequest{Username: "max"}
	repo.EXPECT().
		CreateUser(gomock.AssignableToTypeOf(model.CreateUserRequest{})).
		Return((*model.User)(nil), errors.New("username taken"))

	got, err := svc.CreateUser(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "username taken" {
		t.Fatalf("expected username taken, got %v", err)
	}
}

func TestUserService_UpdateUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.UpdateUserRequest{ID: 9, Username: "newname"}
	want := &model.User{ID: 9, Username: req.Username}

	repo.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(model.UpdateUserRequest{})).
		Return(want, nil)

	got, err := svc.UpdateUser(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserService_UpdateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.UpdateUserRequest{ID: 9}
	repo.EXPECT().
		UpdateUser(gomock.AssignableToTypeOf(model.UpdateUserRequest{})).
		Return((*model.User)(nil), errors.New("bad update"))

	got, err := svc.UpdateUser(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "bad update" {
		t.Fatalf("expected bad update, got %v", err)
	}
}

func TestUserService_DeleteUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.DeleteUserRequest{ID: 13}
	repo.EXPECT().DeleteUser(req).Return(nil)

	if err := svc.DeleteUser(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewUserService(repo)

	req := model.DeleteUserRequest{ID: 13}
	repo.EXPECT().DeleteUser(req).Return(errors.New("not found"))

	if err := svc.DeleteUser(req); err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}
