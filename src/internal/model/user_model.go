package model

type User struct {
	ID           int64             `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	Password     string            `json:"-"`
	Name         string            `json:"name"`
	Avatar       string            `json:"avatar"`
	Age          int               `json:"age"`
	Height       float64           `json:"height"`
	HeightMetric string            `json:"heightMetric"`
	Weight       float64           `json:"weight"`
	WeightMetric string            `json:"weightMetric"`
	GoogleID     string            `json:"googleId,omitempty"`
	Provider     string            `json:"provider,omitempty"`
	IsVerified   bool              `json:"isVerified"`
	Posts        []Post            `json:"posts,omitempty"`
	Followers    []int64           `json:"followers,omitempty"`
	Following    []int64           `json:"following,omitempty"`
	Goals        []Goal            `json:"goals,omitempty"`
	Achievements []UserAchievement `json:"achievements,omitempty"`
	Routines     []ExerciseRoutine `json:"routines,omitempty"`
}

type Goal struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	CreatedAt   string `json:"createdAt"`
	Status      string `json:"status"` // "active", "completed", "paused"
}

type UserAchievement struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	BadgeIcon   string `json:"badgeIcon"`
	EarnedAt    string `json:"earnedAt"`
}

type ExerciseRoutine struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"userId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Exercises   []Exercise `json:"exercises"`
	ExerciseIDs []int64    `json:"exerciseIds"`
	CreatedAt   string     `json:"createdAt"`
	IsActive    bool       `json:"isActive"`
}

type CreateUserRequest struct {
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Avatar       string  `json:"avatar"`
	Age          int     `json:"age"`
	Height       float64 `json:"height"`
	HeightMetric string  `json:"heightMetric"`
	Weight       float64 `json:"weight"`
	WeightMetric string  `json:"weightMetric"`
}

type UpdateUserRequest struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Avatar       string  `json:"avatar"`
	Age          int     `json:"age"`
	Height       float64 `json:"height"`
	HeightMetric string  `json:"heightMetric"`
	Weight       float64 `json:"weight"`
	WeightMetric string  `json:"weightMetric"`
}

type DeleteUserRequest struct {
	ID int64 `json:"id"`
}

type CreateGoalRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

type UpdateGoalRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
	Status      string `json:"status"`
}

type CreateRoutineRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ExerciseIDs []int64 `json:"exerciseIds"`
}

type GoogleAuthRequest struct {
	IDToken string `json:"idToken"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
