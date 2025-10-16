package service

import (
	"errors"
	"testing"
	"workoutpal/src/internal/model"

	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestGoalService_CreateGoal_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockGoalRepository(ctrl)
	svc := NewGoalService(repo)

	const userID int64 = 1
	req := model.CreateGoalRequest{
		Name:        "Run 5K",
		Description: "Finish under 30min",
		Deadline:    "2025-12-31",
	}
	want := &model.Goal{
		ID:          10,
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Deadline:    req.Deadline,
		CreatedAt:   "2025-10-15T00:00:00Z",
	}

	repo.EXPECT().
		CreateGoal(userID, gomock.AssignableToTypeOf(model.CreateGoalRequest{})).
		Return(want, nil)

	got, err := svc.CreateGoal(userID, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.UserID != want.UserID ||
		got.Name != want.Name || got.Description != want.Description ||
		got.Deadline != want.Deadline || got.Status != want.Status {
		t.Fatalf("unexpected goal: %#v", got)
	}
}

func TestGoalService_CreateGoal_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockGoalRepository(ctrl)
	svc := NewGoalService(repo)

	const userID int64 = 2
	req := model.CreateGoalRequest{Name: "Bench 225"}
	repo.EXPECT().
		CreateGoal(userID, gomock.AssignableToTypeOf(model.CreateGoalRequest{})).
		Return((*model.Goal)(nil), errors.New("validation failed"))

	got, err := svc.CreateGoal(userID, req)
	if got != nil {
		t.Fatalf("expected nil goal, got %#v", got)
	}
	if err == nil || err.Error() != "validation failed" {
		t.Fatalf("expected validation failed, got %v", err)
	}
}

func TestGoalService_ReadUserGoals_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockGoalRepository(ctrl)
	svc := NewGoalService(repo)

	const userID int64 = 3
	want := []*model.Goal{
		{ID: 1, UserID: userID, Name: "Lose 5lb", Status: "active"},
		{ID: 2, UserID: userID, Name: "Deadlift 180kg", Status: "paused"},
	}
	repo.EXPECT().ReadUserGoals(userID).Return(want, nil)

	got, err := svc.ReadUserGoals(userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Name != want[1].Name {
		t.Fatalf("unexpected goals: %#v", got)
	}
}
