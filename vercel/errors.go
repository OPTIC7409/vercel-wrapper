package vercel

import "fmt"

// APIError represents an error response from the Vercel API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RawBody    []byte
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("vercel: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("vercel: %s (status=%d)", e.Message, e.StatusCode)
}

// IsAPIError checks if an error is an APIError and returns it if so.
func IsAPIError(err error) (*APIError, bool) {
	if err == nil {
		return nil, false
	}
	apiErr, ok := err.(*APIError)
	return apiErr, ok
}

