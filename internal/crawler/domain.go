package crawler

import (
	"net/url"
	"strings"
)

// ExtractDomain extracts the domain from a URL
func ExtractDomain(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	return parsedURL.Hostname(), nil
}

// IsSameDomain checks if a URL belongs to the same domain as the start URL
func IsSameDomain(startURL, urlToCheck string) (bool, error) {
	startDomain, err := ExtractDomain(startURL)
	if err != nil {
		return false, err
	}

	urlDomain, err := ExtractDomain(urlToCheck)
	if err != nil {
		return false, err
	}

	// Handle subdomains - if we want to consider them as different domains
	// For now, we'll consider subdomains as different domains
	return urlDomain == startDomain, nil
}

// NormalizeURL normalizes a URL by removing trailing slashes and fragments
func NormalizeURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	// Remove fragments
	parsedURL.Fragment = ""

	// Ensure scheme is present
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}

	// Remove trailing slash
	if parsedURL.Path != "/" && strings.HasSuffix(parsedURL.Path, "/") {
		parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")
	}

	return parsedURL.String(), nil
}
