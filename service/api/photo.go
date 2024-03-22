package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"strconv"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

/*
uploadPhoto Upload a photo to the website and create a new post.

	curl -X POST http://localhost:3000/users/USERNAME/photos -H 'Authorization: Bearer USER_ID' -H 'Content-Type: multipart/form-data' -F 'photo=@/path/to/photo' -F 'description=DESCRIPTION'
*/
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	photo.Owner = header

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // limit your maxMultipartMemory
	if handleError(w, "unable to parse multipart form", http.StatusInternalServerError, err) {
		return
	}

	// Get the photo file
	file, _, err := r.FormFile("image")
	if handleError(w, "unable to get photo file", http.StatusInternalServerError, err) {
		return
	}
	defer file.Close()

	// Read the photo file into a byte array
	photo.File, err = io.ReadAll(file)
	if handleError(w, "unable to read photo file", http.StatusInternalServerError, err) {
		return
	}

	// Get the description
	photo.Description = r.FormValue("description")

	// Upload the photo
	photo, err = rt.db.UploadPhoto(photo)
	if handleError(w, "unable to upload photo", http.StatusInternalServerError, err) {
		return
	}

	// Return the created Photo object in the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(photo)
	if handleError(w, "unable to encode photo", http.StatusInternalServerError, err) {
		return
	}
}

/*
deletePhoto Delete an uploaded photo from the website, along with its post.

	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID -H 'Authorization: Bearer USER_ID'
*/
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var photo Photo

	// Get the requesters data from the auth header
	header, err := parseAuthHeader(r.Header.Get("Authorization"))
	if handleError(w, "unable to parse auth header", http.StatusInternalServerError, err) {
		return
	}
	photo.Owner = header

	// Get the photo's data from the db
	photoIdUint64, err := strconv.ParseUint(ps.ByName("photoId"), 10, 64)
	if handleError(w, "unable to parse photoId", http.StatusInternalServerError, err) {
		return
	}
	photo.PhotoId = uint(photoIdUint64)

	// Delete the photo
	err = rt.db.DeletePhoto(photo)
	if handleError(w, "unable to delete photo", http.StatusInternalServerError, err) {
		return
	}

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
