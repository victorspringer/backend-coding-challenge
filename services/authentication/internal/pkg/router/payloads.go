package router

import "github.com/victorspringer/backend-coding-challenge/services/authentication/internal/pkg/domain"

type loginPayload struct {
	Username    string          `json:"username"`
	MD5Password string          `json:"md5Password"`
	Flow        domain.FlowType `json:"flow"`
}

type refreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type jwks struct {
	Keys []key `json:"keys"`
}

type key struct {
	Algorithm string `json:"alg"`
	KeyType   string `json:"kty"`
	Use       string `json:"use"`
	N         string `json:"n"`
	E         string `json:"e"`
}

type validationPayload struct {
	AccessToken string `json:"accessToken"`
}
