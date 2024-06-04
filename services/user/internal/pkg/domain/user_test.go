package domain

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	username := "user123"
	password := "password"
	name := "John Doe"
	picture := "http://example.com/picture.jpg"
	level := "user"

	user := NewUser(username, password, name, picture)

	assert.Equal(t, username, user.ID)
	assert.Equal(t, username, user.Username)
	assert.Equal(t, password, user.Password)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, picture, user.Picture)
	assert.Equal(t, level, user.Level)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
}

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name          string
		user          *User
		expectedError error
	}{
		{
			name: "valid user without picture",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "valid user with picture",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: nil,
		},
		{
			name: "invalid user with invalid picture URL",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "invalid_url",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("provided picture image source is invalid or too slow to load"),
		},
		{
			name: "missing ID",
			user: &User{
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("id is required"),
		},
		{
			name: "missing Username",
			user: &User{
				ID:        "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("username is required"),
		},
		{
			name: "ID and Username do not match",
			user: &User{
				ID:        "user123",
				Username:  "user456",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("id and username must be the same"),
		},
		{
			name: "missing Password",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("password is required"),
		},
		{
			name: "missing Name",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("name is required"),
		},
		{
			name: "CreatedAt after UpdatedAt",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     userLevel,
				CreatedAt: time.Now().Add(1 * time.Hour),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("created_at must be before updated_at"),
		},
		{
			name: "missing Level",
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("invalid user level"),
		},
		{
			name: "admin Level", // not used yet, thus invalid
			user: &User{
				ID:        "user123",
				Username:  "user123",
				Password:  "password",
				Name:      "John Doe",
				Picture:   "http://example.com/picture.jpg",
				Level:     adminLevel,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedError: errors.New("invalid user level"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.validate(false)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
