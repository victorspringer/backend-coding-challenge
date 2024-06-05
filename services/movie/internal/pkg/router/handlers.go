package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/domain"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

// @Summary Get movie by ID
// @Description Get movie information by ID
// @ID get-movie-by-id
// @Param id path string true "ID of the movie"
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Produce json
// @Success 200 {object} response{response=domain.Movie}
// @Failure 401 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /{id} [get]
func (rt *router) findHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level == "anonymous" {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	m, err := rt.repository.FindByID(ctx, id)
	if err != nil {
		rt.logger.Error("movie not found", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, m, http.StatusOK)
}

// @Summary Create a new movie
// @Description Create a new movie
// @ID create-movie
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Accept json
// @Produce json
// @Param movie body createPayload true "Movie object to be created"
// @Success 201 {object} response{response=domain.Movie}
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /create [post]
func (rt *router) createHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level != "anonymous" {
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

	var p createPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	m := domain.NewMovie(p.Title, p.OriginalTitle, p.Poster, p.Genres)

	vm, err := domain.NewValidatedMovie(m)
	if err != nil {
		rt.logger.Error("invalid movie data", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	m, err = rt.repository.Create(ctx, vm)
	if err != nil {
		rt.logger.Error("failed to create movie", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, m, http.StatusCreated)
}
