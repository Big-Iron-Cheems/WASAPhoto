package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// followUser Add a user to your following list via username.
//
//	curl -X POST http://localhost:3000/users/USERNAME/follows -H 'Content-Type: application/json' -d '{"targetUsername": "TARGET_USERNAME"}'
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// unfollowUser Remove a user from your following list via username.
//
//	curl -X DELETE http://localhost:3000/users/USERNAME/follows/TARGET_USERNAME
func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
