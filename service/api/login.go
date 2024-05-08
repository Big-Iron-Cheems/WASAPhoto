package api

import (
	"encoding/json"
	"github.com/Big-Iron-Cheems/WASAPhoto/service/api/reqcontext"
	. "github.com/Big-Iron-Cheems/WASAPhoto/service/model"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
doLogin If the user does not exist, it will be created, and an identifier is returned.
If the user exists, the user identifier is returned.

	curl -X POST http://localhost:3000/session -H 'Content-Type: application/json' -d '{"username": "USERNAME"}'
*/
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	var user User

	// Decode the request body into the user schema
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the username
	if err = validateString(usernamePattern, user.Username); err != nil {
		respondWithJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user exists or log in the user
	user, err = rt.db.CreateUser(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the user schema as response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		respondWithJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
