package model

type Achievement struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
}

type UserAchievement struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Username    string `json:"username"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"date"`
}

type CreateAchievementRequest struct {
	UserID        int64 `json:"userId"`
	AchievementID int64 `json:"achievementId"`
}
