package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMovie(t *testing.T) {
	title := "Movie Title"
	originalTitle := "Original Movie Title"
	poster := "https://example.com/poster.jpg"
	genres := []string{"Action", "Adventure"}

	movie := NewMovie(title, originalTitle, poster, genres)

	assert.NotNil(t, movie)
	assert.NotEmpty(t, movie.ID)
	assert.Equal(t, title, movie.Title)
	assert.Equal(t, originalTitle, movie.OriginalTitle)
	assert.Equal(t, poster, movie.Poster)
	assert.Equal(t, genres, movie.Genres)
	assert.True(t, movie.CreatedAt.Before(time.Now()))
	assert.True(t, movie.UpdatedAt.Before(time.Now()))
}

func TestValidateMovie(t *testing.T) {
	tests := []struct {
		name  string
		movie *Movie
		err   error
	}{
		{
			name: "ValidMovie",
			movie: &Movie{
				ID:            "123",
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Poster:        "https://example.com/poster.jpg",
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now().Add(-time.Hour),
				UpdatedAt:     time.Now(),
			},
			err: nil,
		},
		{
			name: "InvalidMovie_NoID",
			movie: &Movie{
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Poster:        "https://example.com/poster.jpg",
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now().Add(-time.Hour),
				UpdatedAt:     time.Now(),
			},
			err: errors.New("id is required"),
		},
		{
			name: "InvalidMovie_NoTitle",
			movie: &Movie{
				ID:            "123",
				OriginalTitle: "Original Movie Title",
				Poster:        "https://example.com/poster.jpg",
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now().Add(-time.Hour),
				UpdatedAt:     time.Now(),
			},
			err: errors.New("title is required"),
		},
		{
			name: "InvalidMovie_NoOriginalTitle",
			movie: &Movie{
				ID:        "123",
				Title:     "Movie Title",
				Poster:    "https://example.com/poster.jpg",
				Genres:    []string{"Action", "Adventure"},
				CreatedAt: time.Now().Add(-time.Hour),
				UpdatedAt: time.Now(),
			},
			err: errors.New("originalTitle is required"),
		},
		{
			name: "InvalidMovie_NoPoster",
			movie: &Movie{
				ID:            "123",
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now().Add(-time.Hour),
				UpdatedAt:     time.Now(),
			},
			err: errors.New("poster is required"),
		},
		{
			name: "InvalidMovie_InvalidPosterSource",
			movie: &Movie{
				ID:            "123",
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Poster:        "invalidurl", // Invalid poster source
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now().Add(-time.Hour),
				UpdatedAt:     time.Now(),
			},
			err: errors.New("provided poster image source is invalid or too slow to load"),
		},
		{
			name: "InvalidMovie_NoGenres",
			movie: &Movie{
				ID:            "123",
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Poster:        "https://example.com/poster.jpg",
				// No genres provided
				CreatedAt: time.Now().Add(-time.Hour),
				UpdatedAt: time.Now(),
			},
			err: errors.New("at least one genre is required"),
		},
		{
			name: "InvalidMovie_CreatedAtAfterUpdatedAt",
			movie: &Movie{
				ID:            "123",
				Title:         "Movie Title",
				OriginalTitle: "Original Movie Title",
				Poster:        "https://example.com/poster.jpg",
				Genres:        []string{"Action", "Adventure"},
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now().Add(-time.Hour), // UpdatedAt before CreatedAt
			},
			err: errors.New("created_at must be before updated_at"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.movie.validate(false)
			assert.Equal(t, tc.err, err)
		})
	}
}
