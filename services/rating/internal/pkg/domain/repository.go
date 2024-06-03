package domain

import "context"

// Repository is the interface for the domain's repository (e.g. some database).
type Repository interface {
	// Upsert receives a validated input and upserts a Rating.
	Upsert(ctx context.Context, rating *ValidatedRating) (*Rating, error)
	// FindByUserID retrieves a list of Rating by a given user ID.
	FindByUserID(ctx context.Context, userID string) ([]*Rating, error)
	// FindByMovieID retrieves a list of Rating by a given movie ID.
	FindByMovieID(ctx context.Context, movieID string) ([]*Rating, error)
	// Close disconnects the database connection pool.
	Close(ctx context.Context) error
}
