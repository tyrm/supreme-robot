package graphql

import (
	"github.com/gorilla/mux"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/scheduler"
	"net/http"
	"time"
)

type Server struct {
	// data stuff
	db        *models.Client
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

func NewServer(cfg *config.Config, s *scheduler.Client, d *models.Client, r *redis.Client) (*Server, error) {
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
	server.router.Use(server.Middleware)

	// Error Pages
	server.router.NotFoundHandler = server.NotFoundHandler()
	server.router.MethodNotAllowedHandler = server.MethodNotAllowedHandler()

	server.router.HandleFunc("/graphql", server.graphqlHandler).Methods("POST")

	return &server, nil
}

func (s *Server) Close() {
	err := s.server.Close()
	if err != nil {
		logger.Warningf("closing server: %s", err.Error())
	}
}

func (s *Server) ListenAndServe() error {
	s.server = &http.Server{
		Handler:      s.router,
		Addr:         ":5000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return s.server.ListenAndServe()
}
