package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGoalRepository_CreateGoal_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	req := model.CreateGoalRequest{
		Name:        "New Goal",
		Description: "Test goal",
		Deadline:    "2025-01-01",
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "deadline", "created_at", "status"}).
		AddRow(1, 10, req.Name, req.Description, req.Deadline, "2025-01-01T00:00:00Z", "active")

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO goals (user_id, name, description, deadline, status) 
		VALUES ($1, $2, $3, $4, 'active') 
		RETURNING id, user_id, name, description, deadline, created_at, status`)).
		WithArgs(int64(10), req.Name, req.Description, req.Deadline).
		WillReturnRows(rows)

	got, err := repo.CreateGoal(10, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != 1 || got.UserID != 10 || got.Status != "active" {
		t.Fatalf("unexpected model: %#v", got)
	}
}

func TestGoalRepository_CreateGoal_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	mock.ExpectQuery("INSERT INTO goals").
		WillReturnError(errors.New("db exploded"))

	_, err := repo.CreateGoal(10, model.CreateGoalRequest{})
	if err == nil || err.Error() != "db exploded" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestGoalRepository_ReadUserGoals_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "deadline", "created_at", "status"}).
		AddRow(1, 10, "Goal 1", "Desc 1", "2025-01-01", "ts1", "active").
		AddRow(2, 10, "Goal 2", "Desc 2", "2026-01-01", "ts2", "paused")

	mock.ExpectQuery("SELECT id, user_id, name, description, deadline, created_at, status FROM goals WHERE user_id = \\$1").
		WithArgs(int64(10)).
		WillReturnRows(rows)

	got, err := repo.ReadUserGoals(10)
	if err != nil || len(got) != 2 {
		t.Fatalf("unexpected result: %+v err=%v", got, err)
	}
}

func TestGoalRepository_ReadUserGoals_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	// wrong column type
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "deadline", "created_at", "status"}).
		AddRow("BAD", 10, "Goal", "Desc", "2025", "ts", "active")

	mock.ExpectQuery("SELECT id, user_id").
		WithArgs(int64(10)).
		WillReturnRows(rows)

	_, err := repo.ReadUserGoals(10)
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
}

func TestGoalRepository_UpdateGoal_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	req := model.UpdateGoalRequest{
		ID:          1,
		Name:        "Updated",
		Description: "New Desc",
		Deadline:    "2026",
		Status:      "paused",
	}

	row := sqlmock.NewRows([]string{"id", "user_id", "name", "description", "deadline", "created_at", "status"}).
		AddRow(1, 10, "Updated", "New Desc", "2026", "ts", "paused")

	mock.ExpectQuery("UPDATE goals").
		WithArgs(req.ID, req.Name, req.Description, req.Deadline, req.Status).
		WillReturnRows(row)

	got, err := repo.UpdateGoal(req)
	if err != nil || got.Status != "paused" {
		t.Fatalf("unexpected: %+v err=%v", got, err)
	}
}

func TestGoalRepository_UpdateGoal_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	mock.ExpectQuery("UPDATE goals").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateGoal(model.UpdateGoalRequest{ID: 99})
	if err == nil || err.Error() != "goal not found" {
		t.Fatalf("expected 'goal not found', got %v", err)
	}
}

func TestGoalRepository_DeleteGoal_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	mock.ExpectExec("DELETE FROM goals").
		WithArgs(int64(1)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteGoal(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGoalRepository_DeleteGoal_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewGoalRepository(db)

	mock.ExpectExec("DELETE FROM goals").
		WithArgs(int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 0)) // zero rows

	err := repo.DeleteGoal(2)
	if err == nil || err.Error() != "goal not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}
