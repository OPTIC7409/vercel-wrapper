package vercel

import (
	"context"
	"fmt"
)

// ListTeams lists all teams for the authenticated user.
func (c *Client) ListTeams(ctx context.Context) (*ListTeamsResponse, error) {
	var resp ListTeamsResponse
	if err := c.doRequest(ctx, "GET", "/v2/teams", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetTeam retrieves a team by ID.
func (c *Client) GetTeam(ctx context.Context, teamID string) (*Team, error) {
	var team Team
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v2/teams/%s", teamID), nil, nil, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// ListTeamMembers lists all members of a team.
func (c *Client) ListTeamMembers(ctx context.Context, teamID string) (*ListTeamMembersResponse, error) {
	var resp ListTeamMembersResponse
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v2/teams/%s/members", teamID), nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
