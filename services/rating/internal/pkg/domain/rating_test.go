package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewRating(t *testing.T) {
	userID := "user-123"
	movieID := "movie-456"
	value := float32(4.5)

	rating := NewRating(userID, movieID, value)

	assert.NotEmpty(t, rating.ID)
	assert.Equal(t, userID, rating.UserID)
	assert.Equal(t, movieID, rating.MovieID)
	assert.Equal(t, value, rating.Value)
	assert.WithinDuration(t, time.Now(), rating.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), rating.UpdatedAt, time.Second)
}

func TestRating_Validate(t *testing.T) {
	tests := []struct {
		name      string
		rating    *Rating
		wantError bool
	}{
		{
			name: "valid rating",
			rating: &Rating{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				MovieID:   "movie-456",
				Value:     4.5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: false,
		},
		{
			name: "missing ID",
			rating: &Rating{
				UserID:    "user-123",
				MovieID:   "movie-456",
				Value:     4.5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
		{
			name: "missing UserID",
			rating: &Rating{
				ID:        uuid.New().String(),
				MovieID:   "movie-456",
				Value:     4.5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
		{
			name: "missing MovieID",
			rating: &Rating{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				Value:     4.5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
		{
			name: "value out of range (too high)",
			rating: &Rating{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				MovieID:   "movie-456",
				Value:     5.5,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
		{
			name: "value out of range (too low)",
			rating: &Rating{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				MovieID:   "movie-456",
				Value:     0.4,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
		{
			name: "created_at after updated_at",
			rating: &Rating{
				ID:        uuid.New().String(),
				UserID:    "user-123",
				MovieID:   "movie-456",
				Value:     4.5,
				CreatedAt: time.Now().Add(1 * time.Hour),
				UpdatedAt: time.Now(),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rating.validate()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
