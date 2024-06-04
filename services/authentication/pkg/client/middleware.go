package client

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ctxKeyUserLevel contextKey = "userLevel"
)

// Middleware handles authentication / authorization for a request.
func (c *Client) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
		if authToken == "" {
			ck, err := r.Cookie("MRSAccessToken")
			if err == nil && ck.Value != "" {
				authToken = ck.Value
			} else {
				ck, err := r.Cookie("MRSRefreshToken")
				if err == nil && ck.Value != "" {
					tokens, err := c.Refresh(RefreshPayload{ck.Value})
					if err != nil {
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
					tokens, err := c.GenerateAnonymousTokens()
					if err != nil {
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
				}
			}
		}

		ctx := r.Context()

		if authToken != "" {
			claims, err := c.ValidateAccessToken(ValidationPayload{authToken})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, ctxKeyUserLevel, claims.Level)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
