package domain

import "context"

// Repository is the interface for the domain's repository (e.g. some database).
type Repository interface {
	// Create receives a validated input and creates a new Movie.
	Create(ctx context.Context, movie *ValidatedMovie) (*Movie, error)
	// FindById retrieves an Movie by a given unique ID.
	FindByID(ctx context.Context, id string) (*Movie, error)
	// Close disconnects the database connection pool.
	Close(ctx context.Context) error
}
