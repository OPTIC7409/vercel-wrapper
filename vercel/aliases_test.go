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

func TestListAliases_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v4/aliases", r.URL.Path)
		assert.Equal(t, "proj-1", r.URL.Query().Get("projectId"))

		resp := ListAliasesResponse{
			Aliases: []Alias{
				{ID: "alias-1", Alias: "example.com", ProjectID: "proj-1"},
			},
		}
		resp.Pagination.Count = 1
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	aliases, err := c.ListAliases(context.Background(), "proj-1", "", 10)
	require.NoError(t, err)
	assert.Len(t, aliases.Aliases, 1)
	assert.Equal(t, "example.com", aliases.Aliases[0].Alias)
}

func TestListDeploymentAliases_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/deployments/dep-1/aliases", r.URL.Path)

		var resp struct {
			Aliases []Alias `json:"aliases"`
		}
		resp.Aliases = []Alias{
			{ID: "alias-1", Alias: "example.com"},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	aliases, err := c.ListDeploymentAliases(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.Len(t, aliases, 1)
	assert.Equal(t, "example.com", aliases[0].Alias)
}

func TestCreateAlias_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v4/aliases", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateAliasRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "example.com", req.Alias)

		alias := Alias{ID: "alias-1", Alias: req.Alias}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alias)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := CreateAliasRequest{
		Alias: "example.com",
	}
	alias, err := c.CreateAlias(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "example.com", alias.Alias)
}

func TestDeleteAlias_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v4/aliases/alias-1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	err := c.DeleteAlias(context.Background(), "alias-1")
	require.NoError(t, err)
}
