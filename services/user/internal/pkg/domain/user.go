package domain

import (
	"errors"
	"time"

	"github.com/victorspringer/backend-coding-challenge/lib/image"
)

// User entity.
type User struct {
	ID        string    `json:"id" bson:"id"`
	Username  string    `json:"username" bson:"username"`
	Password  string    `json:"password" bson:"password"`
	Name      string    `json:"name" bson:"name"`
	Picture   string    `json:"picture" bson:"picture"`
	Level     string    `json:"level" bson:"level"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

const (
	adminLevel = "admin" // not used yet
	userLevel  = "user"
)

// NewUser returns an instance of the User entity.
func NewUser(username, password, name, picture string) *User {
	return &User{
		ID:        username, // id and username are the same, as the username is unique
		Username:  username,
		Password:  password,
		Name:      name,
		Picture:   picture,
		Level:     userLevel,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (u *User) validate(validateImageContent ...bool) error {
	vc := true
	if len(validateImageContent) > 0 {
		vc = validateImageContent[0]
	}

	if u.ID == "" {
		return errors.New("id is required")
	}
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.ID != u.Username {
		return errors.New("id and username must be the same")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Level != userLevel {
		return errors.New("invalid user level")
	}
	if u.CreatedAt.After(u.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}
	if u.Picture != "" && !image.IsValidSource(u.Picture, vc) {
		return errors.New("provided picture image source is invalid or too slow to load")
	}

	return nil
}
