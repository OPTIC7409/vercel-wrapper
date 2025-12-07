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

func TestListSecrets_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/secrets", r.URL.Path)

		resp := ListSecretsResponse{
			Secrets: []Secret{
				{ID: "secret-1", Name: "API_KEY"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	secrets, err := c.ListSecrets(context.Background())
	require.NoError(t, err)
	assert.Len(t, secrets.Secrets, 1)
	assert.Equal(t, "API_KEY", secrets.Secrets[0].Name)
}

func TestGetSecret_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/secrets/secret-1", r.URL.Path)

		secret := Secret{ID: "secret-1", Name: "API_KEY"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(secret)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	secret, err := c.GetSecret(context.Background(), "secret-1")
	require.NoError(t, err)
	assert.Equal(t, "API_KEY", secret.Name)
}

func TestCreateSecret_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/secrets", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateSecretRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "API_KEY", req.Name)
		assert.Equal(t, "secret-value", req.Value)

		secret := Secret{
			ID:    "secret-1",
			Name:  req.Name,
			Value: req.Value,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(secret)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := CreateSecretRequest{
		Name:  "API_KEY",
		Value: "secret-value",
	}
	secret, err := c.CreateSecret(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "API_KEY", secret.Name)
	assert.Equal(t, "secret-value", secret.Value)
}

func TestDeleteSecret_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/secrets/secret-1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	err := c.DeleteSecret(context.Background(), "secret-1")
	require.NoError(t, err)
}
