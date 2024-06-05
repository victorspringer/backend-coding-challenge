package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	authClient "github.com/victorspringer/backend-coding-challenge/services/authentication/pkg/client"
	_ "github.com/victorspringer/backend-coding-challenge/services/rating/docs"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/domain"
)

// @title Rating Service
// @version 1.0
// @description Rating Service for Movie Rating System.
// @contact.name Victor Springer
// @license.name MIT License
// @host localhost:8082
// @BasePath /

// Router is the interface for the application server http handler.
type Router interface {
	GetHandler() http.Handler
}

type router struct {
	repository domain.Repository
	logger     *log.Logger
	ac         *authClient.Client
}

// New returns a new instance of Router.
func New(repo domain.Repository, logger *log.Logger, ac *authClient.Client) Router {
	return &router{repo, logger, ac}
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

	r.Use(middleware.StripSlashes, rt.ContextMiddleware, rt.ac.Middleware, middleware.Recoverer)

	//health check
	r.Get("/", rt.healthCheckHandler)

	// docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	r.Get("/docs/*", httpSwagger.Handler())

	// endpoints
	r.Get("/user/{id}", rt.findByUserHandler)
	r.Get("/movie/{id}", rt.findByMovieHandler)
	r.Post("/upsert", rt.upsertHandler)

	return r
}
