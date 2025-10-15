package service

import "workoutpal/src/internal/model"

type GoalService interface {
	CreateGoal(userID int64, request model.CreateGoalRequest) (*model.Goal, error)
	ReadUserGoals(userID int64) ([]*model.Goal, error)
}
