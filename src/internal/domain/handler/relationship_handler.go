package handler

import "net/http"

type RelationshipHandler interface {
	ReadFollowers(w http.ResponseWriter, r *http.Request)
	ReadFollowings(w http.ResponseWriter, r *http.Request)
	FollowUser(w http.ResponseWriter, r *http.Request)
	UnfollowUser(w http.ResponseWriter, r *http.Request)
	
	// Follow request endpoints
	SendFollowRequest(w http.ResponseWriter, r *http.Request)
	GetPendingFollowRequests(w http.ResponseWriter, r *http.Request)
	RespondToFollowRequest(w http.ResponseWriter, r *http.Request)
	CancelFollowRequest(w http.ResponseWriter, r *http.Request)
	GetFollowRequestStatus(w http.ResponseWriter, r *http.Request)
}
