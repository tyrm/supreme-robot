package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	errBadLogin            = errors.New("username/password combo invalid")
	errNotFound            = errors.New("not found")
	errRefreshExpired      = errors.New("refresh expired")
	errUnauthorized        = errors.New("unauthorized")
	errUnprocessableEntity = errors.New("unprocessable entity")
)

func (s *Server) returnErrorPage(w http.ResponseWriter, status int, errStr string) {
	errorBody := make(map[string]map[string]interface{})
	errorBody["error"] = map[string]interface{}{
		"status": status,
		"title":  http.StatusText(status),
		"detail": errStr,
	}

	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(errorBody)
	if err != nil {
		logger.Errorf("could marshal json: %s", err.Error())
	}
}

func (s *Server) methodNotAllowedHandler() http.Handler {
	// wrap in middleware since middleware isn't run on error pages
	return s.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, http.StatusMethodNotAllowed, r.Method)
	}))
}

func (s *Server) notFoundHandler() http.Handler {
	// wrap in middleware since middleware isn't run on error pages
	return s.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, http.StatusNotFound, fmt.Sprintf("page not found: %s", r.URL.Path))
	}))
}
