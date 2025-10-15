package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type goalService struct {
	goalRepository repository.GoalRepository
}

func NewGoalService(goalRepository repository.GoalRepository) service.GoalService {
	return &goalService{goalRepository: goalRepository}
}

func (u *goalService) CreateGoal(userID int64, request model.CreateGoalRequest) (*model.Goal, error) {
	return u.goalRepository.CreateGoal(userID, request)
}

func (u *goalService) ReadUserGoals(userID int64) ([]*model.Goal, error) {
	return u.goalRepository.ReadUserGoals(userID)
}
