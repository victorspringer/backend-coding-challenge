package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/log"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

func (rt *router) findHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := chi.URLParam(r, "username")

	u, err := rt.repository.FindByUsername(ctx, username)
	if err != nil {
		rt.logger.Error("user not found", log.Error(err), log.String("requestId", getRequestID(ctx)))
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

	u := domain.NewUser(p.Username, p.MD5Password, p.Name, p.Picture)

	vu, err := domain.NewValidatedUser(u)
	if err != nil {
		rt.logger.Error("invalid user data", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	u, err = rt.repository.Create(ctx, vu)
	if err != nil {
		rt.logger.Error("failed to create user", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, u, http.StatusCreated)
}
