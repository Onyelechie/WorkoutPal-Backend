package repository

import (
	"database/sql"
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type exerciseSettingRepository struct {
	db *sql.DB
}

func NewExerciseSettingRepository(db *sql.DB) repository.ExerciseSettingRepository {
	return &exerciseSettingRepository{db: db}
}

func (e *exerciseSettingRepository) ReadExerciseSetting(req model.ReadExerciseSettingRequest) (*model.ExerciseSetting, error) {
	row := e.db.QueryRow(`
		SELECT user_id, exercise_id, workout_routine_id, weight, reps, sets, break_interval
		FROM user_exercise_settings
		WHERE user_id = $1 AND exercise_id = $2 AND workout_routine_id = $3`,
		req.UserID, req.ExerciseID, req.WorkoutRoutineID,
	)

	var result model.ExerciseSetting
	err := row.Scan(
		&result.UserID,
		&result.ExerciseID,
		&result.WorkoutRoutineID,
		&result.Weight,
		&result.Reps,
		&result.Sets,
		&result.BreakInterval,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (e *exerciseSettingRepository) CreateExerciseSetting(req model.CreateExerciseSettingRequest) (*model.ExerciseSetting, error) {
	_, err := e.db.Exec(`INSERT INTO user_exercise_settings(user_id, exercise_id, workout_routine_id, weight, reps, sets, break_interval)
					VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval)
	if err != nil {
		return nil, err
	}

	return e.ReadExerciseSetting(model.ReadExerciseSettingRequest{UserID: req.UserID, ExerciseID: req.ExerciseID, WorkoutRoutineID: req.WorkoutRoutineID})
}

func (e *exerciseSettingRepository) UpdateExerciseSetting(req model.UpdateExerciseSettingRequest) (*model.ExerciseSetting, error) {
	_, err := e.db.Exec(`UPDATE user_exercise_settings SET weight = $4, reps = $5, 
                                  sets = $6, break_interval = $7
					WHERE user_id = $1 AND exercise_id = $2 AND workout_routine_id = $3`,
		req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval)
	if err != nil {
		return nil, err
	}

	return e.ReadExerciseSetting(model.ReadExerciseSettingRequest{UserID: req.UserID, ExerciseID: req.ExerciseID, WorkoutRoutineID: req.WorkoutRoutineID})
}
