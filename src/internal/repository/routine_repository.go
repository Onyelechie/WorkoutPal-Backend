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
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		// if Commit already happened, Rollback will return sql.ErrTxDone (ignore)
		_ = tx.Rollback()
	}()

	var routine model.ExerciseRoutine
	// Keep this minimal & consistent with model fields we know exist.
	// We’ll store frequency implicitly (if you add it to your model later, extend RETURNING).
	err = tx.QueryRow(`
		INSERT INTO workout_routine (name, user_id, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, user_id, description`,
		request.Name, userID, request.Description,
	).Scan(&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
	if err != nil {
		return nil, err
	}

	// Insert exercises into exercises_in_routine
	for _, exerciseID := range request.ExerciseIDs {
		if _, err := tx.Exec(
			"INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)",
			routine.ID, exerciseID,
		); err != nil {
			return nil, err
		}
	}
	routine.ExerciseIDs = append(routine.ExerciseIDs, request.ExerciseIDs...)
	routine.IsActive = true

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &routine, nil
}

func (r *routineRepository) ReadUserRoutines(userID int64) ([]*model.ExerciseRoutine, error) {
	rows, err := r.db.Query(
		"SELECT id, name, user_id, description FROM workout_routine WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routines []*model.ExerciseRoutine
	for rows.Next() {
		var routine model.ExerciseRoutine
		if err := rows.Scan(&routine.ID, &routine.Name, &routine.UserID, &routine.Description); err != nil {
			return nil, err
		}

		exRows, err := r.db.Query(
			"SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1",
			routine.ID,
		)
		if err != nil {
			return nil, err
		}

		var eids []int64
		for exRows.Next() {
			var eid int64
			if err := exRows.Scan(&eid); err != nil {
				exRows.Close()
				return nil, err
			}
			eids = append(eids, eid)
		}
		if err := exRows.Err(); err != nil {
			exRows.Close()
			return nil, err
		}
		exRows.Close()

		routine.ExerciseIDs = eids
		routine.IsActive = true
		routines = append(routines, &routine)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return routines, nil
}

func (r *routineRepository) DeleteRoutine(routineID int64) error {
	result, err := r.db.Exec("DELETE FROM workout_routine WHERE id = $1", routineID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("routine not found")
	}
	return nil
}

func (r *routineRepository) ReadRoutineWithExercises(routineID int64) (*model.ExerciseRoutine, error) {
	var routine model.ExerciseRoutine
	err := r.db.QueryRow(
		"SELECT id, name, user_id, description FROM workout_routine WHERE id = $1",
		routineID,
	).Scan(&routine.ID, &routine.Name, &routine.UserID, &routine.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("routine not found")
		}
		return nil, err
	}

	exRows, err := r.db.Query(
		"SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1",
		routine.ID,
	)
	if err != nil {
		return nil, err
	}
	var eids []int64
	for exRows.Next() {
		var eid int64
		if err := exRows.Scan(&eid); err != nil {
			exRows.Close()
			return nil, err
		}
		eids = append(eids, eid)
	}
	if err := exRows.Err(); err != nil {
		exRows.Close()
		return nil, err
	}
	exRows.Close()

	routine.ExerciseIDs = eids
	routine.IsActive = true
	return &routine, nil
}

func (r *routineRepository) AddExerciseToRoutine(routineID, exerciseID int64) error {
	_, err := r.db.Exec(
		"INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)",
		routineID, exerciseID,
	)
	return err
}

func (r *routineRepository) RemoveExerciseFromRoutine(routineID, exerciseID int64) error {
	res, err := r.db.Exec(
		"DELETE FROM exercises_in_routine WHERE workout_routine_id = $1 AND exercise_id = $2",
		routineID, exerciseID,
	)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("exercise not found in routine")
	}
	return nil
}
