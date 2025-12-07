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

func TestListDeployments_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v13/deployments", r.URL.Path)

		resp := ListDeploymentsResponse{
			Deployments: []Deployment{
				{ID: "dep-1", Name: "test-deployment", State: "READY"},
			},
		}
		resp.Pagination.Count = 1

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	deployments, err := c.ListDeployments(context.Background(), "proj-1", 10, 0)
	require.NoError(t, err)
	assert.Len(t, deployments.Deployments, 1)
	assert.Equal(t, "test-deployment", deployments.Deployments[0].Name)
}

func TestGetDeployment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v13/deployments/dep-1", r.URL.Path)

		deployment := Deployment{ID: "dep-1", Name: "test-deployment", State: "READY"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deployment)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	deployment, err := c.GetDeployment(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.Equal(t, "test-deployment", deployment.Name)
	assert.Equal(t, "READY", deployment.State)
}

func TestCreateDeployment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v13/deployments", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		var req CreateDeploymentRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "test-deployment", req.Name)

		deployment := Deployment{ID: "dep-1", Name: req.Name, State: "BUILDING"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(deployment)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := CreateDeploymentRequest{
		Name:   "test-deployment",
		Target: "production",
	}
	deployment, err := c.CreateDeployment(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "test-deployment", deployment.Name)
}

func TestCreateDeployment_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "INVALID_REQUEST",
				"message": "Invalid deployment request",
			},
		})
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := CreateDeploymentRequest{Name: "test"}
	deployment, err := c.CreateDeployment(context.Background(), req)
	require.Error(t, err)
	assert.Nil(t, deployment)

	apiErr, ok := IsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "INVALID_REQUEST", apiErr.Code)
}

func TestCancelDeployment_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v13/deployments/dep-1/cancel", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	err := c.CancelDeployment(context.Background(), "dep-1")
	require.NoError(t, err)
}

func TestGetDeploymentLogs_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/deployments/dep-1/logs", r.URL.Path)

		resp := DeploymentLogsResponse{
			Logs: []DeploymentLog{
				{
					ID:        "log-1",
					Timestamp: 1609545600000,
					Message:   "Building application...",
					Type:      "stdout",
				},
				{
					ID:        "log-2",
					Timestamp: 1609545601000,
					Message:   "Build completed successfully",
					Type:      "stdout",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	logs, err := c.GetDeploymentLogs(context.Background(), "dep-1")
	require.NoError(t, err)
	assert.Len(t, logs.Logs, 2)
	assert.Equal(t, "Building application...", logs.Logs[0].Message)
	assert.Equal(t, "stdout", logs.Logs[0].Type)
}
