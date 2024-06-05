package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	logger     *log.Logger
}

type RefreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type ValidationPayload struct {
	AccessToken string `json:"accessToken"`
}

type tokensResponse struct {
	Response Tokens `json:"response"`
}

type claimsResponse struct {
	Response Claims `json:"response"`
}

type Tokens struct {
	RefreshToken           string `json:"refreshToken,omitempty"`
	AccessToken            string `json:"accessToken"`
	AccessTokenExpiration  int64  `json:"accessTokenExpiration"`
	RefreshTokenExpiration int64  `json:"refreshTokenExpiration,omitempty"`
}

type Claims struct {
	Name  string `json:"name,omitempty"`
	Level string `json:"level"`
	jwt.RegisteredClaims
}

func NewClient(baseURL string, timeout time.Duration, logger *log.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		logger: logger,
	}
}

func (c *Client) GenerateAnonymousTokens(ctx context.Context) (*Tokens, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/anonymous", c.baseURL), nil)
	if err != nil {
		c.logger.Error("failed to create request", log.Error(err))
		return nil, err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		c.logger.Error("error from authentication service", log.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var result tokensResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.logger.Error("failed to parse response", log.Error(err))
		return nil, err
	}

	return &result.Response, nil
}

func (c *Client) Refresh(ctx context.Context, payload RefreshPayload) (*Tokens, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		c.logger.Error("failed to parse request body", log.Error(err))
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/refresh", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error("failed to create request", log.Error(err))
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		c.logger.Error("error from authentication service", log.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var result tokensResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.logger.Error("failed to parse response", log.Error(err))
		return nil, err
	}

	return &result.Response, nil
}

func (c *Client) ValidateAccessToken(ctx context.Context, payload ValidationPayload) (*Claims, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		c.logger.Error("failed to parse request body", log.Error(err))
		return nil, err
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/validate", c.baseURL), bytes.NewBuffer(data))
	if err != nil {
		c.logger.Error("failed to create request", log.Error(err))
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(r)
	if err != nil {
		c.logger.Error("error from authentication service", log.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	var result claimsResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		c.logger.Error("failed to parse response", log.Error(err))
		return nil, err
	}

	return &result.Response, nil
}
