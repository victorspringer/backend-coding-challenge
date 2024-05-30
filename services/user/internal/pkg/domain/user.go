package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// User entity.
type User struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  string
	Name      string
	Picture   string
}

// NewUser returns an instance of the User entity.
func NewUser(username, password, name, picture string) *User {
	return &User{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Username:  username,
		Password:  password,
		Name:      name,
		Picture:   picture,
	}
}

func (p *User) validate() error {
	if p.ID == "" {
		return errors.New("id is required")
	}
	if p.Username == "" {
		return errors.New("username is required")
	}
	if p.Password == "" {
		return errors.New("password is required")
	}
	if p.Name == "" {
		return errors.New("name is required")
	}
	// TODO: validate profile picture
	if p.CreatedAt.After(p.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}
