package vercel

import (
	"context"
	"fmt"
	"strconv"
)

// ListProjects lists all projects for the authenticated user or team.
func (c *Client) ListProjects(ctx context.Context, limit, offset int) (*ListProjectsResponse, error) {
	query := make(map[string]string)
	if limit > 0 {
		query["limit"] = strconv.Itoa(limit)
	}
	if offset > 0 {
		query["offset"] = strconv.Itoa(offset)
	}

	var resp ListProjectsResponse
	if err := c.doRequest(ctx, "GET", "/v9/projects", query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetProject retrieves a project by ID or name.
func (c *Client) GetProject(ctx context.Context, idOrName string) (*Project, error) {
	var project Project
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v9/projects/%s", idOrName), nil, nil, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

