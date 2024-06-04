package domain

import (
	"crypto/rsa"
)

type Authenticator interface {
	// GenerateAnonymousTokens generate authentication tokens for anonymous user.
	GenerateAnonymousTokens(userID string, flow FlowType) (*Tokens, error)
	// GenerateUserTokens generate authentication tokens for logged-in user.
	GenerateUserTokens(username, password string, flow FlowType) (*Tokens, error)
	// Revoke revokes an access token registered for a given refresh token.
	Revoke(accessToken string) error
	// Refresh refreshes an user authentication tokens.
	Refresh(refreshToken string) (*Tokens, error)
	// ValidateAccessToken checks if logged-in user authentication token is valid.
	ValidateAccessToken(accessToken string) (*Claims, error)
	// JWTKey returns the authenticator JWT Keys.
	JWTKey() *rsa.PrivateKey
}

// Tokens represents the authentication token json response.
type Tokens struct {
	RefreshToken           string `json:"refreshToken,omitempty"`
	AccessToken            string `json:"accessToken"`
	AccessTokenExpiration  int64  `json:"accessTokenExpiration"`
	RefreshTokenExpiration int64  `json:"refreshTokenExpiration,omitempty"`
}

// FlowType represents login flow type.
type FlowType string

// Valid FlowType values as constants.
const (
	WebsiteSessionFlow FlowType = "websiteSession"
	RememberMeFlow     FlowType = "rememberMe"
)
