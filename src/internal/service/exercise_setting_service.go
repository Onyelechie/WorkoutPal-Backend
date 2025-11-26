package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type exerciseSettingService struct {
	exerciseSettingRepository repository.ExerciseSettingRepository
}

func NewExerciseSettingService(exerciseSettingRepository repository.ExerciseSettingRepository) service.ExerciseSettingService {
	return &exerciseSettingService{exerciseSettingRepository: exerciseSettingRepository}
}

func (s *exerciseSettingService) ReadExerciseSetting(req model.ReadExerciseSettingRequest) (*model.ExerciseSetting, error) {
	return s.exerciseSettingRepository.ReadExerciseSetting(req)
}

func (s *exerciseSettingService) CreateExerciseSetting(req model.CreateExerciseSettingRequest) (*model.ExerciseSetting, error) {
	return s.exerciseSettingRepository.CreateExerciseSetting(req)
}

func (s *exerciseSettingService) UpdateExerciseSetting(req model.UpdateExerciseSettingRequest) (*model.ExerciseSetting, error) {
	return s.exerciseSettingRepository.UpdateExerciseSetting(req)
}
