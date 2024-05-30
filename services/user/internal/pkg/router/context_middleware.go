package router

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const (
	ctxKeyRequestID contextKey = "requestID"
)

// ContextMiddleware generates requests UUIDs.
func (rt *router) ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		w.Header().Set("X-Request-ID", requestID)

		ctx := context.WithValue(r.Context(), ctxKeyRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
