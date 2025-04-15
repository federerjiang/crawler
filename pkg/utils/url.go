package utils

import (
	"net/url"
	"strings"
)

// IsValidURL checks if a string is a valid URL
func IsValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// JoinURL joins a base URL with a path
func JoinURL(baseURL, path string) (string, error) {
	// Parse base URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// If path is already a full URL, return it
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path, nil
	}

	// Join base with path
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	return base.ResolveReference(u).String(), nil
}

// GetDomain extracts the domain from a URL
func GetDomain(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	return u.Hostname(), nil
}

// NormalizeURL normalizes a URL by removing trailing slashes, fragments, and query parameters
func NormalizeURL(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Remove fragment
	u.Fragment = ""

	// Remove query parameters
	u.RawQuery = ""

	// Ensure scheme
	if u.Scheme == "" {
		u.Scheme = "http"
	}

	// Remove trailing slash
	path := u.Path
	if path != "/" && strings.HasSuffix(path, "/") {
		u.Path = strings.TrimSuffix(path, "/")
	}

	return u.String(), nil
}
