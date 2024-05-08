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
getPhotoComments Get all comments under a photo.

	curl -X GET http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json'
*/
func (rt *_router) getPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var photo Photo

	// Get the requesters data from the auth header
	_, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Get the comments from the db
	comments, err := rt.db.GetPhotoComments(photo.PhotoId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the comments
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
commentPhoto Add a comment under a photo.

	curl -X POST http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments -H 'Authorization: Bearer USER_ID' -H 'Content-Type: application/json' -d '{"content": "CONTENT"}'
*/
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var photo Photo
	var user User
	var comment Comment

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.OwnerId = header

	// Validate the username
	if err = validateString(usernamePattern, ps.ByName("username")); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the commenter's data from the db
	user.UserId = header
	user, err = rt.db.GetUserProfile(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusNotFound)
		return
	}
	comment.OwnerUsername = user.Username

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Get the comment's data from the request body
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the comment
	if err = validateString(commentPattern, comment.Content); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Comment the photo
	comment, err = rt.db.CommentPhoto(photo.PhotoId, comment)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created comment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
uncommentPhoto Remove a comment from a photo.

	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments/COMMENT_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var photo Photo
	var comment Comment

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.OwnerId = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Get the comment's data from the db
	commentIdUint64, err := strconv.ParseUint(ps.ByName("commentId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	comment.CommentId = uint(commentIdUint64)

	// Uncomment the photo
	err = rt.db.UncommentPhoto(photo.PhotoId, comment)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
