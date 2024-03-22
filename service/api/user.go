package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

/*
getUserProfile Given a user's username, retrieve all the public info available.

	curl -X GET http://localhost:3000/users/USERNAME/profile -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var requester User
	var profile Profile

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}

	requester.UserId = header
	requester.Username = ps.ByName("username")

	// Get the user's data from the db
	user.Username = ps.ByName("username")
	user, err = rt.db.GetUserProfile(user)
	if handleError(w, "unable to get user", http.StatusInternalServerError, err) {
		return
	}

	// Retrieve the PUBLIC data related to this user's profile
	user, err = rt.db.GetUserProfile(user)
	if handleError(w, "unable to get user", http.StatusInternalServerError, err) {
		return
	}

	profile.UserId = user.UserId
	profile.Username = user.Username

	photoCount, err := rt.db.GetPhotoCount(user.UserId)
	if handleError(w, "unable to get photo count", http.StatusInternalServerError, err) {
		return
	}
	profile.PhotoCount = photoCount

	followersCount, err := rt.db.GetFollowersCount(user.UserId)
	if handleError(w, "unable to get followers count", http.StatusInternalServerError, err) {
		return
	}
	profile.FollowersCount = followersCount

	followingCount, err := rt.db.GetFollowingCount(user.UserId)
	if handleError(w, "unable to get following count", http.StatusInternalServerError, err) {
		return
	}
	profile.FollowingCount = followingCount

	// Return the profile schema as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)
	if handleError(w, "unable to encode profile", http.StatusInternalServerError, err) {
		return
	}
}

/*
getMyStream Given a user's id, retrieve the content stream.
The stream is composed of entries that have images, likes and comments.
These entries are sorted in reverse chronological order.

	curl -X GET http://localhost:3000/users/USERNAME/stream -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getMyStream(w http.ResponseWriter, _ *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	panic("not implemented") // TODO
}

/*
setMyUsername Given a user's username, update its username.

	curl -X PUT http://localhost:3000/users/USERNAME -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"username": "NEW_USERNAME"}'
*/
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	currentUsername := ps.ByName("username")

	// Get the new username from the request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if handleError(w, "unable to decode request body", http.StatusUnprocessableEntity, err) {
		return
	}

	// Update the username
	user, err = rt.db.SetMyUsername(user, currentUsername)
	if handleError(w, "unable to update username", http.StatusInternalServerError, err) {
		return
	}

	// Return the user schema as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if handleError(w, "unable to encode user", http.StatusInternalServerError, err) {
		return
	}
}
