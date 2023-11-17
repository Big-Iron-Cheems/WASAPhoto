package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// getUserProfile Given a user's username, retrieve all the public info available.
//
//	curl -X GET http://localhost:3000/users/USERNAME/profile
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// getMyStream Given a user's id, retrieve the content stream.
// The stream is composed of entries that have images, likes and comments.
// These entries are sorted in reverse chronological order.
//
//	curl -X GET http://localhost:3000/users/USERNAME/stream
func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// setMyUsername Given a user's id, update its username.
//
//	curl -X PUT http://localhost:3000/users/USERNAME -H 'Content-Type: application/json' -d '{"newUsername": "NEW_USERNAME"}'
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
