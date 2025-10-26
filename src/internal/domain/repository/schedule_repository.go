package repository

import "workoutpal/src/internal/model"

type ScheduleRepository interface {
	ReadUserSchedules(userId int64) ([]*model.Schedule, error)
	ReadUserSchedulesByDay(userId int64, dayOfWeek int64) ([]*model.Schedule, error)
	ReadScheduleByID(id int64) (*model.Schedule, error)
	CreateSchedule(request model.CreateScheduleRequest) (*model.Schedule, error)
	UpdateSchedule(request model.UpdateScheduleRequest) (*model.Schedule, error)
	DeleteSchedule(request model.DeleteScheduleRequest) error
}
