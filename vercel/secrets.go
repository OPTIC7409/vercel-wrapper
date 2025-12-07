package vercel

import (
	"context"
	"fmt"
)

// ListSecrets lists all secrets for the authenticated user or team.
func (c *Client) ListSecrets(ctx context.Context) (*ListSecretsResponse, error) {
	var resp ListSecretsResponse
	if err := c.doRequest(ctx, "GET", "/v2/secrets", nil, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetSecret retrieves a secret by ID.
func (c *Client) GetSecret(ctx context.Context, secretID string) (*Secret, error) {
	var secret Secret
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v2/secrets/%s", secretID), nil, nil, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

// CreateSecret creates a new secret.
func (c *Client) CreateSecret(ctx context.Context, req CreateSecretRequest) (*Secret, error) {
	var secret Secret
	if err := c.doRequest(ctx, "POST", "/v2/secrets", nil, req, &secret); err != nil {
		return nil, err
	}

	return &secret, nil
}

// DeleteSecret deletes a secret by ID.
func (c *Client) DeleteSecret(ctx context.Context, secretID string) error {
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/v2/secrets/%s", secretID), nil, nil, nil)
}
