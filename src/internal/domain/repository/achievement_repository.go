package repository

import "workoutpal/src/internal/model"

type AchievementRepository interface {
	ReadAllAchievements() ([]*model.Achievement, error)

	ReadUnlockedAchievementByAchievementID(id int64) (*model.UserAchievement, error)
	ReadUnlockedAchievements(userID int64) ([]*model.UserAchievement, error)
	CreateAchievement(a model.CreateAchievementRequest) (*model.UserAchievement, error)
}
