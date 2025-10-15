package repository

import (
	"workoutpal/src/internal/model"
)

type ExerciseRepository interface {
	ReadExerciseByID(id int64) (*model.Exercise, error)
	ReadAllExercises() ([]*model.Exercise, error)
}
