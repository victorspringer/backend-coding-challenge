package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/image"
)

// User entity.
type User struct {
	ID        string    `bson:"id"`
	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
	Username  string    `bson:"username"`
	Password  string    `bson:"password"`
	Name      string    `bson:"name"`
	Picture   string    `bson:"picture"`
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

func (u *User) validate() error {
	if u.ID == "" {
		return errors.New("id is required")
	}
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Picture != "" && !image.IsValidSource(u.Picture) {
		return errors.New("provided image source is invalid or too slow to load")
	}
	if u.CreatedAt.After(u.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}
