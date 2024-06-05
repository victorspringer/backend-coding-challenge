package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/rating/internal/pkg/domain"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

// @Summary Find ratings by user ID
// @Description Get all ratings given by a specific user
// @Tags ratings
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Param id path string true "User ID"
// @Produce json
// @Success 200 {object} response{response=[]domain.Rating}
// @Failure 401 {object} response
// @Failure 404 {object} response
// @Router /user/{id} [get]
func (rt *router) findByUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level == "anonymous" {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	userID := chi.URLParam(r, "id")

	u, err := rt.repository.FindByUserID(ctx, userID)
	if err != nil {
		rt.logger.Error("rating not found", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
}

// @Summary Find ratings by movie ID
// @Description Get all ratings for a specific movie
// @Tags ratings
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Param id path string true "Movie ID"
// @Produce json
// @Success 200 {object} response{response=[]domain.Rating}
// @Failure 401 {object} response
// @Failure 404 {object} response
// @Router /movie/{id} [get]
func (rt *router) findByMovieHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level == "anonymous" {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	movieID := chi.URLParam(r, "id")

	u, err := rt.repository.FindByMovieID(ctx, movieID)
	if err != nil {
		rt.logger.Error("rating not found", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
}

// @Summary Create a new (or override an old) rating
// @Description Create a new (or override an old) rating for a movie by a user
// @Tags ratings
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Accept json
// @Produce json
// @Param rating body upsertPayload true "Rating"
// @Success 200 {object} response{response=domain.Rating}
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /upsert [post]
func (rt *router) upsertHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level == "anonymous" {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	var p upsertPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.String("body", string(b)), log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	if username := context.GetUserUsername(ctx); username == "" || username != p.UserID {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	rat := domain.NewRating(p.UserID, p.MovieID, p.Value)

	vr, err := domain.NewValidatedRating(rat)
	if err != nil {
		rt.logger.Error("invalid rating data", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	rat, err = rt.repository.Upsert(ctx, vr)
	if err != nil {
		rt.logger.Error("failed to create / update rating", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, rat, http.StatusOK)
}
