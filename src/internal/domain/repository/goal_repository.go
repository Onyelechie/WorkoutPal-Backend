package repository

import "workoutpal/src/internal/model"

type GoalRepository interface {
	CreateGoal(userID int64, request model.CreateGoalRequest) (*model.Goal, error)
	ReadUserGoals(userID int64) ([]*model.Goal, error)
	UpdateGoal(request model.UpdateGoalRequest) (*model.Goal, error)
	DeleteGoal(goalID int64) error
}
