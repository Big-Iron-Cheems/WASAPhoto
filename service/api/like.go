package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// likePhoto Add a like to a photo.
//
//	curl -X POST http://localhost:3000/users/USERNAME/photos/PHOTO_ID/likes -H 'Content-Type: application/json' -d '{"username": "USERNAME"}'
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// unlikePhoto Remove a like from a photo.
//
// curl -X DELETE http://localhost:3000/users/USERNAME/photos/PHOTO_ID/likes/LIKE_ID
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
