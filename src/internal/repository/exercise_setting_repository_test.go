package repository

import (
	"errors"
	"regexp"
	"testing"

	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExerciseSettingRepository_ReadExerciseSetting_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.ReadExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	rows := sqlmock.NewRows([]string{
		"user_id", "exercise_id", "workout_routine_id", "weight", "reps", "sets", "break_interval",
	}).AddRow(
		int64(1), int64(2), int64(3), 50, 8, 3, 90,
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT user_id, exercise_id, workout_routine_id, weight, reps, sets, break_interval
		FROM user_exercise_settings
		WHERE user_id = $1 AND exercise_id = $2 AND workout_routine_id = $3`)).
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID).
		WillReturnRows(rows)

	got, err := repo.ReadExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil result")
	}
	if got.UserID != req.UserID || got.ExerciseID != req.ExerciseID || got.WorkoutRoutineID != req.WorkoutRoutineID {
		t.Fatalf("unexpected keys: %#v", got)
	}
	if got.Weight != 50 || got.Reps != 8 || got.Sets != 3 || got.BreakInterval != 90 {
		t.Fatalf("unexpected values: %#v", got)
	}
}

func TestExerciseSettingRepository_ReadExerciseSetting_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.ReadExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	mock.ExpectQuery("SELECT user_id").
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID).
		WillReturnError(errors.New("query fail"))

	got, err := repo.ReadExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "query fail" {
		t.Fatalf("expected query fail, got %v", err)
	}
}

func TestExerciseSettingRepository_CreateExerciseSetting_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.CreateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           60,
		Reps:             10,
		Sets:             4,
		BreakInterval:    120,
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO user_exercise_settings(user_id, exercise_id, workout_routine_id, weight, reps, sets, break_interval)
					VALUES ($1, $2, $3, $4, $5, $6, $7)`)).
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval).
		WillReturnResult(sqlmock.NewResult(0, 1))

	readRows := sqlmock.NewRows([]string{
		"user_id", "exercise_id", "workout_routine_id", "weight", "reps", "sets", "break_interval",
	}).AddRow(
		int64(1), int64(2), int64(3), 60, 10, 4, 120,
	)

	mock.ExpectQuery("SELECT user_id").
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID).
		WillReturnRows(readRows)

	got, err := repo.CreateExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil result")
	}
	if got.Weight != 60 || got.Reps != 10 || got.Sets != 4 || got.BreakInterval != 120 {
		t.Fatalf("unexpected values: %#v", got)
	}
}

func TestExerciseSettingRepository_CreateExerciseSetting_ErrorOnInsert(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.CreateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	mock.ExpectExec("INSERT INTO user_exercise_settings").
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval).
		WillReturnError(errors.New("insert fail"))

	got, err := repo.CreateExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "insert fail" {
		t.Fatalf("expected insert fail, got %v", err)
	}
}

func TestExerciseSettingRepository_UpdateExerciseSetting_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.UpdateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
		Weight:           70,
		Reps:             12,
		Sets:             5,
		BreakInterval:    90,
	}

	// UPDATE
	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE user_exercise_settings SET weight = $4, reps = $5, 
                                  sets = $6, break_interval = $7
					WHERE user_id = $1 AND exercise_id = $2 AND workout_routine_id = $3`)).
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// ReadExerciseSetting called after update
	readRows := sqlmock.NewRows([]string{
		"user_id", "exercise_id", "workout_routine_id", "weight", "reps", "sets", "break_interval",
	}).AddRow(
		int64(1), int64(2), int64(3), 70, 12, 5, 90,
	)

	mock.ExpectQuery("SELECT user_id").
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID).
		WillReturnRows(readRows)

	got, err := repo.UpdateExerciseSetting(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil {
		t.Fatalf("expected non-nil result")
	}
	if got.Weight != 70 || got.Reps != 12 || got.Sets != 5 || got.BreakInterval != 90 {
		t.Fatalf("unexpected values: %#v", got)
	}
}

func TestExerciseSettingRepository_UpdateExerciseSetting_ErrorOnUpdate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewExerciseSettingRepository(db)

	req := model.UpdateExerciseSettingRequest{
		UserID:           1,
		ExerciseID:       2,
		WorkoutRoutineID: 3,
	}

	mock.ExpectExec("UPDATE user_exercise_settings").
		WithArgs(req.UserID, req.ExerciseID, req.WorkoutRoutineID, req.Weight, req.Reps, req.Sets, req.BreakInterval).
		WillReturnError(errors.New("update fail"))

	got, err := repo.UpdateExerciseSetting(req)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "update fail" {
		t.Fatalf("expected update fail, got %v", err)
	}
}
