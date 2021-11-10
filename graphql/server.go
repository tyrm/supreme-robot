package graphql

import (
	"github.com/gorilla/mux"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/scheduler"
	"net/http"
	"time"
)

// Server is a GraphQL api server
type Server struct {
	// data stuff
	db        db.DB
	redis     *redis.Client
	scheduler *scheduler.Client

	// dns stuff
	primaryNS string

	// web stuff
	accessExpiration  time.Duration
	accessSecret      []byte
	refreshExpiration time.Duration
	refreshSecret     []byte
	router            *mux.Router
	server            *http.Server
}

// NewServer will create a new GraphQL server
func NewServer(cfg *config.Config, s *scheduler.Client, d db.DB, r *redis.Client) (*Server, error) {
	server := Server{
		accessExpiration:  cfg.AccessExpiration,
		accessSecret:      []byte(cfg.AccessSecret),
		db:                d,
		primaryNS:         cfg.PrimaryNS,
		redis:             r,
		refreshExpiration: cfg.RefreshExpiration,
		refreshSecret:     []byte(cfg.RefreshSecret),
		scheduler:         s,
	}

	// Setup Router
	server.router = mux.NewRouter()
	server.router.Use(server.middleware)

	// Error Pages
	server.router.NotFoundHandler = server.notFoundHandler()
	server.router.MethodNotAllowedHandler = server.methodNotAllowedHandler()

	server.router.HandleFunc("/graphql", server.graphqlHandler).Methods("POST")

	return &server, nil
}

// Close will cleanly stop the server.
func (s *Server) Close() {
	err := s.server.Close()
	if err != nil {
		logger.Warningf("closing server: %s", err.Error())
	}
}

// ListenAndServe starts the web server.
func (s *Server) ListenAndServe() error {
	s.server = &http.Server{
		Handler:      s.router,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s.server.ListenAndServe()
}
