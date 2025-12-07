package vercel

import (
	"context"
	"fmt"
	"strconv"
)

// ListDeployments lists deployments, optionally filtered by project.
func (c *Client) ListDeployments(ctx context.Context, projectIDOrName string, limit, since int) (*ListDeploymentsResponse, error) {
	query := make(map[string]string)
	if projectIDOrName != "" {
		query["projectId"] = projectIDOrName
	}
	if limit > 0 {
		query["limit"] = strconv.Itoa(limit)
	}
	if since > 0 {
		query["since"] = strconv.Itoa(since)
	}

	var resp ListDeploymentsResponse
	if err := c.doRequest(ctx, "GET", "/v13/deployments", query, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetDeployment retrieves a deployment by ID.
func (c *Client) GetDeployment(ctx context.Context, id string) (*Deployment, error) {
	var deployment Deployment
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v13/deployments/%s", id), nil, nil, &deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}

// CreateDeployment creates a new deployment.
func (c *Client) CreateDeployment(ctx context.Context, req CreateDeploymentRequest) (*Deployment, error) {
	var deployment Deployment
	if err := c.doRequest(ctx, "POST", "/v13/deployments", nil, req, &deployment); err != nil {
		return nil, err
	}

	return &deployment, nil
}

// CancelDeployment cancels a deployment by ID.
func (c *Client) CancelDeployment(ctx context.Context, id string) error {
	return c.doRequest(ctx, "PATCH", fmt.Sprintf("/v13/deployments/%s/cancel", id), nil, nil, nil)
}

