package domain

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// User entity.
type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Level    string `json:"level"`
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
}

type payload struct {
	Username    string `json:"username"`
	MD5Password string `json:"md5Password"`
}

func NewUserServiceClient(timeout time.Duration, endpoint string) UserServiceClient {
	return &userServiceClient{
		&http.Client{
			Timeout: timeout,
		},
		endpoint,
	}
}

func (c *userServiceClient) CheckCredentials(username, md5Password string) (*User, error) {
	p := payload{username, md5Password}
	b, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/credentials", c.userServiceEndpoint),
		bytes.NewReader(b),
	)
	if err != nil {
		return nil, err
	}

	r.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errUserNotFound
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	u := &User{}
	if err = json.Unmarshal(body, u); err != nil {
		return nil, err
	}

	return u, nil
}
