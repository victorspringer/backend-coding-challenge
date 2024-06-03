package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/domain"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

func (rt *router) findByUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := chi.URLParam(r, "id")

	u, err := rt.repository.FindByUserID(ctx, userID)
	if err != nil {
		rt.logger.Error("rating not found", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
}

func (rt *router) findByMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	movieID := chi.URLParam(r, "id")

	u, err := rt.repository.FindByMovieID(ctx, movieID)
	if err != nil {
		rt.logger.Error("rating not found", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
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

	rat := domain.NewRating(p.UserID, p.MovieID, p.Value)

	vr, err := domain.NewValidatedRating(rat)
	if err != nil {
		rt.logger.Error("invalid rating data", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	rat, err = rt.repository.Create(ctx, vr)
	if err != nil {
		rt.logger.Error("failed to create rating", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, rat, http.StatusCreated)
}
