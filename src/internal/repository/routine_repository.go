package repository

import (
	"database/sql"
	"errors"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type routineRepository struct {
	db *sql.DB
}

func NewRoutineRepository(db *sql.DB) repository.RoutineRepository {
	return &routineRepository{db: db}
}

func (r *routineRepository) CreateRoutine(userID int64, request model.CreateRoutineRequest) (*model.ExerciseRoutine, error) {
	var routine model.ExerciseRoutine
	err := r.db.QueryRow(`
		INSERT INTO workout_routine (name, user_id, frequency) 
		VALUES ($1, $2, 'weekly') 
		RETURNING id, name, user_id, frequency`,
		request.Name, userID).Scan(
		&routine.ID, &routine.Name, &routine.UserID, &routine.Description)

	if err != nil {
		return nil, err
	}

	// Insert exercises into exercises_in_routine table
	for _, exerciseID := range request.ExerciseIDs {
		_, err := r.db.Exec("INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)", routine.ID, exerciseID)
		if err != nil {
			return nil, err
		}
	}

	routine.Description = request.Description
	routine.IsActive = true
	return &routine, nil
}

func (r *routineRepository) ReadUserRoutines(userID int64) ([]*model.ExerciseRoutine, error) {
	rows, err := r.db.Query("SELECT id, name, user_id, frequency FROM workout_routine WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []*model.ExerciseRoutine
	for rows.Next() {
		var routine model.ExerciseRoutine
		err := rows.Scan(&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
		if err != nil {
			return nil, err
		}

		// Fetch exercises for the routine
		exerciseRows, err := r.db.Query("SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1", routine.ID)
		if err != nil {
			return nil, err
		}
		defer exerciseRows.Close()

		var exerciseIDs []int64
		for exerciseRows.Next() {
			var exerciseID int64
			if err := exerciseRows.Scan(&exerciseID); err != nil {
				return nil, err
			}
			exerciseIDs = append(exerciseIDs, exerciseID)
		}
		routine.ExerciseIDs = exerciseIDs

		routine.IsActive = true
		routines = append(routines, &routine)
	}
	return routines, nil
}

func (r *routineRepository) DeleteRoutine(routineID int64) error {
	result, err := r.db.Exec("DELETE FROM workout_routine WHERE id = $1", routineID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("routine not found")
	}
	return nil
}

func (r *routineRepository) ReadRoutineWithExercises(routineID int64) (*model.ExerciseRoutine, error) {
	var routine model.ExerciseRoutine
	err := r.db.QueryRow("SELECT id, name, user_id, frequency FROM workout_routine WHERE id = $1", routineID).Scan(
		&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("routine not found")
		}
		return nil, err
	}
	routine.IsActive = true
	return &routine, nil
}

func (r *routineRepository) AddExerciseToRoutine(routineID, exerciseID int64) error {
	_, err := r.db.Exec("INSERT INTO routine_exercises (routine_id, exercise_id) VALUES ($1, $2)", routineID, exerciseID)
	return err
}

func (r *routineRepository) RemoveExerciseFromRoutine(routineID, exerciseID int64) error {
	result, err := r.db.Exec("DELETE FROM routine_exercises WHERE routine_id = $1 AND exercise_id = $2", routineID, exerciseID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("exercise not found in routine")
	}
	return nil
}
