package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/log"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

func (rt *router) findHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	m, err := rt.repository.FindByID(ctx, id)
	if err != nil {
		rt.logger.Error("movie not found", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, m, http.StatusOK)
}

func (rt *router) createHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	var p createPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	m := domain.NewMovie(p.Title, p.OriginalTitle, p.Poster, p.Genres)

	vm, err := domain.NewValidatedMovie(m)
	if err != nil {
		rt.logger.Error("invalid movie data", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	m, err = rt.repository.Create(ctx, vm)
	if err != nil {
		rt.logger.Error("failed to create movie", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, m, http.StatusCreated)
}
