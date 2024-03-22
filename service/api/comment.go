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
commentPhoto Add a comment under a photo.

	curl -X POST http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"content": "CONTENT"}'
*/
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var photo Photo
	var comment Comment

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	comment.Owner = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if handleError(w, "unable to parse photoId", http.StatusInternalServerError, err) {
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Get the comment's data from the request body
	err = json.NewDecoder(r.Body).Decode(&comment)
	if handleError(w, "unable to decode comment", http.StatusInternalServerError, err) {
		return
	}

	// Comment the photo
	comment, err = rt.db.CommentPhoto(photo.PhotoId, comment)
	if handleError(w, "unable to comment photo", http.StatusInternalServerError, err) {
		return
	}

	// Return the created comment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(comment)
	if handleError(w, "unable to encode comment", http.StatusInternalServerError, err) {
		return
	}
}

/*
uncommentPhoto Remove a comment from a photo.

	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments/COMMENT_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var photo Photo
	var comment Comment

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	comment.Owner = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if handleError(w, "unable to parse photoId", http.StatusInternalServerError, err) {
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Get the comment's data from the db
	commentIdUint64, err := strconv.ParseUint(ps.ByName("commentId"), 10, 64)
	if handleError(w, "unable to parse commentId", http.StatusInternalServerError, err) {
		return
	}
	comment.CommentId = uint(commentIdUint64)

	// Uncomment the photo
	err = rt.db.UncommentPhoto(photo.PhotoId, comment)
	if handleError(w, "unable to uncomment photo", http.StatusInternalServerError, err) {
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
