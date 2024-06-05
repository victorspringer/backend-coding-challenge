package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/victorspringer/backend-coding-challenge/lib/log"
)

// User entity.
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Level    Level  `json:"level"`
}

type UserServiceClient interface {
	CheckCredentials(username, md5Password string) (*User, error)
}

var (
	errBadRequest   = errors.New("bad request to user service")
	errUserNotFound = errors.New("user not found in user service")
)

type userServiceClient struct {
	httpClient          *http.Client
	userServiceEndpoint string
	logger              *log.Logger
}

type payload struct {
	Username    string `json:"username"`
	MD5Password string `json:"md5Password"`
}

type userResponse struct {
	Response User `json:"response"`
}

func NewUserServiceClient(timeout time.Duration, endpoint string, logger *log.Logger) UserServiceClient {
	return &userServiceClient{
		&http.Client{
			Timeout: timeout,
		},
		endpoint,
		logger,
	}
}

func (c *userServiceClient) CheckCredentials(username, md5Password string) (*User, error) {
	p := payload{username, md5Password}
	b, err := json.Marshal(p)
	if err != nil {
		c.logger.Error("failed to parse request body", log.Error(err))
		return nil, err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/credentials", c.userServiceEndpoint),
		bytes.NewReader(b),
	)
	if err != nil {
		c.logger.Error("failed to create request", log.Error(err))
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(r)
	if err != nil {
		c.logger.Error("error from user service", log.Error(err))
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errUserNotFound
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.logger.Error("failed to read response body", log.Error(err))
		return nil, err
	}

	u := &userResponse{}
	if err = json.Unmarshal(body, u); err != nil {
		c.logger.Error("failed to parse response body", log.Error(err))
		return nil, err
	}

	return &u.Response, nil
}
