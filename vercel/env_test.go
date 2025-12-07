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

func TestUpdateEnvVar_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1/env/env-1", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)

		var req UpdateEnvVarRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "new-value", req.Value)
		assert.Equal(t, []EnvTarget{EnvTargetProduction}, req.Target)

		envVar := EnvVar{
			ID:     "env-1",
			Key:    "API_KEY",
			Value:  req.Value,
			Type:   EnvTypePlain,
			Target: req.Target,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(envVar)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := UpdateEnvVarRequest{
		Value:  "new-value",
		Target: []EnvTarget{EnvTargetProduction},
	}
	envVar, err := c.UpdateEnvVar(context.Background(), "proj-1", "env-1", req)
	require.NoError(t, err)
	assert.Equal(t, "API_KEY", envVar.Key)
	assert.Equal(t, "new-value", envVar.Value)
}
