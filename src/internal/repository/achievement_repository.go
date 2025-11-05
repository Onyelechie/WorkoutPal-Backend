package repository

import (
	"database/sql"
	domainrepo "workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/model"
)

type achievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) domainrepo.AchievementRepository {
	return &achievementRepository{db: db}
}

func (r *achievementRepository) ReadAchievementsFeed() ([]*model.UserAchievement, error) {
	rows, err := r.db.Query(`
    SELECT a.id, u.username, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    JOIN users u ON u.id = ua.user_id
    ORDER BY ua.earned_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.UserAchievement = make([]*model.UserAchievement, 0)
	for rows.Next() {
		var a model.UserAchievement
		if err := rows.Scan(&a.ID, &a.Username, &a.UserID, &a.Title, &a.BadgeIcon, &a.Description, &a.EarnedAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *achievementRepository) ReadAllAchievements() ([]*model.Achievement, error) {
	rows, err := r.db.Query(`
    SELECT a.id, a.title, a.badge_icon, a.description
    FROM achievements a`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Achievement = make([]*model.Achievement, 0)
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.Title, &a.BadgeIcon, &a.Description); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *achievementRepository) ReadUnlockedAchievements(userID int64) ([]*model.UserAchievement, error) {
	rows, err := r.db.Query(`
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.user_id = $1
    ORDER BY ua.earned_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.UserAchievement = make([]*model.UserAchievement, 0)
	for rows.Next() {
		var a model.UserAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.Title, &a.BadgeIcon, &a.Description, &a.EarnedAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *achievementRepository) CreateAchievement(req model.CreateAchievementRequest) (*model.UserAchievement, error) {
	_, err := r.db.Exec(`
		INSERT INTO user_achievements (user_id, achievement_id, earned_at)
		VALUES ($1,$2,now())`,
		req.UserID, req.AchievementID,
	)
	if err != nil {
		return nil, err
	}

	return r.ReadUnlockedAchievementByAchievementID(req.AchievementID)
}

func (r *achievementRepository) ReadUnlockedAchievementByAchievementID(id int64) (*model.UserAchievement, error) {
	row := r.db.QueryRow(`
    SELECT a.id, ua.user_id, a.title, a.badge_icon, a.description, ua.earned_at
    FROM achievements a
    JOIN user_achievements ua ON ua.achievement_id = a.id
    WHERE ua.achievement_id = $1
    ORDER BY ua.earned_at DESC`, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var a model.UserAchievement
	if err := row.Scan(&a.ID, &a.UserID, &a.Title, &a.BadgeIcon, &a.Description, &a.EarnedAt); err != nil {
		return nil, err
	}

	return &a, nil
}
