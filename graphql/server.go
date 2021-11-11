package graphql

import (
	"github.com/gorilla/mux"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/db"
	"github.com/tyrm/supreme-robot/kv"
	"github.com/tyrm/supreme-robot/queue"
	"net/http"
	"time"
)

// Server is a GraphQL api server
type Server struct {
	// data stuff
	db        db.DB
	kv        kv.Webapp
	scheduler queue.Scheduler

	// dns stuff
	primaryNS string

	// web stuff
	accessExpiration  time.Duration
	accessSecret      []byte
	port              string
	refreshExpiration time.Duration
	refreshSecret     []byte
	router            *mux.Router
	server            *http.Server
}

// NewServer will create a new GraphQL server
func NewServer(cfg *config.Config, s queue.Scheduler, d db.DB, k kv.Webapp) (*Server, error) {
	server := Server{
		accessExpiration:  cfg.AccessExpiration,
		accessSecret:      []byte(cfg.AccessSecret),
		db:                d,
		kv:                k,
		port:              cfg.HTTPPort,
		primaryNS:         cfg.PrimaryNS,
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
		Addr:         s.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s.server.ListenAndServe()
}
