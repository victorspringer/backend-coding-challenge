package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/victorspringer/backend-coding-challenge/services/movie/internal/pkg/image"
)

// Movie entity.
type Movie struct {
	ID            string    `json:"id" bson:"id"`
	Title         string    `json:"title" bson:"title"`
	OriginalTitle string    `json:"originalTitle" bson:"originalTitle"`
	Overview      string    `json:"overview" bson:"overview"`
	Poster        string    `json:"poster" bson:"poster"`
	Genres        []string  `json:"genres" bson:"genres"`
	Keywords      []string  `json:"keywords" bson:"keywords"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

// NewMovie returns an instance of the Movie entity.
func NewMovie(title, originalTitle, overview, poster string, genres, keywords []string) *Movie {
	return &Movie{
		ID:            uuid.New().String(),
		Title:         title,
		OriginalTitle: originalTitle,
		Overview:      overview,
		Poster:        poster,
		Genres:        genres,
		Keywords:      keywords,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (m *Movie) validate() error {
	if m.ID == "" {
		return errors.New("id is required")
	}
	if m.Title == "" {
		return errors.New("title is required")
	}
	if m.OriginalTitle == "" {
		return errors.New("originalTitle is required")
	}
	if m.Overview == "" {
		return errors.New("overview is required")
	}
	if m.Poster == "" {
		return errors.New("poster is required")
	}
	if !image.IsValidSource(m.Poster) {
		return errors.New("provided poster image source is invalid or too slow to load")
	}
	if len(m.Genres) == 0 {
		return errors.New("at least one genre is required")
	}
	if len(m.Keywords) == 0 {
		return errors.New("at least one keyword is required")
	}
	if m.CreatedAt.After(m.UpdatedAt) {
		return errors.New("created_at must be before updated_at")
	}

	return nil
}
