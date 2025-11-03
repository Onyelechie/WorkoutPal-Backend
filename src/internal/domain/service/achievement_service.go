package service

import "workoutpal/src/internal/model"

type AchievementService interface {
	ReadAllAchievements() ([]*model.Achievement, error)

	ReadUnlockedAchievements(userID int64) ([]*model.UserAchievement, error)
	CreateAchievement(req model.CreateAchievementRequest) (*model.UserAchievement, error)
}
