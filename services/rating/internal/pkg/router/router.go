package router

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/log"
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
	r.Use(rt.ContextMiddleware, rt.RecoverMiddleware)

	r.Get("/", rt.healthCheckHandler)
	r.Get("/user/{id}", rt.findByUserHandler)
	r.Get("/movie/{id}", rt.findByMovieHandler)
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
