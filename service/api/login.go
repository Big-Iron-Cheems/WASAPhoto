package api

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"wasaphoto.uniroma1.it/wasaphoto/service/api/reqcontext"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// doLogin If the user does not exist, it will be created, and an identifier is returned.
// If the user exists, the user identifier is returned.
//
//	curl -X POST http://localhost:3000/session -H 'Content-Type: application/json' -d '{"username": "USERNAME"}'
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params, _ reqcontext.RequestContext) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, fmt.Errorf("unable to decode request body: %w", err).Error(), http.StatusUnprocessableEntity)
		return
	}

	user, err = rt.db.CreateUser(user)
	if err != nil {
		http.Error(w, fmt.Errorf("unable to create user: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// Return the user schema as response
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, fmt.Errorf("unable to encode user: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}
