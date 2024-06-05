package client

import (
	"context"
	"net/http"
	"strings"

	libCtx "github.com/victorspringer/backend-coding-challenge/lib/context"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

// Middleware handles authentication / authorization for a request.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if authToken == "" {
			ck, err := r.Cookie("MRSAccessToken")
			if err == nil && ck.Value != "" {
				authToken = ck.Value
			} else {
				ck, err := r.Cookie("MRSRefreshToken")
				if err == nil && ck.Value != "" {
					tokens, err := c.Refresh(ctx, RefreshPayload{ck.Value})
					if err != nil {
						c.logger.Error("failed to refresh token", log.Error(err))
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					authToken = tokens.AccessToken

					isSecure := r.Header.Get("X-Forwarded-Proto") == "https"
					http.SetCookie(w, &http.Cookie{
						Name:     "MRSAccessToken",
						Value:    tokens.AccessToken,
						HttpOnly: true,
						Secure:   isSecure,
						Path:     "/",
						MaxAge:   int(tokens.AccessTokenExpiration),
						SameSite: http.SameSiteLaxMode,
					})
					http.SetCookie(w, &http.Cookie{
						Name:     "MRSRefreshToken",
						Value:    tokens.RefreshToken,
						HttpOnly: true,
						Secure:   isSecure,
						Path:     "/",
						MaxAge:   int(tokens.RefreshTokenExpiration),
						SameSite: http.SameSiteLaxMode,
					})
				} else {
					tokens, err := c.GenerateAnonymousTokens(r.Context())
					if err != nil {
						c.logger.Error("failed to generate anonymous token", log.Error(err))
						w.WriteHeader(http.StatusUnauthorized)
						return
					}

					authToken = tokens.AccessToken

					isSecure := r.Header.Get("X-Forwarded-Proto") == "https"
					http.SetCookie(w, &http.Cookie{
						Name:     "MRSAccessToken",
						Value:    tokens.AccessToken,
						HttpOnly: true,
						Secure:   isSecure,
						Path:     "/",
						MaxAge:   int(tokens.AccessTokenExpiration),
						SameSite: http.SameSiteLaxMode,
					})
					http.SetCookie(w, &http.Cookie{
						Name:   "MRSRefreshToken",
						MaxAge: -1,
					})
				}
			}
		}

		if authToken != "" {
			claims, err := c.ValidateAccessToken(ctx, ValidationPayload{authToken})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, libCtx.CTX_USER_USERNAME, claims.Subject)
			ctx = context.WithValue(ctx, libCtx.CTX_USER_LEVEL, claims.Level)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
