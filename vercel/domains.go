package vercel

import (
	"context"
	"fmt"
)

// ListDomains lists all domains for a project.
func (c *Client) ListDomains(ctx context.Context, projectIDOrName string) ([]Domain, error) {
	var resp struct {
		Domains []Domain `json:"domains"`
	}

	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v9/projects/%s/domains", projectIDOrName), nil, nil, &resp); err != nil {
		return nil, err
	}

	return resp.Domains, nil
}

// GetDomain retrieves a domain by name for a project.
func (c *Client) GetDomain(ctx context.Context, projectIDOrName, domainName string) (*Domain, error) {
	var domain Domain
	if err := c.doRequest(ctx, "GET", fmt.Sprintf("/v9/projects/%s/domains/%s", projectIDOrName, domainName), nil, nil, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// CreateDomain adds a domain to a project.
func (c *Client) CreateDomain(ctx context.Context, projectIDOrName string, req CreateDomainRequest) (*Domain, error) {
	var domain Domain
	if err := c.doRequest(ctx, "POST", fmt.Sprintf("/v9/projects/%s/domains", projectIDOrName), nil, req, &domain); err != nil {
		return nil, err
	}

	return &domain, nil
}

// DeleteDomain removes a domain from a project.
func (c *Client) DeleteDomain(ctx context.Context, projectIDOrName, domainName string) error {
	return c.doRequest(ctx, "DELETE", fmt.Sprintf("/v9/projects/%s/domains/%s", projectIDOrName, domainName), nil, nil, nil)
}

