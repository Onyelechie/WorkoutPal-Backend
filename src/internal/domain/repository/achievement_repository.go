package repository

import "workoutpal/src/internal/model"

type AchievementRepository interface {
	ReadAchievements(userID int64) ([]*model.Achievement, error)
	CreateAchievement(a model.CreateAchievementRequest) (*model.Achievement, error)
	DeleteAchievement(id int64) error
}
