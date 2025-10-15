package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type routineService struct {
	routineRepository repository.RoutineRepository
}

func NewRoutineService(routineRepository repository.RoutineRepository) service.RoutineService {
	return &routineService{routineRepository: routineRepository}
}

func (u *routineService) CreateRoutine(userID int64, request model.CreateRoutineRequest) (*model.ExerciseRoutine, error) {
	return u.routineRepository.CreateRoutine(userID, request)
}

func (u *routineService) ReadUserRoutines(userID int64) ([]*model.ExerciseRoutine, error) {
	return u.routineRepository.ReadUserRoutines(userID)
}

func (u *routineService) DeleteRoutine(routineID int64) error {
	return u.routineRepository.DeleteRoutine(routineID)
}

func (u *routineService) ReadRoutineWithExercises(routineID int64) (*model.ExerciseRoutine, error) {
	return u.routineRepository.ReadRoutineWithExercises(routineID)
}

func (u *routineService) AddExerciseToRoutine(routineID, exerciseID int64) error {
	return u.routineRepository.AddExerciseToRoutine(routineID, exerciseID)
}

func (u *routineService) RemoveExerciseFromRoutine(routineID, exerciseID int64) error {
	return u.routineRepository.RemoveExerciseFromRoutine(routineID, exerciseID)
}
