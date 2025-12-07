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

func TestListTeams_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/teams", r.URL.Path)

		resp := ListTeamsResponse{
			Teams: []Team{
				{ID: "team-1", Name: "Test Team", Slug: "test-team"},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	teams, err := c.ListTeams(context.Background())
	require.NoError(t, err)
	assert.Len(t, teams.Teams, 1)
	assert.Equal(t, "Test Team", teams.Teams[0].Name)
}

func TestGetTeam_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/teams/team-1", r.URL.Path)

		team := Team{ID: "team-1", Name: "Test Team", Slug: "test-team"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(team)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	team, err := c.GetTeam(context.Background(), "team-1")
	require.NoError(t, err)
	assert.Equal(t, "Test Team", team.Name)
	assert.Equal(t, "test-team", team.Slug)
}

func TestListTeamMembers_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v2/teams/team-1/members", r.URL.Path)

		resp := ListTeamMembersResponse{
			Members: []TeamMember{
				{
					User: struct {
						ID       string `json:"id"`
						Username string `json:"username"`
						Name     string `json:"name,omitempty"`
						Email    string `json:"email,omitempty"`
						Avatar   string `json:"avatar,omitempty"`
					}{
						ID:       "user-1",
						Username: "john",
						Name:     "John Doe",
					},
					Role: "OWNER",
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New("test-token", WithBaseURL(server.URL))

	members, err := c.ListTeamMembers(context.Background(), "team-1")
	require.NoError(t, err)
	assert.Len(t, members.Members, 1)
	assert.Equal(t, "john", members.Members[0].User.Username)
	assert.Equal(t, "OWNER", members.Members[0].Role)
}
