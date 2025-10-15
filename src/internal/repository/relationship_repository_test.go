package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestRelationshipRepository_FollowUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO follows (following_user_id, followed_user_id, created_at) VALUES ($1, $2, NOW())",
	)).
		WithArgs(int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.FollowUser(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelationshipRepository_FollowUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectExec("INSERT INTO follows").
		WillReturnError(errors.New("insert failed"))

	err := repo.FollowUser(1, 2)
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("expected insert failed, got %v", err)
	}
}

func TestRelationshipRepository_UnfollowUser_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM follows WHERE following_user_id = $1 AND followed_user_id = $2",
	)).
		WithArgs(int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.UnfollowUser(1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelationshipRepository_UnfollowUser_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectExec("DELETE FROM follows").
		WithArgs(int64(1), int64(2)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UnfollowUser(1, 2)
	if err == nil || err.Error() != "follow relationship not found" {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRelationshipRepository_UnfollowUser_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectExec("DELETE FROM follows").
		WithArgs(int64(1), int64(2)).
		WillReturnError(errors.New("delete failed"))

	err := repo.UnfollowUser(1, 2)
	if err == nil || err.Error() != "delete failed" {
		t.Fatalf("expected delete failed, got %v", err)
	}
}

func TestRelationshipRepository_ReadUserFollowers_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	rows := sqlmock.NewRows([]string{"following_user_id"}).
		AddRow(int64(3)).
		AddRow(int64(4))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT following_user_id FROM follows WHERE followed_user_id = $1",
	)).
		WithArgs(int64(2)).
		WillReturnRows(rows)

	got, err := repo.ReadUserFollowers(2)
	if err != nil || len(got) != 2 {
		t.Fatalf("unexpected result: %#v err=%v", got, err)
	}
}

func TestRelationshipRepository_ReadUserFollowers_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectQuery("SELECT following_user_id").
		WithArgs(int64(2)).
		WillReturnError(errors.New("query failed"))

	_, err := repo.ReadUserFollowers(2)
	if err == nil || err.Error() != "query failed" {
		t.Fatalf("expected query failed, got %v", err)
	}
}

func TestRelationshipRepository_ReadUserFollowers_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	rows := sqlmock.NewRows([]string{"following_user_id"}).
		AddRow("bad-id")

	mock.ExpectQuery("SELECT following_user_id").
		WithArgs(int64(2)).
		WillReturnRows(rows)

	_, err := repo.ReadUserFollowers(2)
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
}

func TestRelationshipRepository_ReadUserFollowing_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	rows := sqlmock.NewRows([]string{"followed_user_id"}).
		AddRow(int64(7)).
		AddRow(int64(8))

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT followed_user_id FROM follows WHERE following_user_id = $1",
	)).
		WithArgs(int64(5)).
		WillReturnRows(rows)

	got, err := repo.ReadUserFollowing(5)
	if err != nil || len(got) != 2 {
		t.Fatalf("unexpected result: %#v err=%v", got, err)
	}
}

func TestRelationshipRepository_ReadUserFollowing_DBError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	mock.ExpectQuery("SELECT followed_user_id").
		WithArgs(int64(5)).
		WillReturnError(errors.New("query failed"))

	_, err := repo.ReadUserFollowing(5)
	if err == nil || err.Error() != "query failed" {
		t.Fatalf("expected query failed, got %v", err)
	}
}

func TestRelationshipRepository_ReadUserFollowing_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewRelationshipRepository(db)

	rows := sqlmock.NewRows([]string{"followed_user_id"}).
		AddRow("bad-id")

	mock.ExpectQuery("SELECT followed_user_id").
		WithArgs(int64(5)).
		WillReturnRows(rows)

	_, err := repo.ReadUserFollowing(5)
	if err == nil {
		t.Fatal("expected scan error, got nil")
	}
}
