package graphql

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadLogin            = errors.New("username/password combo invalid")
	ErrRefreshExpired      = errors.New("refresh expired")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUnprocessableEntity = errors.New("unprocessable entity")
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

func (s *Server) MethodNotAllowedHandler() http.Handler {
	// wrap in middleware since middleware isn't run on error pages
	return s.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, http.StatusMethodNotAllowed, r.Method)
	}))
}

func (s *Server) NotFoundHandler() http.Handler {
	// wrap in middleware since middleware isn't run on error pages
	return s.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.returnErrorPage(w, http.StatusNotFound, fmt.Sprintf("page not found: %s", r.URL.Path))
	}))
}
