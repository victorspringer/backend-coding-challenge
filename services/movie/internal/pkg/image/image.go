package image

import (
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

// IsValidSource checks if the given image URL is a valid image file.
// If the downloading is too slow (more than 1 second) it also invalidates the image.
func IsValidSource(imgURL string) bool {
	return isValidImageURL(imgURL) && isValidImageContent(imgURL)
}

func isValidImageURL(imgURL string) bool {
	u, err := url.Parse(imgURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	ext := strings.ToLower(path.Ext(u.Path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	}

	return false
}

func isValidImageContent(imgURL string) bool {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(imgURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	contentType := resp.Header.Get("Content-Type")
	return strings.HasPrefix(contentType, "image/")
}
