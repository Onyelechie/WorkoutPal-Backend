package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestAchievementService_ReadAllAchievements_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	want := []*model.Achievement{
		{ID: 1, Title: "First Workout"},
		{ID: 2, Title: "7-Day Streak"},
	}

	repo.EXPECT().
		ReadAllAchievements().
		Return(want, nil)

	got, err := svc.ReadAllAchievements()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestAchievementService_ReadAllAchievements_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	repo.EXPECT().
		ReadAllAchievements().
		Return(nil, errors.New("boom"))

	got, err := svc.ReadAllAchievements()
	if got != nil || err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementService_ReadUnlockedAchievements_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	want := []*model.UserAchievement{
		{ID: 10, UserID: 1, Title: "First Workout", EarnedAt: "2025-01-01T00:00:00Z"},
		{ID: 11, UserID: 1, Title: "7-Day Streak", EarnedAt: "2025-01-05T00:00:00Z"},
	}

	repo.EXPECT().
		ReadUnlockedAchievements(int64(1)).
		Return(want, nil)

	got, err := svc.ReadUnlockedAchievements(1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 || got[0].ID != 10 || got[0].UserID != 1 {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestAchievementService_ReadUnlockedAchievements_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	repo.EXPECT().
		ReadUnlockedAchievements(int64(1)).
		Return(nil, errors.New("boom"))

	got, err := svc.ReadUnlockedAchievements(1)
	if got != nil || err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementService_CreateAchievement_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	req := model.CreateAchievementRequest{UserID: 1, AchievementID: 55}
	want := &model.UserAchievement{ID: 99, UserID: 1, Title: "First Workout", EarnedAt: "2025-01-01T00:00:00Z"}

	repo.EXPECT().
		CreateAchievement(req).
		Return(want, nil)

	got, err := svc.CreateAchievement(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 99 || got.UserID != 1 {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestAchievementService_CreateAchievement_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	req := model.CreateAchievementRequest{UserID: 1, AchievementID: 55}

	repo.EXPECT().
		CreateAchievement(req).
		Return(nil, errors.New("fail"))

	got, err := svc.CreateAchievement(req)
	if got != nil || err == nil {
		t.Fatalf("expected error")
	}
}
