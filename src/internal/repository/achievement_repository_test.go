package repository

import (
	sql2 "database/sql"
	"errors"
	"regexp"
	"testing"

	"workoutpal/src/internal/model"

	"github.com/DATA-DOG/go-sqlmock"
)

func colsAll() []string {
	return []string{"id", "title", "badge_icon", "description"}
}

func colsUnlocked() []string {
	return []string{"id", "user_id", "title", "badge_icon", "description", "earned_at"}
}

// -------------------- ReadAllAchievements --------------------

func Test_ReadAllAchievements_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, a.title, a.badge_icon, a.description
    FROM achievements a`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WillReturnRows(sqlmock.NewRows(colsAll()).
			AddRow(int64(1), "First Workout", "first.png", "Finish your first workout").
			AddRow(int64(2), "7-Day Streak", "streak7.png", "Train 7 days in a row"),
		)

	list, err := repo.ReadAllAchievements()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("want 2, got %d", len(list))
	}
	if list[0].ID != 1 || list[0].Title != "First Workout" {
		t.Fatalf("unexpected first row: %#v", list[0])
	}
}

func Test_ReadAllAchievements_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, a.title, a.badge_icon, a.description
    FROM achievements a`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).WillReturnError(errors.New("db down"))

	_, err := repo.ReadAllAchievements()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_ReadAllAchievements_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, a.title, a.badge_icon, a.description
    FROM achievements a`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WillReturnRows(sqlmock.NewRows(colsAll()).
			AddRow("bad-id", "X", "I", "D"),
		)

	_, err := repo.ReadAllAchievements()
	if err == nil {
		t.Fatalf("expected scan error")
	}
}

// -------------------- ReadUnlockedAchievements --------------------

func Test_ReadUnlockedAchievements_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.user_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(42)).
		WillReturnRows(sqlmock.NewRows(colsUnlocked()).
			AddRow(int64(9), int64(42), "First Workout", "first.png", "Finish your first workout", "2025-01-02T03:04:05Z"),
		)

	list, err := repo.ReadUnlockedAchievements(42)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("want 1, got %d", len(list))
	}
	got := list[0]
	if got.ID != 9 || got.UserID != 42 || got.Title != "First Workout" || got.EarnedAt != "2025-01-02T03:04:05Z" {
		t.Fatalf("unexpected row: %#v", got)
	}
}

func Test_ReadUnlockedAchievements_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.user_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(7)).
		WillReturnError(errors.New("query fail"))

	_, err := repo.ReadUnlockedAchievements(7)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_ReadUnlockedAchievements_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.user_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(colsUnlocked()).
			AddRow("oops", int64(1), "X", "I", "D", "2025-01-01T00:00:00Z"),
		)

	_, err := repo.ReadUnlockedAchievements(1)
	if err == nil {
		t.Fatalf("expected scan error")
	}
}

// -------------------- ReadUnlockedAchievementByAchievementID --------------------

func Test_ReadUnlockedAchievementByAchievementID_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.achievement_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(55)).
		WillReturnRows(sqlmock.NewRows(colsUnlocked()).
			AddRow(int64(55), int64(7), "7-Day Streak", "streak7.png", "Train 7 days in a row", "2025-02-03T04:05:06Z"),
		)

	got, err := repo.ReadUnlockedAchievementByAchievementID(55)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 55 || got.UserID != 7 || got.Title != "7-Day Streak" || got.EarnedAt != "2025-02-03T04:05:06Z" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func Test_ReadUnlockedAchievementByAchievementID_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.achievement_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(999)).
		WillReturnError(errors.New("boom"))

	_, err := repo.ReadUnlockedAchievementByAchievementID(999)
	if err == nil {
		t.Fatalf("expected error")
	}
}

func Test_ReadUnlockedAchievementByAchievementID_NoRows(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	sql := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.achievement_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(sql)).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows(colsUnlocked()))

	_, err := repo.ReadUnlockedAchievementByAchievementID(1)
	if !errors.Is(err, sql2.ErrNoRows) {
		t.Fatalf("expected sql.ErrNoRows, got %v", err)
	}
}

// -------------------- CreateAchievement --------------------

func Test_CreateAchievement_OK(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	req := model.CreateAchievementRequest{
		UserID:        7,
		AchievementID: 55,
	}

	insertSQL := `
		INSERT INTO user_achievements (user_id, achievement_id, earned_at)
		VALUES ($1,$2,now())`
	mock.ExpectExec(regexp.QuoteMeta(insertSQL)).
		WithArgs(req.UserID, req.AchievementID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	selectSQL := `
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.achievement_id = $1
    ORDER BY ua.earned_at DESC`

	mock.ExpectQuery(regexp.QuoteMeta(selectSQL)).
		WithArgs(req.AchievementID).
		WillReturnRows(sqlmock.NewRows(colsUnlocked()).
			AddRow(int64(55), int64(7), "First Workout", "first.png", "Finish your first workout", "2025-01-01T00:00:00Z"),
		)

	got, err := repo.CreateAchievement(req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got == nil || got.ID != 55 || got.UserID != 7 || got.Title != "First Workout" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func Test_CreateAchievement_InsertError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	repo := NewAchievementRepository(db)

	req := model.CreateAchievementRequest{UserID: 1, AchievementID: 2}

	insertSQL := `
		INSERT INTO user_achievements (user_id, achievement_id, earned_at)
		VALUES ($1,$2,now())`
	mock.ExpectExec(regexp.QuoteMeta(insertSQL)).
		WithArgs(req.UserID, req.AchievementID).
		WillReturnError(errors.New("insert fail"))

	_, err := repo.CreateAchievement(req)
	if err == nil {
		t.Fatalf("expected error")
	}
}
