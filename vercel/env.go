package vercel

import (
	"context"
	"fmt"
)

// ListEnvVars lists all environment variables for a project.
func (c *Client) ListEnvVars(ctx context.Context, projectIDOrName string) ([]EnvVar, error) {
	var resp struct {
		EnvVars []EnvVar `json:"env"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v9/projects/%s/env", projectIDOrName), nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp.EnvVars, nil
}

// CreateEnvVar creates a new environment variable for a project.
func (c *Client) CreateEnvVar(ctx context.Context, projectIDOrName string, req CreateEnvVarRequest) (*EnvVar, error) {
	var envVar EnvVar
	if err := c.doRequest(ctx, "POST", fmt.Sprintf("/v9/projects/%s/env", projectIDOrName), nil, req, &envVar); err != nil {
		return nil, err
	}

	return &envVar, nil
}

// UpdateEnvVar updates an environment variable by ID.
func (c *Client) UpdateEnvVar(ctx context.Context, projectIDOrName, envID string, req UpdateEnvVarRequest) (*EnvVar, error) {
	var envVar EnvVar
	if err := c.doRequest(ctx, "PATCH", fmt.Sprintf("/v9/projects/%s/env/%s", projectIDOrName, envID), nil, req, &envVar); err != nil {
		return nil, err
	}

	return &envVar, nil
}

// DeleteEnvVar deletes an environment variable by ID.
func (c *Client) DeleteEnvVar(ctx context.Context, projectIDOrName, envID string) error {
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/v9/projects/%s/env/%s", projectIDOrName, envID), nil, nil, nil)
}
