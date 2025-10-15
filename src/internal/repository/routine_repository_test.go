package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRoutineRepository_CreateRoutine_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	req := model.CreateRoutineRequest{
		Name:        "Push",
		Description: "Upper push day",
		ExerciseIDs: []int64{1, 2},
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO workout_routine (name, user_id, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, user_id, description`)).
		WithArgs(req.Name, int64(10), req.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "description"}).
			AddRow(100, req.Name, 10, req.Description))

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)",
	)).WithArgs(int64(100), int64(1)).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)",
	)).WithArgs(int64(100), int64(2)).WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	got, err := repo.CreateRoutine(10, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 100 || got.UserID != 10 || got.Name != "Push" || len(got.ExerciseIDs) != 2 || !got.IsActive {
		t.Fatalf("unexpected routine: %#v", got)
	}
}

func TestRoutineRepository_CreateRoutine_InsertExerciseFails_RollsBack(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	req := model.CreateRoutineRequest{Name: "A", Description: "D", ExerciseIDs: []int64{1}}

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO workout_routine").
		WithArgs("A", int64(5), "D").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "description"}).
			AddRow(77, "A", 5, "D"))
	mock.ExpectExec("INSERT INTO exercises_in_routine").
		WithArgs(int64(77), int64(1)).
		WillReturnError(errors.New("fk violation"))
	mock.ExpectRollback()

	_, err := repo.CreateRoutine(5, req)
	if err == nil || err.Error() != "fk violation" {
		t.Fatalf("expected fk violation, got %v", err)
	}
}

func TestRoutineRepository_ReadUserRoutines_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, user_id, description FROM workout_routine WHERE user_id = $1",
	)).WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "description"}).
			AddRow(1, "A", 10, "d1").
			AddRow(2, "B", 10, "d2"))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1",
	)).WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id"}).AddRow(int64(5)).AddRow(int64(6)))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1",
	)).WithArgs(int64(2)).
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id"}))

	got, err := repo.ReadUserRoutines(10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].ID != 1 || len(got[0].ExerciseIDs) != 2 || !got[1].IsActive {
		t.Fatalf("unexpected routines: %#v", got)
	}
}

func TestRoutineRepository_ReadUserRoutines_OuterQueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery("SELECT id, name, user_id, description FROM workout_routine").
		WithArgs(int64(9)).
		WillReturnError(errors.New("db down"))

	_, err := repo.ReadUserRoutines(9)
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down, got %v", err)
	}
}

func TestRoutineRepository_ReadUserRoutines_InnerQueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery("SELECT id, name, user_id, description FROM workout_routine").
		WithArgs(int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "description"}).
			AddRow(1, "A", 9, "d1"))

	mock.ExpectQuery("SELECT exercise_id FROM exercises_in_routine").
		WithArgs(int64(1)).
		WillReturnError(errors.New("join failed"))

	_, err := repo.ReadUserRoutines(9)
	if err == nil || err.Error() != "join failed" {
		t.Fatalf("expected join failed, got %v", err)
	}
}

func TestRoutineRepository_ReadRoutineWithExercises_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, user_id, description FROM workout_routine WHERE id = $1",
	)).WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "user_id", "description"}).
			AddRow(7, "Pull", 2, "d"))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT exercise_id FROM exercises_in_routine WHERE workout_routine_id = $1",
	)).WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"exercise_id"}).AddRow(int64(11)).AddRow(int64(12)))

	got, err := repo.ReadRoutineWithExercises(7)
	if err != nil || got.ID != 7 || len(got.ExerciseIDs) != 2 {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestRoutineRepository_ReadRoutineWithExercises_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectQuery("SELECT id, name, user_id, description FROM workout_routine WHERE id = \\$1").
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.ReadRoutineWithExercises(99)
	if err == nil || err.Error() != "routine not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRoutineRepository_DeleteRoutine_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM workout_routine WHERE id = $1")).
		WithArgs(int64(3)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteRoutine(3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoutineRepository_DeleteRoutine_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("DELETE FROM workout_routine").
		WithArgs(int64(4)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := repo.DeleteRoutine(4); err == nil || err.Error() != "routine not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRoutineRepository_AddExerciseToRoutine_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO exercises_in_routine (workout_routine_id, exercise_id) VALUES ($1, $2)",
	)).WithArgs(int64(5), int64(9)).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.AddExerciseToRoutine(5, 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoutineRepository_AddExerciseToRoutine_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("INSERT INTO exercises_in_routine").
		WithArgs(int64(5), int64(9)).
		WillReturnError(errors.New("fk failed"))

	if err := repo.AddExerciseToRoutine(5, 9); err == nil || err.Error() != "fk failed" {
		t.Fatalf("expected fk failed, got %v", err)
	}
}

func TestRoutineRepository_RemoveExerciseFromRoutine_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM exercises_in_routine WHERE workout_routine_id = $1 AND exercise_id = $2",
	)).WithArgs(int64(5), int64(9)).WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.RemoveExerciseFromRoutine(5, 9); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRoutineRepository_RemoveExerciseFromRoutine_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRoutineRepository(db)

	mock.ExpectExec("DELETE FROM exercises_in_routine").
		WithArgs(int64(5), int64(9)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	if err := repo.RemoveExerciseFromRoutine(5, 9); err == nil || err.Error() != "exercise not found in routine" {
		t.Fatalf("expected not found, got %v", err)
	}
}
