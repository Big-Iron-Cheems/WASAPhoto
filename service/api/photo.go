package api

import (
	"encoding/json"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"strings"
)

/*
getPhotoList Get the list of photos uploaded by a user.

	curl -X GET http://localhost:3000/users/USERNAME/photos -H 'Authorization Bearer USER_ID' -H 'Content-Type: application/json'
*/
func (rt *_router) getPhotoList(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	var targetUser User

	// Get the requesters data from the auth header
	_, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Username = ps.ByName("username")

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Get target user's ID
	targetUser, err = rt.db.GetUserProfile(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the photos from the db
	photos, err := rt.db.GetPhotoList(targetUser.UserId)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the photos in the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(photos)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
uploadPhoto Upload a photo to the website and create a new post.

	curl -X POST http://localhost:3000/users/USERNAME/photos -H 'Authorization: Bearer USER_ID' -H 'Content-Type: multipart/form-data' -F 'image=@/path/to/photo' -F 'caption=CAPTION'
*/
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.OwnerId = header

	// Validate the username
	if err = validateString(usernamePattern, ps.ByName("username")); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // limit your maxMultipartMemory
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the photo file
	file, _, err := r.FormFile("image")
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the photo file into a byte array
	photo.Image, err = io.ReadAll(file)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the mime type
	photo.MimeType = r.FormValue("mimeType")

	// Get the caption
	var caption = strings.TrimSpace(r.FormValue("caption"))
	// Validate the caption
	if len(caption) > 0 {
		if err = validateString(captionPattern, caption); err != nil {
			respondWithJSONError(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
	}
	photo.Caption = caption

	// Upload the photo
	photo, err = rt.db.UploadPhoto(photo)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the created Photo object in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/*
deletePhoto Delete an uploaded photo from the website, along with its comments and likes.

	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, _ reqcontext.RequestContext) {
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.OwnerId = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Delete the photo
	err = rt.db.DeletePhoto(photo)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
