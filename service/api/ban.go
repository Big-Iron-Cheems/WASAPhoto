package api

import (
	"encoding/json"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
getBansList Get the list of users a user has banned via username.

	curl -X GET http://localhost:3000/users/USERNAME/bans/list -H 'Authorization : Bearer USER_ID'
*/
func (rt *_router) getBansList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header
	user.Username = ps.ByName("username")

	// Get the bans list
	bans, err := rt.db.GetBansList(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the bans list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bans)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getBanStatus Get the ban status of a user via username.
That is, check if the user in the 1st param has banned the user in the 2nd param.

	curl -X GET http://localhost:3000/users/USERNAME/bans/list/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getBanStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
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
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the 2nd username
	targetUser.Username = ps.ByName("targetUsername")
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the 1st user (owner of the ban list) data from the db
	user, err = rt.db.GetUserProfile(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get the 2nd user (target of the ban check) data from the db
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Check if the target user (2nd) is banned
	isBanned, err := rt.db.GetBanStatus(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ban status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]bool{"isBanned": isBanned})
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
banUser Add a user to your ban list via username.

	curl -X POST http://localhost:3000/users/USERNAME/bans -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"username": "TARGET_USERNAME"}'
*/
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
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

	// Validate the target username (path)
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the target user's username from the request body
	err = json.NewDecoder(r.Body).Decode(&targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the target username (body)
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the target user's data from the db
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Ban the target user
	err = rt.db.BanUser(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the banned user's info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(targetUser.Username)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
unbanUser Remove a user from your ban list via username.

	curl -X DELETE http://localhost:3000/users/USERNAME/bans/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
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

	// Validate the username (path)
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	targetUser.Username = ps.ByName("targetUsername")
	// Validate the target username (path)
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the target user's data from the db
	targetUser, err = rt.db.GetUserProfile(targetUser)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Unban the target user
	err = rt.db.UnbanUser(user.UserId, targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
