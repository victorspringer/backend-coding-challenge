package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type Tokens struct {
	RefreshToken           string `json:"refreshToken,omitempty"`
	AccessToken            string `json:"accessToken"`
	AccessTokenExpiration  int64  `json:"accessTokenExpiration"`
	RefreshTokenExpiration int64  `json:"refreshTokenExpiration,omitempty"`
}

type RefreshPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type ValidationPayload struct {
	AccessToken string `json:"accessToken"`
}

type Claims struct {
	Name  string `json:"name,omitempty"`
	Level string `json:"level"`
	jwt.RegisteredClaims
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) GenerateAnonymousTokens() (*Tokens, error) {
	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/anonymous", c.BaseURL), "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Tokens
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) Refresh(payload RefreshPayload) (*Tokens, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/refresh", c.BaseURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Tokens
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) ValidateAccessToken(payload ValidationPayload) (*Claims, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(fmt.Sprintf("%s/validate", c.BaseURL), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result Claims
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
