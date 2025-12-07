package vercel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := New("test-token")
	assert.Equal(t, "test-token", c.token)
	assert.Equal(t, DefaultBaseURL, c.baseURL)
	assert.NotNil(t, c.httpClient)
}

func TestWithTeamID(t *testing.T) {
	c := New("test-token", WithTeamID("team-123"))
	assert.Equal(t, "team-123", c.teamID)
}

func TestWithBaseURL(t *testing.T) {
	c := New("test-token", WithBaseURL("https://custom.api.com"))
	assert.Equal(t, "https://custom.api.com", c.baseURL)
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 10 * time.Second}
	c := New("test-token", WithHTTPClient(customClient))
	assert.Equal(t, customClient, c.httpClient)
}

func TestBuildURL(t *testing.T) {
	c := New("test-token", WithTeamID("team-123"))

	url, err := c.buildURL("/v9/projects", map[string]string{"limit": "10"})
	require.NoError(t, err)
	assert.Contains(t, url, "teamId=team-123")
	assert.Contains(t, url, "limit=10")
}

func TestDoRequest_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	var result map[string]string
	err := c.doRequest(context.Background(), "GET", "/test", nil, nil, &result)
	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
}

func TestDoRequest_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "NOT_FOUND",
				"message": "Resource not found",
			},
		})
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	var result map[string]string
	err := c.doRequest(context.Background(), "GET", "/test", nil, nil, &result)
	require.Error(t, err)

	apiErr, ok := IsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "NOT_FOUND", apiErr.Code)
	assert.Equal(t, "Resource not found", apiErr.Message)
}

func TestDoRequest_ErrorNoJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	var result map[string]string
	err := c.doRequest(context.Background(), "GET", "/test", nil, nil, &result)
	require.Error(t, err)

	apiErr, ok := IsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, "Internal Server Error", string(apiErr.RawBody))
}
