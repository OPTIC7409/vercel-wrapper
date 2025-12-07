package vercel

import (
	"context"
	"fmt"
	"strconv"
)

// ListAliases lists all aliases, optionally filtered by project or deployment.
func (c *Client) ListAliases(ctx context.Context, projectID, deploymentID string, limit int) (*ListAliasesResponse, error) {
	query := make(map[string]string)
	if projectID != "" {
		query["projectId"] = projectID
	}
	if deploymentID != "" {
		query["deploymentId"] = deploymentID
	}
	if limit > 0 {
		query["limit"] = strconv.Itoa(limit)
	}

	var resp ListAliasesResponse
	if err := c.doRequest(ctx, "GET", "/v4/aliases", query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// ListDeploymentAliases lists all aliases for a specific deployment.
func (c *Client) ListDeploymentAliases(ctx context.Context, deploymentID string) ([]Alias, error) {
	var resp struct {
		Aliases []Alias `json:"aliases"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v2/deployments/%s/aliases", deploymentID), nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Aliases, nil
}

// CreateAlias creates a new alias.
func (c *Client) CreateAlias(ctx context.Context, req CreateAliasRequest) (*Alias, error) {
	var alias Alias
	if err := c.doRequest(ctx, "POST", "/v4/aliases", nil, req, &alias); err != nil {
		return nil, err
	}

	return &alias, nil
}

// DeleteAlias deletes an alias by ID.
func (c *Client) DeleteAlias(ctx context.Context, aliasID string) error {
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/v4/aliases/%s", aliasID), nil, nil, nil)
}
