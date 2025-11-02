package repository

import "workoutpal/src/internal/model"

type AchievementRepository interface {
	ReadAchievements() ([]*model.Achievement, error)
	CreateAchievement(a model.CreateAchievementRequest) (*model.Achievement, error)
	DeleteAchievement(id int64) error
}
