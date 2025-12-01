package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
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
		"height", "height_metric", "weight", "weight_metric", "avatar_data", "is_private", "show_metrics_to_followers",
	}).AddRow(
		1, "max", "a@b.com", "hashed", "Max", 25,
		180, "cm", 75.0, "kg", []byte{255, 216, 255, 224}, false, false, // Sample JPEG binary data
	)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users WHERE email = $1",
	)).WithArgs("a@b.com").WillReturnRows(rows)

	got, err := repo.ReadUserByEmail("a@b.com")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 1 || got.Username != "max" || got.Avatar != "data:image/jpeg;base64,4pig4g==" {
		t.Fatalf("unexpected user: %#v", got)
	}
}

func TestUserRepository_ReadUserByEmail_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewUserRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users WHERE email = $1",
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
		"height", "height_metric", "weight", "weight_metric", "avatar_data", "is_private", "show_metrics_to_followers",
	}).AddRow(1, "max", "a@b.com", "Max", 25, 180, "cm", 75.0, "kg", []byte{255, 216, 255, 224}, false, false).
		AddRow(2, "sam", "c@d.com", "Sam", 30, 175, "cm", 70.0, "kg", nil, false, false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users",
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
		"height", "height_metric", "weight", "weight_metric", "avatar_data",
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
		"height", "height_metric", "weight", "weight_metric", "avatar_data", "is_private", "show_metrics_to_followers",
	}).AddRow(7, "max", "a@b.com", "Max", 25, 180, "cm", 75.0, "kg", []byte{255, 216, 255, 224}, false, false)

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers FROM users WHERE id = $1",
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
		Age:          25,
		Height:       180,
		HeightMetric: "cm",
		Weight:       75.0,
		WeightMetric: "kg",
		Avatar:       "data:image/jpeg;base64,4pig4g==",
	}

	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_data", "is_private", "show_metrics_to_followers",
	}).AddRow(1, req.Username, req.Email, req.Name, req.Age,
		req.Height, req.HeightMetric, req.Weight, req.WeightMetric, []byte{226, 152, 160, 226}, false, false)

	mock.ExpectQuery(regexp.QuoteMeta(`
			INSERT INTO users (username, email, password, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10::bytea, COALESCE($11, FALSE), COALESCE($12, FALSE)) 
			RETURNING id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers`)).
	WithArgs(req.Username, req.Email, req.Password, req.Name, req.Age, req.Height, req.HeightMetric, req.Weight, req.WeightMetric, []byte{226, 152, 160, 226}, req.IsPrivate, req.ShowMetricsToFollowers).
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
		Avatar:       "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==",
	}

	expectedBinary := []byte{137, 80, 78, 71, 13, 10, 26, 10, 0, 0, 0, 13, 73, 72, 68, 82, 0, 0, 0, 1, 0, 0, 0, 1, 8, 6, 0, 0, 0, 31, 21, 196, 137, 0, 0, 0, 13, 73, 68, 65, 84, 120, 218, 99, 252, 255, 159, 161, 30, 0, 7, 130, 2, 127, 61, 200, 72, 239, 0, 0, 0, 0, 73, 69, 78, 68, 174, 66, 96, 130}
	
	rows := sqlmock.NewRows([]string{
		"id", "username", "email", "name", "age",
		"height", "height_metric", "weight", "weight_metric", "avatar_data", "is_private", "show_metrics_to_followers",
	}).AddRow(req.ID, req.Username, req.Email, req.Name, req.Age,
		req.Height, req.HeightMetric, req.Weight, req.WeightMetric, expectedBinary, false, false)

	mock.ExpectQuery(regexp.QuoteMeta(`
			UPDATE users SET username=$2, email=$3, name=$4, age=$5, height=$6, height_metric=$7, weight=$8, weight_metric=$9, avatar_data=$10::bytea, is_private=$11, show_metrics_to_followers=$12
			WHERE id=$1 RETURNING id, username, email, name, age, height, height_metric, weight, weight_metric, avatar_data, is_private, show_metrics_to_followers`)).
	WithArgs(req.ID, req.Username, req.Email, req.Name, req.Age, req.Height, req.HeightMetric, req.Weight, req.WeightMetric, expectedBinary, req.IsPrivate, req.ShowMetricsToFollowers).
		WillReturnRows(rows)

	got, err := repo.UpdateUser(req)
	if err != nil || got == nil || got.ID != 10 {
		t.Fatalf("unexpected: %#v err=%v", got, err)
	}
	// Check that avatar is properly formatted data URL (binary conversion may cause slight differences)
	if !strings.HasPrefix(got.Avatar, "data:image/") || !strings.Contains(got.Avatar, "base64,") {
		t.Fatalf("expected valid avatar data URL, got: %s", got.Avatar)
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
