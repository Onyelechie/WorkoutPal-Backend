package repository

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"workoutpal/src/internal/model"
)

func TestAchievementRepository_ReadAchievements_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, user_id, title, badge_icon, description, created_at
		FROM achievements
		WHERE user_id = $1
		ORDER BY created_at DESC`,
	)).WillReturnRows(sqlmock.NewRows(
		[]string{"id", "user_id", "title", "badge_icon", "description", "created_at"},
	).AddRow(int64(1), int64(2), "T", "I", "D", "2025-01-01T00:00:00Z"))

	got, err := repo.ReadAchievements(2)
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if len(got) != 1 || got[0].ID != 1 || got[0].Title != "T" {
		t.Fatalf("unexpected: %+v", got)
	}
}

func TestAchievementRepository_ReadAchievements_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	mock.ExpectQuery("SELECT id, user_id, title").WillReturnError(errors.New("db down"))

	_, err := repo.ReadAchievements(0)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementRepository_CreateAchievement_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	req := model.CreateAchievementRequest{UserID: 1, Title: "T", BadgeIcon: "I", Description: "D", EarnedAt: "2025-01-01T00:00:00Z"}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO achievements (user_id, title, badge_icon, description, created_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, user_id, title, badge_icon, description, created_at`,
	)).WithArgs(req.UserID, req.Title, req.BadgeIcon, req.Description, req.EarnedAt).
		WillReturnRows(sqlmock.NewRows(
			[]string{"id", "user_id", "title", "badge_icon", "description", "created_at"},
		).AddRow(int64(9), int64(1), "T", "I", "D", "2025-01-01T00:00:00Z"))

	got, err := repo.CreateAchievement(req)
	if err != nil || got.ID != 9 || got.Title != "T" {
		t.Fatalf("unexpected: %+v err=%v", got, err)
	}
}

func TestAchievementRepository_CreateAchievement_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	req := model.CreateAchievementRequest{UserID: 1, Title: "T", BadgeIcon: "I", Description: "D", EarnedAt: "2025-01-01T00:00:00Z"}

	mock.ExpectQuery("INSERT INTO achievements").WithArgs(req.UserID, req.Title, req.BadgeIcon, req.Description, req.EarnedAt).
		WillReturnError(errors.New("insert fail"))

	_, err := repo.CreateAchievement(req)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestAchievementRepository_DeleteAchievement_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM achievements WHERE id = $1`)).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))

	if err := repo.DeleteAchievement(5); err != nil {
		t.Fatalf("err=%v", err)
	}
}

func TestAchievementRepository_DeleteAchievement_Error(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	mock.ExpectExec("DELETE FROM achievements").
		WithArgs(int64(6)).
		WillReturnError(errors.New("fail"))

	if err := repo.DeleteAchievement(6); err == nil {
		t.Fatalf("expected error")
	}
}
