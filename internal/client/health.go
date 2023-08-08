package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HealthResponse structure from MSR.
type HealthResponse struct {
	Error   string `json:"error"`
	Healthy bool   `json:"healthy"`
}

type MSRVersion struct {
	Version string `json:"version"`
}

// IsHealthy method - checking the MSR health endpoint.
func (c *Client) IsHealthy(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/health", c.MsrURL), nil)
	if err != nil {
		return false, fmt.Errorf("health ping failed. %w: %s", ErrRequestCreation, err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return false, fmt.Errorf("health ping failed. %w", err)
	}

	hResponse := HealthResponse{}

	if err := json.Unmarshal(body, &hResponse); err != nil {
		return false, fmt.Errorf("health ping failed. %w: %s", ErrUnmarshaling, err)
	}

	return hResponse.Healthy, nil
}

// VerifGetMSRVersion returns the MSR instance version given the credentials are valid.
func (c *Client) GetMSRVersion(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.createMsrUrl("admin/version"), nil)
	if err != nil {
		return "", fmt.Errorf("getting MSR version failed. %w: %s", ErrRequestCreation, err)
	}
	body, err := c.doRequest(req)
	if err != nil {
		return "", fmt.Errorf("getting MSR version failed. %w", err)
	}

	ver := MSRVersion{}
	if err := json.Unmarshal(body, &ver); err != nil {
		return "", fmt.Errorf("getting MSR version failed. %w: %s", ErrUnmarshaling, err)
	}

	return ver.Version, err
}
