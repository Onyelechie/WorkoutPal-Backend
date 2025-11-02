package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type AchievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) *AchievementService {
	return &AchievementService{repo: repo}
}

func (s *AchievementService) ReadAchievements(userID int64) ([]*model.Achievement, error) {
	return s.repo.ReadAchievements(userID)
}

func (s *AchievementService) CreateAchievement(req model.CreateAchievementRequest) (*model.Achievement, error) {
	return s.repo.CreateAchievement(req)
}

func (s *AchievementService) DeleteAchievement(id int64) error {
	return s.repo.DeleteAchievement(id)
}
