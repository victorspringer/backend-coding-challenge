package database

import (
	"fmt"

	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
)

// TODO: implement database

type database struct{}

// New returns a new instance of database.
func New() (domain.Repository, error) {
	return &database{}, nil
}

// Create implements domain.Repository interface's Create method.
func (db *database) Create(user *domain.ValidatedUser) (*domain.User, error) {
	if user.IsValid() {
		return &user.User, nil
	}

	return nil, fmt.Errorf(
		"invalid user (id: %s, username: %s, md5 password: %s, name: %s, picture: %s)",
		user.ID, user.Username, user.Password, user.Name, user.Picture,
	)
}

// FindByID implements domain.Repository interface's FindByID method.
func (db *database) FindByID(id string) (*domain.User, error) {
	return nil, nil
}

// FindByUsername implements domain.Repository interface's FindByUsername method.
func (db *database) FindByUsername(username string) (*domain.User, error) {
	return nil, nil
}
