package domain

import "context"

// Repository is the interface for the domain's repository (e.g. some database).
type Repository interface {
	// Create receives a validated input and creates a new User.
	Create(ctx context.Context, user *ValidatedUser) (*User, error)
	// FindById retrieves an User by a given unique ID.
	FindByID(ctx context.Context, id string) (*User, error)
	// FindByUsername retrieves an User by a given unique username.
	FindByUsername(ctx context.Context, username string) (*User, error)
	// Close disconnects the database connection pool.
	Close(ctx context.Context) error
}
