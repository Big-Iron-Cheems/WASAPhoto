package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

/*
banUser Add a user to your ban list via username.

	curl -X POST http://localhost:3000/users/USERNAME/bans -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"username": "TARGET_USERNAME"}'
*/
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	user.UserId = header

	// get the target user's username from the request body
	err = json.NewDecoder(r.Body).Decode(&targetUser)
	if handleError(w, "unable to decode target user", http.StatusInternalServerError, err) {
		return
	}
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if handleError(w, "unable to get target user", http.StatusInternalServerError, err) {
		return
	}

	// Ban the target user
	err = rt.db.BanUser(user.UserId, targetUser.UserId)
	if handleError(w, "unable to ban target user", http.StatusInternalServerError, err) {
		return
	}

	// TODO: remove the comments and like of the target user from the requesters photos

	// Return the banned user's info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(targetUser.Username)
	if handleError(w, "unable to encode target user", http.StatusInternalServerError, err) {
		return
	}
}

/*
unbanUser Remove a user from your ban list via username.

	curl -X DELETE http://localhost:3000/users/USERNAME/bans/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Get the target user's data from the db
	targetUser.Username = ps.ByName("targetUsername")
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if handleError(w, "unable to get target user", http.StatusInternalServerError, err) {
		return
	}

	// Unban the target user
	err = rt.db.UnbanUser(user.UserId, targetUser.UserId)
	if handleError(w, "unable to unban target user", http.StatusInternalServerError, err) {
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
