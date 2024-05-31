package router

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/log"
)

// RecoverMiddleware recovers the application from panic during a request.
func (rt *router) RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				e, ok := err.(error)
				if !ok {
					err = fmt.Errorf("%#v\n%s", err, string(debug.Stack()))
				}

				rt.logger.Error("recovered handler panic", log.Error(e), log.String("requestId", getRequestID(r.Context())))

				debug.PrintStack()

				rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
