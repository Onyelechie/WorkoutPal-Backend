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
	Posts        []Post            `json:"posts,omitempty"`
	Followers    []int64           `json:"followers,omitempty"`
	Following    []int64           `json:"following,omitempty"`
	Goals        []Goal            `json:"goals,omitempty"`
	Achievements []UserAchievement `json:"achievements,omitempty"`
	Routines     []ExerciseRoutine `json:"routines,omitempty"`
}

type Goal struct {
	ID          int64   `json:"id"`
	UserID      int64   `json:"userId"`
	Type        string  `json:"type"` // "weight", "body_fat", "lift"
	TargetValue float64 `json:"targetValue"`
	CurrentValue float64 `json:"currentValue"`
	TargetDate  string  `json:"targetDate"`
	CreatedAt   string  `json:"createdAt"`
	Status      string  `json:"status"` // "active", "completed", "paused"
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
	Type        string  `json:"type"`
	TargetValue float64 `json:"targetValue"`
	TargetDate  string  `json:"targetDate"`
}

type UpdateGoalRequest struct {
	ID           int64   `json:"id"`
	CurrentValue float64 `json:"currentValue"`
	Status       string  `json:"status"`
}

type CreateRoutineRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ExerciseIDs []int64 `json:"exerciseIds"`
}
