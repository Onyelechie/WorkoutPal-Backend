package service

import (
	"workoutpal/src/internal/model"
)

type ExerciseService interface {
	ReadExerciseByID(id int64) (*model.Exercise, error)
	ReadAllExercises() ([]*model.Exercise, error)
}
