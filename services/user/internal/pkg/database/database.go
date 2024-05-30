package database

import (
	"github.com/victorspringer/backend-coding-challenge/services/user/internal/pkg/domain"
)

// TODO: implement database
// TODO: add comments

type database struct{}

func New() (domain.Repository, error) {
	return &database{}, nil
}

func (db *database) Create(user *domain.ValidatedUser) (*domain.User, error) {
	return &user.User, nil
}

func (db *database) FindByID(id string) (*domain.User, error) {
	return nil, nil
}

func (db *database) FindByUsername(username string) (*domain.User, error) {
	return nil, nil
}
