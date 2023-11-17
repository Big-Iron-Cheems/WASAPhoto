package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// commentPhoto Add a comment under a photo.
//
//	curl -X POST http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments -H 'Content-Type: application/json' -d '{"username": "USERNAME", "content": "CONTENT"}'
func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// uncommentPhoto Remove a comment from a photo.
//
//	curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID/comments/COMMENT_ID
func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
