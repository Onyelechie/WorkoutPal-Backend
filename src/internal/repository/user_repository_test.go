package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"workoutpal/src/internal/model"
)

func TestUserRepository_ReadUserByEmail_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "password", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow(
		1, "max", "a@b.com", "hashed", "Max", 25,
		180, "cm", 75.0, "kg", "https://img/avatar.png",
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE email = $1",
	)).WithArgs("a@b.com").WillReturnRows(rows)

	got, err := repo.ReadUserByEmail("a@b.com")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 1 || got.Username != "max" || got.Avatar != "https://img/avatar.png" {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserRepository_ReadUserByEmail_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE email = $1",
	)).WithArgs("x@y.com").WillReturnError(sql.ErrNoRows)

	got, err := repo.ReadUserByEmail("x@y.com")
	if got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestUserRepository_ReadUserByEmail_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, username").
		WithArgs("a@b.com").
		WillReturnError(errors.New("db down"))

	_, err := repo.ReadUserByEmail("a@b.com")
	if err == nil || err.Error() != "db down" {
		t.Fatalf("expected db down, got %v", err)
	}
}

func TestUserRepository_ReadUsers_OK_IncludingNullAvatar(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow(1, "max", "a@b.com", "Max", 25, 180, "cm", 75.0, "kg", "https://img/a.png").
		AddRow(2, "sam", "c@d.com", "Sam", 30, 175, "cm", 70.0, "kg", nil)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_url FROM users",
	)).WillReturnRows(rows)

	got, err := repo.ReadUsers()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(got) != 2 || got[0].Username != "max" || got[0].Avatar == "" || got[1].Avatar != "" {
		t.Fatalf("unexpected users: %#v", got)
	}
}

func TestUserRepository_ReadUsers_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, username, email, name").
		WillReturnError(errors.New("query fail"))

	_, err := repo.ReadUsers()
	if err == nil || err.Error() != "query fail" {
		t.Fatalf("expected query fail, got %v", err)
	}
}

func TestUserRepository_ReadUsers_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow("bad", "u", "e", "n", 25, 170, "cm", 60.0, "kg", nil)

	mock.ExpectQuery("SELECT id, username, email, name").
		WillReturnRows(rows)

	_, err := repo.ReadUsers()
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
}

func TestUserRepository_ReadUserByID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow(7, "max", "a@b.com", "Max", 25, 180, "cm", 75.0, "kg", "https://img/a.png")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_url FROM users WHERE id = $1",
	)).WithArgs(int64(7)).WillReturnRows(rows)

	got, err := repo.ReadUserByID(7)
	if err != nil || got == nil || got.ID != 7 || got.Avatar == "" {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestUserRepository_ReadUserByID_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, username, email, name").
		WithArgs(int64(99)).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.ReadUserByID(99)
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestUserRepository_ReadUserByID_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("SELECT id, username, email, name").
		WithArgs(int64(3)).
		WillReturnError(errors.New("db fail"))

	_, err := repo.ReadUserByID(3)
	if err == nil || err.Error() != "db fail" {
		t.Fatalf("expected db fail, got %v", err)
	}
}

func TestUserRepository_CreateUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	req := model.CreateUserRequest{
		Username:     "max",
		Email:        "a@b.com",
		Password:     "hashed",
		Name:         "Max",
		Height:       180,
		HeightMetric: "cm",
		Weight:       75.0,
		WeightMetric: "kg",
		Avatar:       "https://img/a.png",
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow(1, req.Username, req.Email, req.Name,
		req.Height, req.HeightMetric, req.Weight, req.WeightMetric, req.Avatar)

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO users (username, email, password, name, height, height_metric, weight, weight_metric, avatar_url) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id, username, email, name, height, height_metric, weight, weight_metric, avatar_url`)).
		WithArgs(req.Username, req.Email, req.Password, req.Name, req.Height, req.HeightMetric, req.Weight, req.WeightMetric, req.Avatar).
		WillReturnRows(rows)

	got, err := repo.CreateUser(req)
	if err != nil || got == nil || got.ID != 1 || got.Avatar != req.Avatar {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestUserRepository_CreateUser_Duplicate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("INSERT INTO users").
		WillReturnError(&pq.Error{Code: pq.ErrorCode("23505")})

	_, err := repo.CreateUser(model.CreateUserRequest{})
	if err == nil || err.Error() != "user already exists" {
		t.Fatalf("expected 'user already exists', got %v", err)
	}
}

func TestUserRepository_CreateUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("INSERT INTO users").
		WillReturnError(errors.New("insert fail"))

	_, err := repo.CreateUser(model.CreateUserRequest{})
	if err == nil || err.Error() != "insert fail" {
		t.Fatalf("expected insert fail, got %v", err)
	}
}

func TestUserRepository_UpdateUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	req := model.UpdateUserRequest{
		ID:           10,
		Username:     "newname",
		Email:        "new@e.com",
		Name:         "New",
		Age:          26,
		Height:       181,
		HeightMetric: "cm",
		Weight:       76.0,
		WeightMetric: "kg",
		Avatar:       "https://img/new.png",
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_url",
	}).AddRow(req.ID, req.Username, req.Email, req.Name, req.Age,
		req.Height, req.HeightMetric, req.Weight, req.WeightMetric, req.Avatar)

	mock.ExpectQuery(regexp.QuoteMeta(`
		UPDATE users SET username=$2, email=$3, name=$4, age=$5, height=$6, height_metric=$7, weight=$8, weight_metric=$9, avatar_url=$10
		WHERE id=$1 RETURNING id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_url`)).
		WithArgs(req.ID, req.Username, req.Email, req.Name, req.Age, req.Height, req.HeightMetric, req.Weight, req.WeightMetric, req.Avatar).
		WillReturnRows(rows)

	got, err := repo.UpdateUser(req)
	if err != nil || got == nil || got.ID != 10 || got.Avatar != req.Avatar {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
}

func TestUserRepository_UpdateUser_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("UPDATE users").
		WillReturnError(sql.ErrNoRows)

	_, err := repo.UpdateUser(model.UpdateUserRequest{ID: 7})
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestUserRepository_UpdateUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery("UPDATE users").
		WillReturnError(errors.New("update fail"))

	_, err := repo.UpdateUser(model.UpdateUserRequest{ID: 7})
	if err == nil || err.Error() != "update fail" {
		t.Fatalf("expected update fail, got %v", err)
	}
}

func TestUserRepository_DeleteUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM users WHERE id = $1",
	)).WithArgs(int64(22)).WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.DeleteUser(model.DeleteUserRequest{ID: 22})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUserRepository_DeleteUser_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec("DELETE FROM users").
		WithArgs(int64(22)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteUser(model.DeleteUserRequest{ID: 22})
	if err == nil || err.Error() != "user not found" {
		t.Fatalf("expected user not found, got %v", err)
	}
}

func TestUserRepository_DeleteUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectExec("DELETE FROM users").
		WithArgs(int64(22)).
		WillReturnError(errors.New("delete fail"))

	err := repo.DeleteUser(model.DeleteUserRequest{ID: 22})
	if err == nil || err.Error() != "delete fail" {
		t.Fatalf("expected delete fail, got %v", err)
	}
}
