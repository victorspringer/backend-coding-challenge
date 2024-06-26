package router

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

// @Summary Get user by username
// @Description Get user information by username
// @ID get-user-by-username
// @Param username path string true "Username of the user"
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token"
// @Produce json
// @Success 200 {object} response{response=domain.User}
// @Failure 401 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /{username} [get]
func (rt *router) findHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if level := context.GetUserLevel(ctx); level == "anonymous" {
		rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	username := chi.URLParam(r, "username")

	u, err := rt.repository.FindByID(ctx, username)
	if err != nil {
		rt.logger.Error("user not found", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusNotFound)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
}

// @Summary Create a new user
// @Description Create a new user
// @ID create-user
// @Accept json
// @Produce json
// @Param user body createPayload true "User object to be created"
// @Success 201 {object} response{response=domain.User}
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /create [post]
func (rt *router) createHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

	u := domain.NewUser(p.Username, p.MD5Password, p.Name, p.Picture)

	vu, err := domain.NewValidatedUser(u)
	if err != nil {
		rt.logger.Error("invalid user data", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	u, err = rt.repository.Create(ctx, vu)
	if err != nil {
		rt.logger.Error("failed to create user", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, u, http.StatusCreated)
}

// @Summary Get user by credentials (username and password)
// @Description Get user information by username and password
// @ID get-user-by-credentials
// @Param user body credentialsPayload true "User object to be found"
// @Produce json
// @Success 200 {object} response{response=domain.User}
// @Failure 400 {object} response
// @Failure 404 {object} response
// @Failure 500 {object} response
// @Router /credentials [post]
func (rt *router) findByCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	var p credentialsPayload
	err = json.Unmarshal(b, &p)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := rt.repository.FindByCredentials(ctx, p.Username, p.MD5Password)
	if err != nil {
		if err.Error() == "invalid credentials" {
			rt.logger.Info("user not found", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
			rt.respond(w, r, err.Error(), http.StatusNotFound)
			return
		}
		rt.logger.Error("internal error", log.Error(err), log.String("requestId", context.GetRequestID(ctx)))
		rt.respond(w, r, err.Error(), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, u, http.StatusOK)
}
