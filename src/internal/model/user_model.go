package model

type User struct {
	ID           int64   `json:"id"`
	Username     string  `json:"username"`
	Email        string  `json:"email"`
	Password     string  `json:"-"`
	Name         string  `json:"name"`
	Avatar       string  `json:"avatar"`
	Age          int     `json:"age"`
	Height       float64 `json:"height"`
	HeightMetric string  `json:"heightMetric"`
	Weight       float64 `json:"weight"`
	WeightMetric string  `json:"weightMetric"`
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
