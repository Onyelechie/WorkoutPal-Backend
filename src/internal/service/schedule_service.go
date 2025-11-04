package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type scheduleService struct {
	repository repository.ScheduleRepository
}

func NewScheduleService(repository repository.ScheduleRepository) service.ScheduleService {
	return &scheduleService{repository: repository}
}

func (s *scheduleService) ReadUserSchedules(userId int64) ([]*model.Schedule, error) {
	return s.repository.ReadUserSchedules(userId)
}

func (s *scheduleService) ReadUserSchedulesByDay(userId int64, dayOfWeek int64) ([]*model.Schedule, error) {
	return s.repository.ReadUserSchedulesByDay(userId, dayOfWeek)
}

func (s *scheduleService) ReadScheduleByID(id int64) (*model.Schedule, error) {
	return s.repository.ReadScheduleByID(id)
}

func (s *scheduleService) CreateSchedule(request model.CreateScheduleRequest) (*model.Schedule, error) {
	schedule, err := s.repository.CreateSchedule(request)
	if err != nil {
		return nil, err
	}
	return schedule, nil
}

func (s *scheduleService) UpdateSchedule(request model.UpdateScheduleRequest) (*model.Schedule, error) {
	return s.repository.UpdateSchedule(request)
}

func (s *scheduleService) DeleteSchedule(request model.DeleteScheduleRequest) error {
	return s.repository.DeleteSchedule(request)
}
