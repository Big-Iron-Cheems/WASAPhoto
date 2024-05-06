package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	. "wasaphoto.uniroma1.it/wasaphoto/service/model"
)

// usernamePattern is the regex pattern for a valid username.
const usernamePattern = `^[A-Za-z0-9_\-]{3,32}$`

// captionPattern is the regex pattern for a valid post caption.
const captionPattern = `^[\p{L}\p{N}\p{M}\p{P}\p{S} ]{0,32}$`

// commentPattern is the regex pattern for a valid post comment.
const commentPattern = `^[\p{L}\p{N}\p{M}\p{P}\p{S} \n]{1,256}$`

// parseAuthHeader is a helper function to parse the authorization header and return the user id.
func parseAuthHeader(header string) (uint, error) {
	// Remove the "Bearer " prefix
	token := strings.TrimPrefix(header, "Bearer ")
	if token == header {
		return 0, errors.New("invalid authorization header")
	}

	// Convert the token to an uint
	userId, err := strconv.ParseUint(token, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

/*
validateString validates a string against a regex pattern.
*/
func validateString(pattern string, str string) error {
	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		// Return error if regex matching fails
		return err
	}
	if !matched {
		// Return error if str doesn't match pattern
		return &InvalidPatternError{Pattern: pattern, Str: str}
	}
	// Return nil if validation succeeds
	return nil
}

/*
respondWithJSONError will respond to a request with a JSON error message.

This should be used whenever we'd call http.Error, but we want to return a JSON error message instead.
The JSON object returned is the model.Error struct.
*/
func respondWithJSONError(w http.ResponseWriter, message string, statusCode int) {
	errResponse := Error{Code: statusCode, Message: message}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Encode the error message as JSON
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding error message: %v", err), http.StatusInternalServerError)
	}
}
