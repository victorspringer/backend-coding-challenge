package image

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsValidSource(t *testing.T) {
	tests := []struct {
		name            string
		imgURL          string
		validateContent bool
		expected        bool
	}{
		{"Valid image URL without validation", "http://example.com/image.png", false, true},
		{"Valid image URL with validation", "http://example.com/image.png", true, true},
		{"Invalid URL", "invalid-url", false, false},
		{"Non-image URL", "http://example.com/file.txt", false, false},
		{"Non-image URL with validation", "http://example.com/file.txt", true, false},
		{"Valid image URL with slow response", "http://slow.example.com/image.png", true, false},
		{"Valid image URL with invalid content type", "http://example.com/image.txt", true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.imgURL == "http://slow.example.com/image.png" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					http.ServeFile(w, r, "./test/files/image.png")
				}))
				defer server.Close()
				tt.imgURL = server.URL + "/image.png"
			} else if tt.imgURL == "http://example.com/image.png" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/png")
					http.ServeFile(w, r, "./test/files/image.png")
				}))
				defer server.Close()
				tt.imgURL = server.URL + "/image.png"
			} else if tt.imgURL == "http://example.com/file.txt" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					http.ServeFile(w, r, "./test/files/file.txt")
				}))
				defer server.Close()
				tt.imgURL = server.URL + "/file.txt"
			}

			got := IsValidSource(tt.imgURL, tt.validateContent)
			if got != tt.expected {
				t.Errorf("IsValidSource(%q, %v) = %v; want %v", tt.imgURL, tt.validateContent, got, tt.expected)
			}
		})
	}
}

func TestIsValidImageURL(t *testing.T) {
	tests := []struct {
		name     string
		imgURL   string
		expected bool
	}{
		{"Valid image URL", "http://example.com/image.png", true},
		{"Invalid URL", "invalid-url", false},
		{"Non-image URL", "http://example.com/file.txt", false},
		{"Unsupported image extension", "http://example.com/image.bmp", true},
		{"Missing scheme", "example.com/image.png", false},
		{"Missing host", "http:///image.png", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidImageURL(tt.imgURL)
			if got != tt.expected {
				t.Errorf("isValidImageURL(%q) = %v; want %v", tt.imgURL, got, tt.expected)
			}
		})
	}
}

func TestIsValidImageContent(t *testing.T) {
	tests := []struct {
		name     string
		imgURL   string
		expected bool
	}{
		{"Valid image content", "http://example.com/image.png", true},
		{"Slow response", "http://slow.example.com/image.png", false},
		{"Invalid content type", "http://example.com/image.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.imgURL == "http://slow.example.com/image.png" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					time.Sleep(2 * time.Second)
					http.ServeFile(w, r, "./test/files/image.png")
				}))
				defer server.Close()
				tt.imgURL = server.URL
			} else if tt.imgURL == "http://example.com/image.png" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "image/png")
					http.ServeFile(w, r, "./test/files/image.png")
				}))
				defer server.Close()
				tt.imgURL = server.URL
			} else if tt.imgURL == "http://example.com/image.txt" {
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					http.ServeFile(w, r, "./test/files/file.txt")
				}))
				defer server.Close()
				tt.imgURL = server.URL
			}

			got := isValidImageContent(tt.imgURL)
			if got != tt.expected {
				t.Errorf("isValidImageContent(%q) = %v; want %v", tt.imgURL, got, tt.expected)
			}
		})
	}
}
