package model

type Achievement struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"earnedAt"`
}

type CreateAchievementRequest struct {
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	BadgeIcon   string `json:"badgeIcon"`
	Description string `json:"description"`
	EarnedAt    string `json:"earnedAt"`
}
