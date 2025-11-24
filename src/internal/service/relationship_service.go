package service

import (
	"workoutpal/src/internal/domain/repository"
	"workoutpal/src/internal/domain/service"
	"workoutpal/src/internal/model"
)

type relationshipService struct {
	relationshipRepository repository.RelationshipRepository
	userRepository         repository.UserRepository
}

func NewRelationshipService(relationshipRepository repository.RelationshipRepository, userRepository repository.UserRepository) service.RelationshipService {
	return &relationshipService{
		relationshipRepository: relationshipRepository,
		userRepository:         userRepository,
	}
}

func (u *relationshipService) FollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.FollowUser(followerID, followeeID)
}

func (u *relationshipService) UnfollowUser(followerID, followeeID int64) error {
	return u.relationshipRepository.UnfollowUser(followerID, followeeID)
}

func (u *relationshipService) ReadUserFollowers(userID int64) ([]model.User, error) {
	followerIDs, err := u.relationshipRepository.ReadUserFollowers(userID)
	if err != nil {
		return nil, err
	}

	var followers []model.User
	for _, id := range followerIDs {
		user, err := u.userRepository.ReadUserByID(id)
		if err == nil {
			followers = append(followers, *user)
		}
	}
	return followers, nil
}

func (u *relationshipService) ReadUserFollowing(userID int64) ([]model.User, error) {
	followingIDs, err := u.relationshipRepository.ReadUserFollowing(userID)
	if err != nil {
		return nil, err
	}

	var following []model.User
	for _, id := range followingIDs {
		user, err := u.userRepository.ReadUserByID(id)
		if err == nil {
			following = append(following, *user)
		}
	}
	return following, nil
}

// Follow request methods

func (u *relationshipService) SendFollowRequest(requesterID, requestedID int64) error {
	return u.relationshipRepository.CreateFollowRequest(requesterID, requestedID)
}

func (u *relationshipService) GetFollowRequest(requesterID, requestedID int64) (*model.FollowRequestModel, error) {
	return u.relationshipRepository.GetFollowRequest(requesterID, requestedID)
}

func (u *relationshipService) GetPendingFollowRequests(userID int64) ([]*model.FollowRequestWithUser, error) {
	return u.relationshipRepository.GetPendingFollowRequests(userID)
}

func (u *relationshipService) AcceptFollowRequest(requestID int64) error {
	// Get the request details
	req, err := u.relationshipRepository.GetFollowRequestByID(requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return nil // Request not found
	}
	
	// Create the follow relationship
	err = u.relationshipRepository.FollowUser(req.RequesterID, req.RequestedID)
	if err != nil {
		return err
	}
	
	// Update status to accepted
	return u.relationshipRepository.UpdateFollowRequestStatus(requestID, "accepted")
}

func (u *relationshipService) RejectFollowRequest(requestID int64) error {
	return u.relationshipRepository.UpdateFollowRequestStatus(requestID, "rejected")
}

func (u *relationshipService) CancelFollowRequest(requesterID, requestedID int64) error {
	req, err := u.relationshipRepository.GetFollowRequest(requesterID, requestedID)
	if err != nil {
		return err
	}
	if req == nil {
		return nil // No request to cancel
	}
	return u.relationshipRepository.DeleteFollowRequest(req.ID)
}
