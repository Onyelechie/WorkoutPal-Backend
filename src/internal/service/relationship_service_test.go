package service

import (
	"errors"
	"testing"
	"workoutpal/src/internal/model"

	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestRelationshipService_FollowUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().FollowUser(int64(1), int64(2)).Return(nil)

	if err := svc.FollowUser(1, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelationshipService_FollowUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().FollowUser(int64(1), int64(2)).Return(errors.New("already following"))

	if err := svc.FollowUser(1, 2); err == nil || err.Error() != "already following" {
		t.Fatalf("expected already following, got %v", err)
	}
}

func TestRelationshipService_UnfollowUser_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().UnfollowUser(int64(3), int64(5)).Return(nil)

	if err := svc.UnfollowUser(3, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelationshipService_UnfollowUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().UnfollowUser(int64(3), int64(5)).Return(errors.New("not following"))

	if err := svc.UnfollowUser(3, 5); err == nil || err.Error() != "not following" {
		t.Fatalf("expected not following, got %v", err)
	}
}

func TestRelationshipService_ReadUserFollowers_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	followerIds := []int64{10, 11, 12}
	want := []model.User{
		{ID: 10, Name: "User1", Username: "user1"},
		{ID: 11, Name: "User2", Username: "user2"},
		{ID: 12, Name: "User3", Username: "user3"},
	}
	repo.EXPECT().ReadUserFollowers(int64(7)).Return(followerIds, nil)
	userRepo.EXPECT().ReadUserByID(int64(10)).Return(&want[0], nil)
	userRepo.EXPECT().ReadUserByID(int64(11)).Return(&want[1], nil)
	userRepo.EXPECT().ReadUserByID(int64(12)).Return(&want[2], nil)

	got, err := svc.ReadUserFollowers(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[2].ID != want[2].ID {
		t.Fatalf("unexpected followers: %#v", got)
	}
}

func TestRelationshipService_ReadUserFollowers_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().ReadUserFollowers(int64(7)).Return(nil, errors.New("user not found"))

	got, err := svc.ReadUserFollowers(7)
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestRelationshipService_ReadUserFollowing_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	followingIds := []int64{20, 21}
	want := []model.User{
		{ID: 20, Name: "User20", Username: "user20"},
		{ID: 21, Name: "User21", Username: "user21"},
	}
	repo.EXPECT().ReadUserFollowing(int64(8)).Return(followingIds, nil)
	userRepo.EXPECT().ReadUserByID(int64(20)).Return(&want[0], nil)
	userRepo.EXPECT().ReadUserByID(int64(21)).Return(&want[1], nil)

	got, err := svc.ReadUserFollowing(8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[1].ID != want[1].ID {
		t.Fatalf("unexpected following: %#v", got)
	}
}

func TestRelationshipService_ReadUserFollowing_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRelationshipRepository(ctrl)
	userRepo := mock_repository.NewMockUserRepository(ctrl)
	svc := NewRelationshipService(repo, userRepo)

	repo.EXPECT().ReadUserFollowing(int64(8)).Return(nil, errors.New("user not found"))

	got, err := svc.ReadUserFollowing(8)
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}
