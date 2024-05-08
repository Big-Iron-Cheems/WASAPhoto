package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// maxPageSize is the maximum number of users that can be retrieved in a single request.
const maxPageSize = 100

/*
getAllUsers retrieves all users from the database via paginated requests.

	curl -X GET http://localhost:3000/users?page=1&pageSize=50 -H 'Authorization Bearer USER_ID'
*/
func (rt *_router) getAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	var page, pageSize int
	var err error

	// Validate the page number
	if pageStr == "" {
		page = 1 // Default value
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			respondWithJSONError(w, "Page number must be a valid integer", http.StatusBadRequest)
			return
		}
	}
	if page < 1 {
		respondWithJSONError(w, "page number must be greater than 0", http.StatusBadRequest)
		return
	}

	// Validate the page size
	if pageSizeStr == "" {
		pageSize = 50 // Default value
	} else {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if err != nil {
			respondWithJSONError(w, "Page size must be a valid integer", http.StatusBadRequest)
			return
		}
	}
	if pageSize < 1 || pageSize > maxPageSize {
		respondWithJSONError(w, fmt.Sprintf("page size must be between 1 and %d", maxPageSize), http.StatusBadRequest)
		return
	}

	users, err := rt.db.GetAllUsers(page, pageSize)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	requester.UserId = header
	requester.Username = ps.ByName("username")

	// Get the user's data from the db
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists
	user, err = rt.db.GetUserProfile(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	profile.UserId = user.UserId
	profile.Username = user.Username

	// Fetch the user's photo count, followers count, following count and banned count
	photoCount, err := rt.db.GetPhotoCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.PhotoCount = photoCount

	followersCount, err := rt.db.GetFollowersCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.FollowersCount = followersCount

	followingCount, err := rt.db.GetFollowingCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.FollowingCount = followingCount

	bannedCount, err := rt.db.GetBansCount(user.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.BannedCount = bannedCount

	// Return the profile schema as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(profile)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
getMyStream Given a user's id, retrieve the content stream.
The stream is composed of Photo entries.
These entries are sorted in reverse chronological order.

	curl -X GET http://localhost:3000/stream -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header

	// Get the stream
	stream, err := rt.db.GetMyStream(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the stream as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(stream)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
setMyUsername Given a user's username, update its username.

	curl -X PUT http://localhost:3000/users/USERNAME -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"username": "NEW_USERNAME"}'
*/
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	currentUsername := ps.ByName("username")

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.UserId = header

	// Check if the username in the path is of an existing user in the database
	_, err = rt.db.GetUserProfile(User{Username: currentUsername})
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Get the new username from the request body
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the new username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the username
	user, err = rt.db.SetMyUsername(user, currentUsername)
	if err != nil {
		if errors.Is(err, &UsernameTakenError{Username: user.Username}) {
			respondWithJSONError(w, err.Error(), http.StatusConflict)
		} else {
			respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the user schema as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
