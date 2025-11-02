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

func (r *achievementRepository) ReadAchievements(userID int64) ([]*model.Achievement, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, title, badge_icon, description, created_at
		FROM achievements
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Achievement
	for rows.Next() {
		var a model.Achievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.Title, &a.BadgeIcon, &a.Description, &a.EarnedAt); err != nil {
			return nil, err
		}
		result = append(result, &a)
	}
	return result, nil
}

func (r *achievementRepository) CreateAchievement(req model.CreateAchievementRequest) (*model.Achievement, error) {
	row := r.db.QueryRow(`
		INSERT INTO achievements (user_id, title, badge_icon, description, created_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, user_id, title, badge_icon, description, created_at`,
		req.UserID, req.Title, req.BadgeIcon, req.Description, req.EarnedAt,
	)

	var a model.Achievement
	if err := row.Scan(&a.ID, &a.UserID, &a.Title, &a.BadgeIcon, &a.Description, &a.EarnedAt); err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *achievementRepository) DeleteAchievement(id int64) error {
	_, err := r.db.Exec(`DELETE FROM achievements WHERE id = $1`, id)
	return err
}
