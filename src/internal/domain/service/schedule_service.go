package service

import "workoutpal/src/internal/model"

type ScheduleService interface {
	ReadUserSchedules(userId int64) ([]*model.Schedule, error)
	ReadUserSchedulesByDay(userId int64, dayOfWeek int64) ([]*model.Schedule, error)
	ReadScheduleByID(id int64) (*model.Schedule, error)
	CreateSchedule(request model.CreateScheduleRequest) (*model.Schedule, error)
	UpdateSchedule(request model.UpdateScheduleRequest) (*model.Schedule, error)
	DeleteSchedule(request model.DeleteScheduleRequest) error
}
