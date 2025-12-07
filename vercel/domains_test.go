package vercel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDomains_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1/domains", r.URL.Path)

		resp := struct {
			Domains []Domain `json:"domains"`
		}{
			Domains: []Domain{
				{
					ID:        "dom-1",
					Name:      "example.com",
					Verified:  true,
					ProjectID: "proj-1",
				},
			},
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	domains, err := c.ListDomains(context.Background(), "proj-1")
	require.NoError(t, err)
	assert.Len(t, domains, 1)
	assert.Equal(t, "example.com", domains[0].Name)
}

func TestGetDomain_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1/domains/example.com", r.URL.Path)

		domain := Domain{
			ID:        "dom-1",
			Name:      "example.com",
			Verified:  true,
			ProjectID: "proj-1",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(domain)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	domain, err := c.GetDomain(context.Background(), "proj-1", "example.com")
	require.NoError(t, err)
	assert.Equal(t, "example.com", domain.Name)
	assert.True(t, domain.Verified)
}

func TestCreateDomain_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1/domains", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateDomainRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "newdomain.com", req.Name)

		domain := Domain{
			ID:        "dom-new",
			Name:      req.Name,
			Verified:  false,
			ProjectID: "proj-1",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(domain)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := CreateDomainRequest{
		Name: "newdomain.com",
	}
	domain, err := c.CreateDomain(context.Background(), "proj-1", req)
	require.NoError(t, err)
	assert.Equal(t, "newdomain.com", domain.Name)
}

func TestDeleteDomain_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1/domains/example.com", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	err := c.DeleteDomain(context.Background(), "proj-1", "example.com")
	require.NoError(t, err)
}

func TestGetDomain_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "NOT_FOUND",
				"message": "Domain not found",
			},
		})
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	domain, err := c.GetDomain(context.Background(), "proj-1", "nonexistent.com")
	require.Error(t, err)
	assert.Nil(t, domain)

	apiErr, ok := IsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
	assert.Equal(t, "NOT_FOUND", apiErr.Code)
}

