package vercel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// logResponse logs a cleaned, formatted JSON response for debugging
func logResponse(t *testing.T, label string, data interface{}) {
	t.Helper()
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		t.Logf("%s: Failed to marshal response: %v", label, err)
		return
	}
	t.Logf("\n=== %s ===\n%s\n", label, string(jsonData))
}

// TestAllEndpoints_WithLogging tests all endpoints and logs cleaned responses
func TestAllEndpoints_WithLogging(t *testing.T) {
	t.Run("Projects", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects", r.URL.Path)
			assert.Equal(t, "10", r.URL.Query().Get("limit"))
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

			resp := ListProjectsResponse{
				Projects: []Project{
					{
						ID:        "proj-abc123",
						Name:      "my-awesome-project",
						Framework: "nextjs",
						CreatedAt: 1609459200000,
						UpdatedAt: 1609545600000,
						TeamID:    "team-xyz789",
					},
					{
						ID:        "proj-def456",
						Name:      "another-project",
						Framework: "react",
						CreatedAt: 1609632000000,
						UpdatedAt: 1609718400000,
					},
				},
			}
			resp.Pagination.Count = 2
			resp.Pagination.Limit = 10
			resp.Pagination.Offset = 0
			resp.Pagination.Total = 2

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		projects, err := c.ListProjects(context.Background(), 10, 0)
		require.NoError(t, err)
		logResponse(t, "ListProjects Response", projects)

		assert.Len(t, projects.Projects, 2)
		assert.Equal(t, "my-awesome-project", projects.Projects[0].Name)
		assert.Equal(t, 2, projects.Pagination.Total)
	})

	t.Run("GetProject", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/test-project", r.URL.Path)

			project := Project{
				ID:        "proj-abc123",
				Name:      "test-project",
				Framework: "nextjs",
				CreatedAt: 1609459200000,
				UpdatedAt: 1609545600000,
				TeamID:    "team-xyz789",
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(project)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		project, err := c.GetProject(context.Background(), "test-project")
		require.NoError(t, err)
		logResponse(t, "GetProject Response", project)

		assert.Equal(t, "test-project", project.Name)
		assert.Equal(t, "proj-abc123", project.ID)
	})

	t.Run("Deployments", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v13/deployments", r.URL.Path)
			assert.Equal(t, "proj-abc123", r.URL.Query().Get("projectId"))
			assert.Equal(t, "10", r.URL.Query().Get("limit"))

			resp := ListDeploymentsResponse{
				Deployments: []Deployment{
					{
						ID:        "dpl-123456",
						Name:      "my-deployment",
						URL:       "https://my-deployment.vercel.app",
						State:     "READY",
						Target:    "production",
						CreatedAt: 1609459200000,
						ReadyAt:   1609459260000,
						ProjectID: "proj-abc123",
					},
					{
						ID:        "dpl-789012",
						Name:      "preview-deployment",
						URL:       "https://preview-deployment.vercel.app",
						State:     "BUILDING",
						Target:    "preview",
						CreatedAt: 1609459300000,
						ProjectID: "proj-abc123",
					},
				},
			}
			resp.Pagination.Count = 2
			resp.Pagination.Limit = 10

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		deployments, err := c.ListDeployments(context.Background(), "proj-abc123", 10, 0)
		require.NoError(t, err)
		logResponse(t, "ListDeployments Response", deployments)

		assert.Len(t, deployments.Deployments, 2)
		assert.Equal(t, "READY", deployments.Deployments[0].State)
		assert.Equal(t, "BUILDING", deployments.Deployments[1].State)
	})

	t.Run("GetDeployment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v13/deployments/dpl-123456", r.URL.Path)

			deployment := Deployment{
				ID:        "dpl-123456",
				Name:      "my-deployment",
				URL:       "https://my-deployment.vercel.app",
				State:     "READY",
				Target:    "production",
				CreatedAt: 1609459200000,
				ReadyAt:   1609459260000,
				ProjectID: "proj-abc123",
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(deployment)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		deployment, err := c.GetDeployment(context.Background(), "dpl-123456")
		require.NoError(t, err)
		logResponse(t, "GetDeployment Response", deployment)

		assert.Equal(t, "dpl-123456", deployment.ID)
		assert.Equal(t, "READY", deployment.State)
	})

	t.Run("CreateDeployment", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v13/deployments", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var req CreateDeploymentRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			logResponse(t, "CreateDeployment Request", req)

			deployment := Deployment{
				ID:        "dpl-new123",
				Name:      req.Name,
				URL:       fmt.Sprintf("https://%s.vercel.app", req.Name),
				State:     "BUILDING",
				Target:    req.Target,
				CreatedAt: 1609459200000,
				ProjectID: req.Project,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(deployment)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		req := CreateDeploymentRequest{
			Name:    "new-deployment",
			Project: "proj-abc123",
			Target:  "production",
			Files: []DeploymentFile{
				{
					File: "index.html",
					Data: "PGh0bWw+SGVsbG8gV29ybGQ8L2h0bWw+",
				},
			},
			Env: map[string]string{
				"API_KEY": "secret-value",
			},
		}

		deployment, err := c.CreateDeployment(context.Background(), req)
		require.NoError(t, err)
		logResponse(t, "CreateDeployment Response", deployment)

		assert.Equal(t, "new-deployment", deployment.Name)
		assert.Equal(t, "BUILDING", deployment.State)
	})

	t.Run("Environment Variables", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/env", r.URL.Path)

			resp := struct {
				EnvVars []EnvVar `json:"env"`
			}{
				EnvVars: []EnvVar{
					{
						ID:        "env-123",
						Key:       "API_KEY",
						Type:      EnvTypeSecret,
						Target:    []EnvTarget{EnvTargetProduction, EnvTargetPreview},
						CreatedAt: 1609459200000,
						UpdatedAt: 1609545600000,
					},
					{
						ID:        "env-456",
						Key:       "DATABASE_URL",
						Type:      EnvTypePlain,
						Target:    []EnvTarget{EnvTargetProduction},
						CreatedAt: 1609459200000,
						UpdatedAt: 1609545600000,
					},
				},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		envVars, err := c.ListEnvVars(context.Background(), "proj-abc123")
		require.NoError(t, err)
		logResponse(t, "ListEnvVars Response", envVars)

		assert.Len(t, envVars, 2)
		assert.Equal(t, "API_KEY", envVars[0].Key)
		assert.Equal(t, EnvTypeSecret, envVars[0].Type)
	})

	t.Run("CreateEnvVar", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/env", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var req CreateEnvVarRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			logResponse(t, "CreateEnvVar Request", req)

			envVar := EnvVar{
				ID:        "env-new789",
				Key:       req.Key,
				Type:      req.Type,
				Target:    req.Target,
				CreatedAt: 1609459200000,
				UpdatedAt: 1609459200000,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(envVar)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		req := CreateEnvVarRequest{
			Key:    "NEW_VAR",
			Value:  "new-value",
			Type:   EnvTypePlain,
			Target: []EnvTarget{EnvTargetProduction},
		}

		envVar, err := c.CreateEnvVar(context.Background(), "proj-abc123", req)
		require.NoError(t, err)
		logResponse(t, "CreateEnvVar Response", envVar)

		assert.Equal(t, "NEW_VAR", envVar.Key)
		assert.Equal(t, EnvTypePlain, envVar.Type)
	})

	t.Run("Domains", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/domains", r.URL.Path)

			resp := struct {
				Domains []Domain `json:"domains"`
			}{
				Domains: []Domain{
					{
						ID:          "dom-123",
						Name:        "example.com",
						ServiceType: "vercel-dns",
						Nameservers: []string{"ns1.vercel-dns.com", "ns2.vercel-dns.com"},
						Verified:    true,
						CreatedAt:   1609459200000,
						UpdatedAt:   1609545600000,
						ProjectID:   "proj-abc123",
						CDNEnabled:  true,
					},
					{
						ID:          "dom-456",
						Name:        "www.example.com",
						ServiceType: "external",
						Verified:    false,
						CreatedAt:   1609459300000,
						UpdatedAt:   1609545700000,
						ProjectID:   "proj-abc123",
					},
				},
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		domains, err := c.ListDomains(context.Background(), "proj-abc123")
		require.NoError(t, err)
		logResponse(t, "ListDomains Response", domains)

		assert.Len(t, domains, 2)
		assert.Equal(t, "example.com", domains[0].Name)
		assert.True(t, domains[0].Verified)
	})

	t.Run("GetDomain", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/domains/example.com", r.URL.Path)

			domain := Domain{
				ID:          "dom-123",
				Name:        "example.com",
				ServiceType: "vercel-dns",
				Nameservers: []string{"ns1.vercel-dns.com", "ns2.vercel-dns.com"},
				Verified:    true,
				CreatedAt:   1609459200000,
				UpdatedAt:   1609545600000,
				ProjectID:   "proj-abc123",
				CDNEnabled:  true,
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(domain)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		domain, err := c.GetDomain(context.Background(), "proj-abc123", "example.com")
		require.NoError(t, err)
		logResponse(t, "GetDomain Response", domain)

		assert.Equal(t, "example.com", domain.Name)
		assert.True(t, domain.Verified)
		assert.Equal(t, "vercel-dns", domain.ServiceType)
	})

	t.Run("CreateDomain", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/domains", r.URL.Path)
			assert.Equal(t, "POST", r.Method)

			var req CreateDomainRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			logResponse(t, "CreateDomain Request", req)

			domain := Domain{
				ID:          "dom-new789",
				Name:        req.Name,
				ServiceType: "vercel-dns",
				Verified:    false,
				CreatedAt:   1609459200000,
				UpdatedAt:   1609459200000,
				ProjectID:   "proj-abc123",
			}

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(domain)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		req := CreateDomainRequest{
			Name:      "newdomain.com",
			GitBranch: "main",
		}

		domain, err := c.CreateDomain(context.Background(), "proj-abc123", req)
		require.NoError(t, err)
		logResponse(t, "CreateDomain Response", domain)

		assert.Equal(t, "newdomain.com", domain.Name)
		assert.Equal(t, "proj-abc123", domain.ProjectID)
	})

	t.Run("DeleteDomain", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/v9/projects/proj-abc123/domains/example.com", r.URL.Path)
			assert.Equal(t, "DELETE", r.Method)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		c := New("test-token", WithBaseURL(server.URL))

		err := c.DeleteDomain(context.Background(), "proj-abc123", "example.com")
		require.NoError(t, err)
		t.Log("DeleteDomain: Successfully deleted domain")
	})
}

