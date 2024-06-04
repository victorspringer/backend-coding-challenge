package router

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
	"github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/domain"
)

func (rt *router) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

// @Summary Login anonymously
// @Description Generates anonymous tokens for a new user
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} response{response=domain.Tokens}
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /anonymous [post]
func (rt *router) loginAnonymous(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := uuid.New().String()
	tokens, err := rt.authenticator.GenerateAnonymousTokens(userID, domain.WebsiteSessionFlow)
	if err != nil {
		rt.logger.Info("failed to generate anonymous tokens", log.Error(err), log.String("requestId", getRequestID(ctx)))
		if err == domain.ErrUnauthorized {
			rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, tokens, http.StatusOK)
}

// @Summary Login
// @Description Generates tokens for an existing user
// @Tags authentication
// @Accept json
// @Produce json
// @Param loginPayload body loginPayload true "Login payload"
// @Success 200 {object} response{response=domain.Tokens}
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /login [post]
func (rt *router) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body := &loginPayload{}
	err = json.Unmarshal(b, body)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tokens, err := rt.authenticator.GenerateUserTokens(body.Username, body.MD5Password, body.Flow)
	if err != nil {
		rt.logger.Error("failed to generate user tokens", log.Error(err), log.String("requestId", getRequestID(ctx)))
		if err == domain.ErrUnauthorized {
			rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, tokens, http.StatusOK)
}

// @Summary Refresh tokens
// @Description Refreshes existing tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param refreshPayload body refreshPayload true "Refresh payload"
// @Success 200 {object} response{response=domain.Tokens}
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /refresh [post]
func (rt *router) refresh(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body := &refreshPayload{}
	err = json.Unmarshal(b, body)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	tokens, err := rt.authenticator.Refresh(body.RefreshToken)
	if err != nil {
		rt.logger.Error("failed to refresh tokens", log.Error(err), log.String("requestId", getRequestID(ctx)))
		if err == domain.ErrUnauthorized {
			rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, tokens, http.StatusOK)
}

// @Summary Logout
// @Description Revokes the access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header with Bearer token"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /logout [post]
func (rt *router) logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		rt.logger.Info("Authorization header is not set")
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	accessToken = strings.Replace(accessToken, "Bearer ", "", 1)
	err := rt.authenticator.Revoke(accessToken)
	if err != nil {
		rt.logger.Error("failed to revoke access token", log.Error(err), log.String("requestId", getRequestID(ctx)))
		if err == domain.ErrUnauthorized {
			rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, http.StatusText(http.StatusOK), http.StatusOK)
}

// @Summary Validate access token
// @Description Validates the provided access token
// @Tags authentication
// @Accept json
// @Produce json
// @Param validationPayload body validationPayload true "Validation payload"
// @Success 200 {object} response{response=domain.Claims}
// @Failure 400 {object} response
// @Failure 401 {object} response
// @Failure 500 {object} response
// @Router /validate [post]
func (rt *router) validateAccessToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		rt.logger.Error("failed to read request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body := &validationPayload{}
	err = json.Unmarshal(b, body)
	if err != nil {
		rt.logger.Error("failed to parse request body", log.Error(err), log.String("requestId", getRequestID(ctx)))
		rt.respond(w, r, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	claims, err := rt.authenticator.ValidateAccessToken(body.AccessToken)
	if err != nil {
		rt.logger.Debug("invalid access token", log.Error(err), log.String("requestId", getRequestID(ctx)))
		if err == domain.ErrUnauthorized {
			rt.respond(w, r, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		rt.respond(w, r, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	rt.respond(w, r, claims, http.StatusOK)
}

// @Summary JSON Web Key Set
// @Description Returns the JSON Web Key Set
// @Tags authentication
// @Accept json
// @Produce json
// @Success 200 {object} response{response=jwks}
// @Router /.well-known/jwks.json [get]
func (rt *router) jwks(w http.ResponseWriter, r *http.Request) {
	bigExp := big.NewInt(int64(rt.authenticator.JWTKey().PublicKey.E))
	exponent := base64.RawURLEncoding.EncodeToString(bigExp.Bytes())
	modulus := base64.RawURLEncoding.EncodeToString(rt.authenticator.JWTKey().PublicKey.N.Bytes())
	jwks := &jwks{
		Keys: []key{
			{
				Algorithm: "RSA256",
				KeyType:   "RSA",
				Use:       "sig",
				N:         modulus,
				E:         exponent,
			},
		},
	}
	rt.respond(w, r, jwks, http.StatusOK)
}
