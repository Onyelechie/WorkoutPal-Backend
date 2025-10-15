package service

import (
	"errors"
	"testing"
	"workoutpal/src/internal/model"

	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestExerciseService_ReadAllExercises_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseRepository(ctrl)
	svc := NewExerciseService(repo)

	want := []*model.Exercise{
		{ID: 1, Name: "Squat"},
		{ID: 2, Name: "Bench Press"},
	}

	repo.EXPECT().ReadAllExercises().Return(want, nil)

	got, err := svc.ReadAllExercises()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Name != want[1].Name {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestExerciseService_ReadAllExercises_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseRepository(ctrl)
	svc := NewExerciseService(repo)

	repo.EXPECT().ReadAllExercises().Return(nil, errors.New("db down"))

	_, err := svc.ReadAllExercises()
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down error, got %v", err)
	}
}

func TestExerciseService_ReadExerciseByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseRepository(ctrl)
	svc := NewExerciseService(repo)

	want := &model.Exercise{ID: 42, Name: "Deadlift"}
	repo.EXPECT().ReadExerciseByID(int64(42)).Return(want, nil)

	got, err := svc.ReadExerciseByID(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Name != want.Name {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestExerciseService_ReadExerciseByID_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseRepository(ctrl)
	svc := NewExerciseService(repo)

	repo.EXPECT().ReadExerciseByID(int64(7)).Return((*model.Exercise)(nil), errors.New("not found"))

	got, err := svc.ReadExerciseByID(7)
	if got != nil {
		t.Fatalf("expected nil exercise, got %#v", got)
	}
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found error, got %v", err)
	}
}
