package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExerciseRepository_ReadExerciseByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "targets", "image", "demo"}).
		AddRow(42, "Deadlift", "posterior chain", "back,glutes,hamstrings", "img.png", "demo.mp4")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises WHERE id = $1",
	)).
		WithArgs(int64(42)).
		WillReturnRows(rows)

	got, err := repo.ReadExerciseByID(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil || got.ID != 42 || got.Name != "Deadlift" {
		t.Fatalf("unexpected exercise: %#v", got)
	}
	if len(got.Targets) != 3 || got.Targets[0] != "back" || got.Image != "img.png" || got.Demo != "demo.mp4" {
		t.Fatalf("unexpected fields: %#v", got)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestExerciseRepository_ReadExerciseByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises WHERE id = $1",
	)).
		WithArgs(int64(7)).
		WillReturnError(sql.ErrNoRows)

	got, err := repo.ReadExerciseByID(7)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "exercise not found" {
		t.Fatalf("expected 'exercise not found', got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestExerciseRepository_ReadExerciseByID_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises WHERE id = $1",
	)).
		WithArgs(int64(9)).
		WillReturnError(assertErr)

	got, err := repo.ReadExerciseByID(9)
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != assertErr.Error() {
		t.Fatalf("expected %q, got %v", assertErr.Error(), err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestExerciseRepository_ReadAllExercises_OK_WithNulls(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "targets", "image", "demo"}).
		AddRow(1, "Push-ups", "desc1", "chest,shoulders,triceps", "img1.png", nil).
		AddRow(2, "Pull-ups", "desc2", "back,biceps", nil, "demo2.mp4")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises",
	)).WillReturnRows(rows)

	got, err := repo.ReadAllExercises()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("want 2 rows, got %d", len(got))
	}
	if got[0].ID != 1 || got[0].Name != "Push-ups" || len(got[0].Targets) != 3 || got[0].Image != "img1.png" || got[0].Demo != "" {
		t.Fatalf("bad row1: %#v", got[0])
	}
	if got[1].ID != 2 || got[1].Name != "Pull-ups" || len(got[1].Targets) != 2 || got[1].Image != "" || got[1].Demo != "demo2.mp4" {
		t.Fatalf("bad row2: %#v", got[1])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestExerciseRepository_ReadAllExercises_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "targets", "image", "demo"}).
		AddRow("bad", "X", "Y", "a,b", nil, nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises",
	)).WillReturnRows(rows)

	got, err := repo.ReadAllExercises()
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil {
		t.Fatalf("expected scan error, got nil")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
func TestExerciseRepository_ReadAllExercises_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewExerciseRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "description", "targets", "image", "demo"}).
		AddRow(1, "A", "d", "x,y", nil, nil).
		AddRow(2, "B", "e", "p,q", nil, nil).
		RowError(1, assertErr)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, name, description, targets, image, demo FROM exercises",
	)).WillReturnRows(rows)

	got, err := repo.ReadAllExercises()
	if got != nil {
		t.Fatalf("expected nil slice, got %#v", got)
	}
	if err == nil || err.Error() != assertErr.Error() {
		t.Fatalf("expected %q, got %v", assertErr.Error(), err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

var assertErr = &testErr{"boom"}

type testErr struct{ s string }

func (e *testErr) Error() string { return e.s }
