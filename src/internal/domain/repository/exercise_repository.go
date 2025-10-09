package repository

import "workoutpal/src/internal/model"

type ExerciseRepository interface {
	GetAllExercises() ([]model.Exercise, error)
}