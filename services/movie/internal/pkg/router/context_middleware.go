package router

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/log"
)

type contextKey string

const (
	ctxKeyRequestID contextKey = "requestID"
)

// ContextMiddleware adds data into the request's context (e.g. UUID).
func (rt *router) ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		start := time.Now()
		defer func() {
			duration := time.Since(start)
			rt.logger.Debug(
				"request finished",
				log.String("requestId", requestID),
				log.String("duration", duration.String()),
			)
		}()

		ctx := context.WithValue(r.Context(), ctxKeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
