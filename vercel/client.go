package vercel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// DefaultBaseURL is the default Vercel API base URL.
	DefaultBaseURL = "https://api.vercel.com"
	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second
)

// Client is a client for interacting with the Vercel API.
type Client struct {
	token      string
	teamID     string
	baseURL    string
	httpClient *http.Client
}

// Option is a function that configures a Client.
type Option func(*Client)

// WithTeamID sets the team ID for the client.
func WithTeamID(teamID string) Option {
	return func(c *Client) {
		c.teamID = teamID
	}
}

// WithBaseURL sets the base URL for the client.
func WithBaseURL(base string) Option {
	return func(c *Client) {
		c.baseURL = base
	}
}

// WithHTTPClient sets the HTTP client for the client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// New creates a new Vercel API client.
func New(token string, opts ...Option) *Client {
	c := &Client{
		token:   token,
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// buildURL constructs a full URL from a path and query parameters.
// It automatically adds the teamId query parameter if set.
func (c *Client) buildURL(path string, query map[string]string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}

	u.Path = path

	q := u.Query()
	if c.teamID != "" {
		q.Set("teamId", c.teamID)
	}
	for k, v := range query {
		if v != "" {
			q.Set(k, v)
		}
	}
	u.RawQuery = q.Encode()

	return u.String(), nil
}

// doRequest performs an HTTP request and handles the response.
func (c *Client) doRequest(ctx context.Context, method, path string, query map[string]string, body interface{}, v interface{}) error {
	reqURL, err := c.buildURL(path, query)
	if err != nil {
		return err
	}

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		apiErr := &APIError{
			StatusCode: resp.StatusCode,
			RawBody:    respBody,
			Message:    http.StatusText(resp.StatusCode),
		}

		// Try to unmarshal error response
		var errorResp struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(respBody, &errorResp); err == nil {
			apiErr.Code = errorResp.Error.Code
			if errorResp.Error.Message != "" {
				apiErr.Message = errorResp.Error.Message
			}
		}

		return apiErr
	}

	if v != nil {
		if err := json.Unmarshal(respBody, v); err != nil {
			return fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return nil
}

