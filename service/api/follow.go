package api

import (
	"encoding/json"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
getFollowersList Get the list of followers for a user via username.

	curl -X GET http://localhost:3000/users/USERNAME/followers/list -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getFollowersList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the followers list
	followers, err := rt.db.GetFollowersList(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the followers list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(followers)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getFollowersCount Get the count of followers for a user via username.

	curl -X GET http://localhost:3000/users/USERNAME/followers/count -H 'Authorization : Bearer USER_ID'
*/
func (rt *_router) getFollowersCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the followers count
	count, err := rt.db.GetFollowersCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the followers count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getFollowingList Get the list of users a user is following via username.

	curl -X GET http://localhost:3000/users/USERNAME/following/list -H 'Authorization : Bearer USER_ID'
*/
func (rt *_router) getFollowingList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the following list
	following, err := rt.db.GetFollowingList(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the following list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(following)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getFollowingCount Get the count of users a user is following via username.

	curl -X GET http://localhost:3000/users/USERNAME/following/count -H 'Authorization : Bearer USER_ID'
*/
func (rt *_router) getFollowingCount(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the following count
	count, err := rt.db.GetFollowingCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the following count
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(count)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getFollowStatus Get the follow status of a user via username.
That is, check if the user in the 1st param is following the user in the 2nd param.

	curl -X GET http://localhost:3000/users/USERNAME/followers/list/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getFollowStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	_, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the 1st username
	user.Username = ps.ByName("username")
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Validate the 2nd username
	targetUser.Username = ps.ByName("targetUsername")
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the 1st user (owner of the followers list) data from the db
	user, err = rt.db.GetUserProfile(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get the 2nd user (target of the follow check) data from the db
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if the target user (2nd) is followed
	isFollowing, err := rt.db.GetFollowStatus(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the follow status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]bool{"isFollowing": isFollowing})
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
followUser Add a user to your following list via username.

	curl -X POST http://localhost:3000/users/USERNAME/followers -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"username": "TARGET_USERNAME"}'
*/
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// get the target user's username from the request body
	err = json.NewDecoder(r.Body).Decode(&targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the target username
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Follow the target user
	err = rt.db.FollowUser(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the followed user's info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(targetUser.Username)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
unfollowUser Remove a user from your following list via username.

	curl -X DELETE http://localhost:3000/users/USERNAME/followers/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Validate the target username
	targetUser.Username = ps.ByName("targetUsername")
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get the target user's data from the db
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Unfollow the target user
	err = rt.db.UnfollowUser(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
