package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	_ "github.com/victorspringer/backend-coding-challenge/services/authentication/docs"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/domain"
)

// @title Authentication Service
// @version 1.0
// @description Authentication Service for Movie Rating System.
// @contact.name Victor Springer
// @license.name MIT License
// @host localhost:8084
// @BasePath /

// Router is the interface for the application server http handler.
type Router interface {
	GetHandler() http.Handler
}

type router struct {
	authenticator domain.Authenticator
	logger        *log.Logger
}

// New returns a new instance of Router.
func New(auth domain.Authenticator, logger *log.Logger) Router {
	return &router{auth, logger}
}

// GetHandler returns the router's http handler.
func (rt *router) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID", "X-Forwarded-Proto"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(rt.ContextMiddleware, middleware.Recoverer)

	//health check
	r.Get("/", rt.healthCheckHandler)

	// docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	r.Get("/docs/*", httpSwagger.Handler())

	// endpoints
	r.Post("/anonymous", rt.loginAnonymous)
	r.Post("/login", rt.login)
	r.Post("/refresh", rt.refresh)
	r.Post("/logout", rt.logout)
	r.Post("/validate", rt.validateAccessToken)
	r.Get("/.well-known/jwks.json", rt.jwks)

	return r
}
