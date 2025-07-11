package utils

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// PingURLOptions holds configuration for URL ping
type PingURLOptions struct {
	Timeout      time.Duration
	MaxRedirects int
}

// DefaultPingOptions returns default options for URL ping
func DefaultPingOptions() PingURLOptions {
	return PingURLOptions{
		Timeout:      10 * time.Second,
		MaxRedirects: 5,
	}
}

// PingURLResult holds the result of a URL ping
type PingURLResult struct {
	Available    bool
	StatusCode   int
	ResponseTime time.Duration
	Error        string
	FinalURL     string // URL after redirects
}

// PingURL checks if a URL is available and accessible
func PingURL(targetURL string, options ...PingURLOptions) PingURLResult {
	var opts PingURLOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = DefaultPingOptions()
	}

	result := PingURLResult{
		Available: false,
	}

	// Validate URL format
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		result.Error = fmt.Sprintf("Invalid URL format: %v", err)
		return result
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
		targetURL = parsedURL.String()
	}

	// Create HTTP client with timeout and redirect policy
	client := &http.Client{
		Timeout: opts.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= opts.MaxRedirects {
				return fmt.Errorf("too many redirects (%d)", len(via))
			}
			return nil
		},
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "HEAD", targetURL, nil)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		return result
	}

	// Set a realistic User-Agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; URLCrawler/1.0)")

	// Record start time
	startTime := time.Now()

	// Make the request
	resp, err := client.Do(req)

	// Calculate response time
	result.ResponseTime = time.Since(startTime)

	if err != nil {
		// Try with GET if HEAD fails (some servers don't support HEAD)
		req.Method = "GET"
		resp, err = client.Do(req)
		if err != nil {
			result.Error = fmt.Sprintf("Request failed: %v", err)
			return result
		}
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.FinalURL = resp.Request.URL.String()

	// Consider 2xx and 3xx status codes as available
	// Some sites might return 403 or other codes but still be "available"
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		result.Available = true
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		// Client errors - might still be available but with restrictions
		result.Available = true
		result.Error = fmt.Sprintf("Client error: HTTP %d", resp.StatusCode)
	} else {
		// Server errors - consider unavailable
		result.Available = false
		result.Error = fmt.Sprintf("Server error: HTTP %d", resp.StatusCode)
	}

	return result
}

// IsURLAvailable is a simple convenience function that returns just the availability status
func IsURLAvailable(targetURL string) (bool, int) {
	result := PingURL(targetURL)
	return result.Available, result.StatusCode
}
