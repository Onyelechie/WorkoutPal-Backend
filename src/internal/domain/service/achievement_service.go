package service

import "workoutpal/src/internal/model"

type AchievementService interface {
	ReadAchievements(userID int64) ([]*model.Achievement, error)
	CreateAchievement(req model.CreateAchievementRequest) (*model.Achievement, error)
	DeleteAchievement(id int64) error
}
