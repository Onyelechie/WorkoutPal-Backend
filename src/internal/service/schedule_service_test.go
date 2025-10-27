package service

import (
	"errors"
	"testing"

	"workoutpal/src/internal/model"
	mock_repository "workoutpal/src/mock_internal/domain/repository"

	"github.com/golang/mock/gomock"
)

func TestScheduleService_ReadUserSchedules_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	userID := int64(55)
	want := []*model.Schedule{
		{ID: 1, UserID: userID, Name: "A"},
		{ID: 2, UserID: userID, Name: "B"},
	}

	mockRepo.EXPECT().
		ReadUserSchedules(userID).
		Return(want, nil)

	got, err := svc.ReadUserSchedules(userID)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 || got[0].Name != "A" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestScheduleService_ReadUserSchedules_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	userID := int64(55)

	mockRepo.EXPECT().
		ReadUserSchedules(userID).
		Return(nil, errors.New("boom"))

	got, err := svc.ReadUserSchedules(userID)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "boom" {
		t.Fatalf("expected boom, got %v", err)
	}
}

func TestScheduleService_ReadUserSchedulesByDay_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	userID := int64(10)
	day := int64(3)
	want := []*model.Schedule{
		{ID: 9, UserID: userID, DayOfWeek: day, Name: "Legs"},
	}

	mockRepo.EXPECT().
		ReadUserSchedulesByDay(userID, day).
		Return(want, nil)

	got, err := svc.ReadUserSchedulesByDay(userID, day)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Legs" || got[0].DayOfWeek != day {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestScheduleService_ReadScheduleByID_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	id := int64(123)
	want := &model.Schedule{ID: id, Name: "Pull"}

	mockRepo.EXPECT().
		ReadScheduleByID(id).
		Return(want, nil)

	got, err := svc.ReadScheduleByID(id)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != id {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestScheduleService_CreateSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.CreateScheduleRequest{
		Name:                 "AM Upper",
		UserID:               77,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}

	want := &model.Schedule{
		ID:                   900,
		Name:                 "AM Upper",
		UserID:               77,
		DayOfWeek:            1,
		TimeSlot:             "07:15",
		RoutineLengthMinutes: 80,
		RoutineIDs:           []int64{10, 11},
	}

	mockRepo.EXPECT().
		CreateSchedule(req).
		Return(want, nil)

	got, err := svc.CreateSchedule(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 900 {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestScheduleService_CreateSchedule_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.CreateScheduleRequest{
		Name:      "oops",
		UserID:    1,
		DayOfWeek: 2,
		TimeSlot:  "bad",
	}

	mockRepo.EXPECT().
		CreateSchedule(req).
		Return(nil, errors.New("create fail"))

	got, err := svc.CreateSchedule(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "create fail" {
		t.Fatalf("expected create fail, got %v", err)
	}
}

func TestScheduleService_UpdateSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.UpdateScheduleRequest{
		ID:                   321,
		UserID:               44,
		Name:                 "PM Lower",
		DayOfWeek:            5,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}
	want := &model.Schedule{
		ID:                   321,
		UserID:               44,
		Name:                 "PM Lower",
		DayOfWeek:            5,
		TimeSlot:             "20:00",
		RoutineLengthMinutes: 55,
		RoutineIDs:           []int64{7, 8},
	}

	mockRepo.EXPECT().
		UpdateSchedule(req).
		Return(want, nil)

	got, err := svc.UpdateSchedule(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.Name != "PM Lower" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestScheduleService_UpdateSchedule_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.UpdateScheduleRequest{
		ID:        999,
		UserID:    44,
		Name:      "doesn't matter",
		DayOfWeek: 6,
		TimeSlot:  "xx",
	}

	mockRepo.EXPECT().
		UpdateSchedule(req).
		Return(nil, errors.New("update fail"))

	got, err := svc.UpdateSchedule(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "update fail" {
		t.Fatalf("expected update fail, got %v", err)
	}
}

func TestScheduleService_DeleteSchedule_OK(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.DeleteScheduleRequest{ID: 123}

	mockRepo.EXPECT().
		DeleteSchedule(req).
		Return(nil)

	if err := svc.DeleteSchedule(req); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

func TestScheduleService_DeleteSchedule_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockRepo := mock_repository.NewMockScheduleRepository(ctrl)
	svc := &scheduleService{repository: mockRepo}

	req := model.DeleteScheduleRequest{ID: 123}

	mockRepo.EXPECT().
		DeleteSchedule(req).
		Return(errors.New("delete fail"))

	err := svc.DeleteSchedule(req)
	if err == nil || err.Error() != "delete fail" {
		t.Fatalf("expected delete fail, got %v", err)
	}
}
