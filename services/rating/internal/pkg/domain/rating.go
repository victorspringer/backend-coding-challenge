package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Rating entity.
type Rating struct {
	ID        string    `json:"id" bson:"id"`
	UserID    string    `json:"userId" bson:"userId"`
	MovieID   string    `json:"movieId" bson:"movieId"`
	Value     float32   `json:"value" bson:"value"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// NewRating returns an instance of the Rating entity.
func NewRating(userID, movieID string, value float32) *Rating {
	return &Rating{
		ID:        uuid.New().String(),
		UserID:    userID,
		MovieID:   movieID,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *Rating) validate() error {
	if r.ID == "" {
		return errors.New("id is required")
	}
	if r.UserID == "" {
		return errors.New("userId is required")
	}
	if r.MovieID == "" {
		return errors.New("movieId is required")
	}
	if r.Value > 5 || r.Value < 0.5 {
		return errors.New("rating value must be from 0.5 to 5")
	}
	if r.CreatedAt.After(r.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}
