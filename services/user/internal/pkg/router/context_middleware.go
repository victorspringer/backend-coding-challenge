package router

import (
	"context"
	"net/http"
	"time"

	libCtx "github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

// ContextMiddleware adds data into the request's context (e.g. UUID).
func (rt *router) ContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = libCtx.GetRequestID(r.Context())
			r.Header.Set("X-Request-ID", requestID)
		}
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

		ctx := context.WithValue(r.Context(), libCtx.CTX_REQUEST_ID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
