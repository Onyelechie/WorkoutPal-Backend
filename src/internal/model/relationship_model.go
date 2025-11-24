package model

import "time"

type FollowRequest struct {
	UserID int64 `json:"userID"`
}

type UnfollowRequest struct {
	UserID int64 `json:"userID"`
}

type FollowRequestModel struct {
	ID          int64     `json:"id"`
	RequesterID int64     `json:"requesterID"`
	RequestedID int64     `json:"requestedID"`
	Status      string    `json:"status"` // pending, accepted, rejected
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type FollowRequestWithUser struct {
	ID          int64     `json:"id"`
	RequesterID int64     `json:"requesterID"`
	RequestedID int64     `json:"requestedID"`
	Status      string    `json:"status"`
	User        *User     `json:"user"` // The user who sent the request
	CreatedAt   time.Time `json:"createdAt"`
}

type FollowRequestResponse struct {
	RequestID int64  `json:"requestID"`
	Action    string `json:"action"` // accept, reject
}
