package repository

import "workoutpal/src/internal/model"

type RelationshipRepository interface {
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	ReadUserFollowers(userID int64) ([]int64, error)
	ReadUserFollowing(userID int64) ([]int64, error)
	
	// Follow request methods
	CreateFollowRequest(requesterID, requestedID int64) error
	GetFollowRequest(requesterID, requestedID int64) (*model.FollowRequestModel, error)
	GetFollowRequestByID(requestID int64) (*model.FollowRequestModel, error)
	GetPendingFollowRequests(userID int64) ([]*model.FollowRequestWithUser, error)
	UpdateFollowRequestStatus(requestID int64, status string) error
	DeleteFollowRequest(requestID int64) error
}
