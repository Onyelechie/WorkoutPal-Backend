package service

import "workoutpal/src/internal/model"

type ExerciseService interface {
	GetAllExercises() ([]model.Exercise, error)
}