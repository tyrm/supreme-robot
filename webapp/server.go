package webapp

import (
	"context"
	"encoding/gob"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/pkger"
	"github.com/rbcervilla/redisstore/v8"
	"github.com/tyrm/supreme-robot/config"
	"github.com/tyrm/supreme-robot/models"
	"github.com/tyrm/supreme-robot/redis"
	"github.com/tyrm/supreme-robot/scheduler"
	"github.com/tyrm/supreme-robot/startup"
	"html/template"
	"net/http"
	"time"

	redisCon "github.com/go-redis/redis/v8"
)

type Server struct {
	// data stuff
	config    config.Config
	db        *models.Client
	scheduler *scheduler.Client

	// web stuff
	store     *redisstore.RedisStore
	router    *mux.Router
	server    *http.Server
	templates *template.Template
}

func NewServer(scfg *startup.StartupConfig, s *scheduler.Client, d *models.Client, c config.Config) (*Server, error) {
	server := Server{
		config:    c,
		db:        d,
		scheduler: s,
	}

	// Load Templates
	templateDir := pkger.Include("/webapp/templates")
	t, err := compileTemplates(templateDir)
	if err != nil {
		return nil, err
	}
	server.templates = t

	// Redis client
	client := redisCon.NewClient(&redisCon.Options{
		Addr:     scfg.RedisWebappAddress,
		DB:       scfg.RedisWebappDB,
		Password: scfg.RedisWebappPassword,
	})

	// Fetch new store.
	server.store, err = redisstore.NewRedisStore(context.Background(), client)
	if err != nil {
		logger.Errorf("create redis store: %s", err.Error())
		return nil, err
	}

	server.store.KeyPrefix(redis.KeySession)
	server.store.Options(sessions.Options{
		Path:   "/",
		Domain: scfg.ExtHostname,
		MaxAge: 86400 * 60,
	})

	// Register models for GOB
	gob.Register(models.User{})
	gob.Register(templateAlert{})

	// Setup Router
	server.router = mux.NewRouter()
	server.router.Use(server.Middleware)

	// Error Pages
	server.router.NotFoundHandler = server.NotFoundHandler()
	server.router.MethodNotAllowedHandler = server.MethodNotAllowedHandler()

	// Static Files
	server.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(pkger.Dir("/webapp/static"))))

	server.router.HandleFunc("/login", server.LoginGetHandler).Methods("GET")
	server.router.HandleFunc("/login", server.LoginPostHandler).Methods("POST")
	server.router.HandleFunc("/logout", server.LogoutGetHandler).Methods("GET")

	// Protected Pages
	protected := server.router.PathPrefix("/app/").Subrouter()
	protected.Use(server.MiddlewareRequireAuth)
	protected.HandleFunc("/", server.HomeGetHandler).Methods("GET")
	protected.HandleFunc("/admin/users", server.AdminUsersGetHandler).Methods("GET")
	protected.HandleFunc("/admin/users/add", server.AdminUserAddGetHandler).Methods("GET")
	protected.HandleFunc("/admin/users/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/edit", server.AdminUserEditGetHandler).Methods("GET")
	protected.HandleFunc("/dns", server.DnsGetHandler).Methods("GET")
	protected.HandleFunc("/dns/add", server.DnsDomainAddGetHandler).Methods("GET")
	protected.HandleFunc("/dns/add", server.DnsDomainAddPostHandler).Methods("POST")
	protected.HandleFunc("/dns/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}", server.DnsDomainGetHandler).Methods("GET")
	protected.HandleFunc("/dns/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/delete", server.DnsDomainDeleteGetHandler).Methods("GET")
	protected.HandleFunc("/dns/{id:[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}}/delete", server.DnsDomainDeletePostHandler).Methods("POST")

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
