package handler

import "net/http"

type RelationshipHandler interface {
	ReadFollowers(w http.ResponseWriter, r *http.Request)
	ReadFollowings(w http.ResponseWriter, r *http.Request)
	FollowUser(w http.ResponseWriter, r *http.Request)
	UnfollowUser(w http.ResponseWriter, r *http.Request)
}
