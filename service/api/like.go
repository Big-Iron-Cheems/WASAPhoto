package api

import (
	"encoding/json"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

/*
likePhoto Add a like to a photo.

	curl -X POST BASE_URL/users/USERNAME/photos/PHOTO_ID/likes -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var liker User
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	liker.UserId = header
	liker.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, liker.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo, err = rt.db.GetPhoto(uint(photoIdUint64))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Like the photo
	err = rt.db.LikePhoto(liker.UserId, photo.PhotoId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photo.LikeCount++

	// Return the photo info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
unlikePhoto Remove a like from a photo.

	curl -X DELETE BASE_URL/users/USERNAME/photos/PHOTO_ID/likes/LIKER_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var unliker User
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	unliker.UserId = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Unlike the photo
	err = rt.db.UnlikePhoto(unliker.UserId, photo.PhotoId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

/*
getLikeStatus Get the like status of a user for a photo.

	curl -X GET BASE_URL/users/USERNAME/photos/PHOTO_ID/likes/list/TARGET_USERNAME -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) getLikeStatus(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var photo Photo
	var targetUser User

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	targetUser.UserId = header
	targetUser.Username = ps.ByName("targetUsername")

	// Validate the target username
	if err = validateString(usernamePattern, targetUser.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the username
	user.Username = ps.ByName("username")
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the photo's ID from the path
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Check if the target user has liked the photo
	hasLiked, err := rt.db.GetLikeStatus(targetUser.UserId, photo.PhotoId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the like status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]bool{"hasLiked": hasLiked})
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
