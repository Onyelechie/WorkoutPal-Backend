package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestRoutineService_CreateRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	const userID int64 = 1
	req := model.CreateRoutineRequest{Name: "Push Day"}
	want := &model.ExerciseRoutine{ID: 10, UserID: userID, Name: req.Name}

	repo.EXPECT().
		CreateRoutine(userID, gomock.AssignableToTypeOf(model.CreateRoutineRequest{})).
		Return(want, nil)

	got, err := svc.CreateRoutine(userID, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != want.ID || got.UserID != want.UserID || got.Name != want.Name {
		t.Fatalf("unexpected routine: %#v", got)
	}
}

func TestRoutineService_CreateRoutine_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	const userID int64 = 2
	req := model.CreateRoutineRequest{Name: "Leg Day"}

	repo.EXPECT().
		CreateRoutine(userID, gomock.AssignableToTypeOf(model.CreateRoutineRequest{})).
		Return((*model.ExerciseRoutine)(nil), errors.New("validation failed"))

	got, err := svc.CreateRoutine(userID, req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "validation failed" {
		t.Fatalf("expected validation failed, got %v", err)
	}
}

func TestRoutineService_ReadUserRoutines_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	const userID int64 = 3
	want := []*model.ExerciseRoutine{
		{ID: 1, UserID: userID, Name: "A"},
		{ID: 2, UserID: userID, Name: "B"},
	}

	repo.EXPECT().ReadUserRoutines(userID).Return(want, nil)

	got, err := svc.ReadUserRoutines(userID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != len(want) || got[0].ID != want[0].ID || got[1].Name != want[1].Name {
		t.Fatalf("unexpected routines: %#v", got)
	}
}

func TestRoutineService_ReadUserRoutines_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	const userID int64 = 4
	repo.EXPECT().ReadUserRoutines(userID).Return(nil, errors.New("user not found"))

	got, err := svc.ReadUserRoutines(userID)
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestRoutineService_DeleteRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().DeleteRoutine(int64(5)).Return(nil)

	if err := svc.DeleteRoutine(5); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestRoutineService_DeleteRoutine_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().DeleteRoutine(int64(6)).Return(errors.New("not found"))

	if err := svc.DeleteRoutine(6); err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRoutineService_ReadRoutineWithExercises_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	want := &model.ExerciseRoutine{ID: 7, Name: "Pull Day"}
	repo.EXPECT().ReadRoutineWithExercises(int64(7)).Return(want, nil)

	got, err := svc.ReadRoutineWithExercises(7)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != want.ID || got.Name != want.Name {
		t.Fatalf("unexpected routine: %#v", got)
	}
}

func TestRoutineService_ReadRoutineWithExercises_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().ReadRoutineWithExercises(int64(8)).
		Return((*model.ExerciseRoutine)(nil), errors.New("not found"))

	got, err := svc.ReadRoutineWithExercises(8)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRoutineService_AddExerciseToRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().AddExerciseToRoutine(int64(9), int64(100)).Return(nil)

	if err := svc.AddExerciseToRoutine(9, 100); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestRoutineService_AddExerciseToRoutine_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().AddExerciseToRoutine(int64(9), int64(100)).
		Return(errors.New("already added"))

	if err := svc.AddExerciseToRoutine(9, 100); err == nil || err.Error() != "already added" {
		t.Fatalf("expected already added, got %v", err)
	}
}

func TestRoutineService_RemoveExerciseFromRoutine_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().RemoveExerciseFromRoutine(int64(10), int64(101)).Return(nil)

	if err := svc.RemoveExerciseFromRoutine(10, 101); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestRoutineService_RemoveExerciseFromRoutine_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockRoutineRepository(ctrl)
	svc := NewRoutineService(repo)

	repo.EXPECT().RemoveExerciseFromRoutine(int64(10), int64(101)).
		Return(errors.New("not in routine"))

	if err := svc.RemoveExerciseFromRoutine(10, 101); err == nil || err.Error() != "not in routine" {
		t.Fatalf("expected not in routine, got %v", err)
	}
}
