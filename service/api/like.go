package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

/*
likePhoto Add a like to a photo.

	curl -X POST http://localhost:3000/users/USERNAME/photos/PHOTO_ID/likes -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var liker User
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	liker.UserId = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if handleError(w, "unable to parse photoId", http.StatusInternalServerError, err) {
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Like the photo
	err = rt.db.LikePhoto(liker.UserId, photo.PhotoId)
	if handleError(w, "unable to like photo", http.StatusInternalServerError, err) {
		return
	}

	// Get the updated photo data from the db
	photo, err = rt.db.GetPhoto(photo.PhotoId)
	if handleError(w, "unable to get photo", http.StatusInternalServerError, err) {
		return
	}

	// Return the photo info
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(photo)
	if handleError(w, "unable to encode photo", http.StatusInternalServerError, err) {
		return
	}
}

/*
unlikePhoto Remove a like from a photo.

	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID/likes/LIKER_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var unliker User
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	unliker.UserId = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if handleError(w, "unable to parse photoId", http.StatusInternalServerError, err) {
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Unlike the photo
	err = rt.db.UnlikePhoto(unliker.UserId, photo.PhotoId)
	if handleError(w, "unable to unlike photo", http.StatusInternalServerError, err) {
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
