package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/domain"
)

// Router is the interface for the application server http handler.
type Router interface {
	GetHandler() http.Handler
}

type router struct {
	repository domain.Repository
	logger     *log.Logger
}

// New returns a new instance of Router.
func New(repo domain.Repository, logger *log.Logger) Router {
	return &router{repo, logger}
}

// GetHandler returns the router's http handler.
func (rt *router) GetHandler() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(rt.ContextMiddleware, middleware.Recoverer)

	r.Get("/", rt.healthCheckHandler)
	r.Get("/{id}", rt.findHandler)
	r.Post("/create", rt.createHandler)

	return r
}

func getRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(ctxKeyRequestID).(string)
	if !ok {
		requestID = "unknown"
	}
	return requestID
}
