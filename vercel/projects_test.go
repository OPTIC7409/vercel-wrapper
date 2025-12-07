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

func TestListProjects_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects", r.URL.Path)
		assert.Equal(t, "10", r.URL.Query().Get("limit"))

		resp := ListProjectsResponse{
			Projects: []Project{
				{ID: "proj-1", Name: "test-project"},
			},
		}
		resp.Pagination.Count = 1
		resp.Pagination.Total = 1

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	projects, err := c.ListProjects(context.Background(), 10, 0)
	require.NoError(t, err)
	assert.Len(t, projects.Projects, 1)
	assert.Equal(t, "test-project", projects.Projects[0].Name)
}

func TestGetProject_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/test-project", r.URL.Path)

		project := Project{ID: "proj-1", Name: "test-project"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(project)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	project, err := c.GetProject(context.Background(), "test-project")
	require.NoError(t, err)
	assert.Equal(t, "test-project", project.Name)
}

func TestGetProject_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]string{
				"code":    "NOT_FOUND",
				"message": "Project not found",
			},
		})
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	project, err := c.GetProject(context.Background(), "nonexistent")
	require.Error(t, err)
	assert.Nil(t, project)

	apiErr, ok := IsAPIError(err)
	require.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.StatusCode)
}

func TestUpdateProject_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1", r.URL.Path)
		assert.Equal(t, "PATCH", r.Method)

		var req UpdateProjectRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "updated-project", req.Name)
		assert.Equal(t, "nextjs", req.Framework)

		project := Project{ID: "proj-1", Name: req.Name, Framework: req.Framework}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(project)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	req := UpdateProjectRequest{
		Name:      "updated-project",
		Framework: "nextjs",
	}
	project, err := c.UpdateProject(context.Background(), "proj-1", req)
	require.NoError(t, err)
	assert.Equal(t, "updated-project", project.Name)
	assert.Equal(t, "nextjs", project.Framework)
}

func TestDeleteProject_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v9/projects/proj-1", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	err := c.DeleteProject(context.Background(), "proj-1")
	require.NoError(t, err)
}
