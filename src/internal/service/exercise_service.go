package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type exerciseService struct {
	exerciseRepository repository.ExerciseRepository
}

func NewExerciseService(er repository.ExerciseRepository) service.ExerciseService {
	return &exerciseService{
		exerciseRepository: er,
	}
}

func (e *exerciseService) ReadAllExercises() ([]*model.Exercise, error) {
	return e.exerciseRepository.ReadAllExercises()
}

func (e *exerciseService) ReadExerciseByID(id int64) (*model.Exercise, error) {
	return e.exerciseRepository.ReadExerciseByID(id)
}
