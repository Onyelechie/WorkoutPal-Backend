package service

import "workoutpal/src/internal/model"

type RelationshipService interface {
	FollowUser(followerID, followeeID int64) error
	UnfollowUser(followerID, followeeID int64) error
	ReadUserFollowers(userID int64) ([]model.User, error)
	ReadUserFollowing(userID int64) ([]model.User, error)
	
	// Follow request methods
	SendFollowRequest(requesterID, requestedID int64) error
	GetFollowRequest(requesterID, requestedID int64) (*model.FollowRequestModel, error)
	GetPendingFollowRequests(userID int64) ([]*model.FollowRequestWithUser, error)
	AcceptFollowRequest(requestID int64) error
	RejectFollowRequest(requestID int64) error
	CancelFollowRequest(requesterID, requestedID int64) error
}
