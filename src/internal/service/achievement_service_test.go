package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestAchievementService_Read_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	want := []*model.Achievement{{ID: 1}, {ID: 2}}
	repo.EXPECT().ReadAchievements().Return(want, nil)

	got, err := svc.ReadAchievements()
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d", len(got))
	}
}

func TestAchievementService_Read_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	repo.EXPECT().ReadAchievements().Return(nil, errors.New("boom"))

	got, err := svc.ReadAchievements()
	if got != nil || err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementService_Create_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	req := model.CreateAchievementRequest{UserID: 1, Title: "A"}
	want := &model.Achievement{ID: 10, UserID: 1, Title: "A"}
	repo.EXPECT().CreateAchievement(gomock.AssignableToTypeOf(model.CreateAchievementRequest{})).
		Return(want, nil)

	got, err := svc.CreateAchievement(req)
	if err != nil || got.ID != want.ID {
		t.Fatalf("unexpected: %+v err=%v", got, err)
	}
}

func TestAchievementService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	req := model.CreateAchievementRequest{UserID: 1, Title: "A"}
	repo.EXPECT().CreateAchievement(gomock.AssignableToTypeOf(model.CreateAchievementRequest{})).
		Return(nil, errors.New("fail"))

	got, err := svc.CreateAchievement(req)
	if got != nil || err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementService_Delete_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	repo.EXPECT().DeleteAchievement(int64(5)).Return(nil)

	if err := svc.DeleteAchievement(5); err != nil {
		t.Fatalf("err=%v", err)
	}
}

func TestAchievementService_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockAchievementRepository(ctrl)
	svc := NewAchievementService(repo)

	repo.EXPECT().DeleteAchievement(int64(6)).Return(errors.New("not found"))

	if err := svc.DeleteAchievement(6); err == nil {
		t.Fatalf("expected error")
	}
}
