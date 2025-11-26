package repository

import "workoutpal/src/internal/model"

type ExerciseSettingRepository interface {
	ReadExerciseSetting(req model.ReadExerciseSettingRequest) (*model.ExerciseSetting, error)
	CreateExerciseSetting(req model.CreateExerciseSettingRequest) (*model.ExerciseSetting, error)
	UpdateExerciseSetting(req model.UpdateExerciseSettingRequest) (*model.ExerciseSetting, error)
}
