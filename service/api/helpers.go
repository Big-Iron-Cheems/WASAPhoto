package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
	// TODO: verify that the userid is in the authorised db
	return uint(userId), nil
}

// handleError is a helper function to handle errors in the api package.
func handleError(w http.ResponseWriter, customMessage string, statusCode int, err error) bool {
	if err != nil {
		http.Error(w, fmt.Sprintf("%s: %v", customMessage, err), statusCode)
		return true
	}
	return false
}
