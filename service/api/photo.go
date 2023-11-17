package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// uploadPhoto Upload a photo to the website and create a new post.
//
//	curl -X POST http://localhost:3000/users/USERNAME/photos -H 'Content-Type: multipart/form-data' -F 'photo=@/path/to/photo' -F 'description=DESCRIPTION'
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// deletePhoto Delete an uploaded photo from the website, along with its post.
//
//	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
