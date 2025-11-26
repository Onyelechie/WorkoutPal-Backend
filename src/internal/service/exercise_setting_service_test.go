package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestExerciseSettingService_ReadExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.ReadExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}
	want := &model.ExerciseSetting{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           50,
		Reps:             8,
		Sets:             3,
		BreakInterval:    90,
	}

	repo.EXPECT().
		ReadExerciseSetting(req).
		Return(want, nil)

	got, err := svc.ReadExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.UserID != want.UserID || got.ExerciseID != want.ExerciseID || got.WorkoutRoutineID != want.WorkoutRoutineID {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestExerciseSettingService_ReadExerciseSetting_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.ReadExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	repo.EXPECT().
		ReadExerciseSetting(req).
		Return((*model.ExerciseSetting)(nil), errors.New("read fail"))

	got, err := svc.ReadExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "read fail" {
		t.Fatalf("expected read fail, got %v", err)
	}
}

func TestExerciseSettingService_CreateExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.CreateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           60,
		Reps:             10,
		Sets:             4,
		BreakInterval:    120,
	}
	want := &model.ExerciseSetting{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           60,
		Reps:             10,
		Sets:             4,
		BreakInterval:    120,
	}

	repo.EXPECT().
		CreateExerciseSetting(req).
		Return(want, nil)

	got, err := svc.CreateExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.Weight != want.Weight || got.Reps != want.Reps || got.Sets != want.Sets {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestExerciseSettingService_CreateExerciseSetting_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.CreateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	repo.EXPECT().
		CreateExerciseSetting(req).
		Return((*model.ExerciseSetting)(nil), errors.New("create fail"))

	got, err := svc.CreateExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "create fail" {
		t.Fatalf("expected create fail, got %v", err)
	}
}

func TestExerciseSettingService_UpdateExerciseSetting_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.UpdateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}
	want := &model.ExerciseSetting{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}

	repo.EXPECT().
		UpdateExerciseSetting(req).
		Return(want, nil)

	got, err := svc.UpdateExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.Weight != want.Weight || got.Reps != want.Reps || got.Sets != want.Sets {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestExerciseSettingService_UpdateExerciseSetting_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mock_repository.NewMockExerciseSettingRepository(ctrl)
	svc := NewExerciseSettingService(repo)

	req := model.UpdateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	repo.EXPECT().
		UpdateExerciseSetting(req).
		Return((*model.ExerciseSetting)(nil), errors.New("update fail"))

	got, err := svc.UpdateExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "update fail" {
		t.Fatalf("expected update fail, got %v", err)
	}
}
