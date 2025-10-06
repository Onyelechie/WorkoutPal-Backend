package model

type FollowRequest struct {
	UserID int64 `json:"userID"`
}

type UnfollowRequest struct {
	UserID int64 `json:"userID"`
}
