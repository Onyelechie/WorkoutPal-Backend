package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type AchievementService struct {
	repo repository.AchievementRepository
}

func NewAchievementService(repo repository.AchievementRepository) service.AchievementService {
	return &AchievementService{repo: repo}
}

func (s *AchievementService) ReadAchievementsFeed() ([]*model.UserAchievement, error) {
	return s.repo.ReadAchievementsFeed()
}

func (s *AchievementService) ReadAllAchievements() ([]*model.Achievement, error) {
	return s.repo.ReadAllAchievements()
}

func (s *AchievementService) ReadUnlockedAchievements(userID int64) ([]*model.UserAchievement, error) {
	return s.repo.ReadUnlockedAchievements(userID)
}

func (s *AchievementService) CreateAchievement(req model.CreateAchievementRequest) (*model.UserAchievement, error) {
	return s.repo.CreateAchievement(req)
}
