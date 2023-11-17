package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
)

// banUser Add a user to your ban list via username.
//
//	curl -X POST http://localhost:3000/users/USERNAME/bans -H 'Content-Type: application/json' -d '{"targetUsername": "TARGET_USERNAME"}'
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}

// unbanUser Remove a user from your ban list via username.
//
//	curl -X DELETE http://localhost:3000/users/USERNAME/bans/TARGET_USERNAME
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

}
