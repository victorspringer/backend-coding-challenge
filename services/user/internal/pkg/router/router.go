package router

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	_ "github.com/victorspringer/backend-coding-challenge/services/user/docs"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

// @title User Service
// @version 1.0
// @description User Service for Movie Rating System.
// @contact.name Victor Springer
// @license.name MIT License
// @host localhost:8081
// @BasePath /

// Router is the interface for the application server http handler.
type Router interface {
	GetHandler() http.Handler
}

type router struct {
	repository      domain.Repository
	logger          *log.Logger
	cacheMiddleware func(next http.Handler) http.Handler
}

// New returns a new instance of Router.
func New(repo domain.Repository, logger *log.Logger) Router {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}
	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return &router{repo, logger, cacheClient.Middleware}
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

	r.Use(middleware.StripSlashes, rt.ContextMiddleware, rt.cacheMiddleware, middleware.Recoverer)

	// health check
	r.Get("/", rt.healthCheckHandler)

	// docs
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})
	r.Get("/docs/*", httpSwagger.Handler())

	// endpoints
	r.Get("/{username}", rt.findHandler)
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
