package repository

import "workoutpal/src/internal/model"

type RoutineRepository interface {
	CreateRoutine(userID int64, request model.CreateRoutineRequest) (*model.ExerciseRoutine, error)
	ReadUserRoutines(userID int64) ([]*model.ExerciseRoutine, error)
	ReadRoutineWithExercises(routineID int64) (*model.ExerciseRoutine, error)
	AddExerciseToRoutine(routineID, exerciseID int64) error
	RemoveExerciseFromRoutine(routineID, exerciseID int64) error
	DeleteRoutine(routineID int64) error
}
