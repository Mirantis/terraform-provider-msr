package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CreateRepo struct {
	ImmutableTags    bool   `json:"immutableTags"`
	LongDescription  string `json:"longDescription"`
	Name             string `json:"name"`
	ScanOnPush       bool   `json:"scanOnPush"`
	ShortDescription string `json:"shortDescription"`
	TagLimit         int    `json:"tagLimit"`
	Visibility       string `json:"visibility" enum:"public|private"`
}

type UpdateRepo struct {
	ImmutableTags    bool   `json:"immutableTags"`
	LongDescription  string `json:"longDescription"`
	ScanOnPush       bool   `json:"scanOnPush"`
	ShortDescription string `json:"shortDescription"`
	TagLimit         int    `json:"tagLimit"`
	Visibility       string `json:"visibility" enum:"public|private"`
}

type ResponseRepo struct {
	ID               string `json:"id"`
	ImmutableTags    bool   `json:"immutableTags"`
	LongDescription  string `json:"longDescription"`
	Name             string `json:"name"`
	Namespace        string `json:"namespace"`
	NamespaceType    string `json:"namespaceType"`
	Pulls            int    `json:"pulls"`
	Pushes           int    `json:"pushes"`
	ScanOnPush       bool   `json:"scanOnPush"`
	ShortDescription string `json:"shortDescription"`
	TagLimit         int    `json:"tagLimit"`
	Visibility       string `json:"visibility" enum:"public|private"`
}

// CreateRepo creates a repo in MSR.
func (c *Client) CreateRepo(ctx context.Context, orgName string, repo CreateRepo) (ResponseRepo, error) {
	if (repo == CreateRepo{}) {
		return ResponseRepo{}, fmt.Errorf("creating repo failed. %w: %+v", ErrEmptyStruct, repo)
	}
	body, err := json.Marshal(repo)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("creating repo %s failed. %w: %s", repo.Name, ErrMarshaling, err)
	}
	url := fmt.Sprintf("%s/%s", c.createMsrUrl("repositories"), orgName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("creating repo %s failed. %w: %s", repo.Name, ErrRequestCreation, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("creating repo %s failed. %w", repo.Name, err)
	}

	resRepo := ResponseRepo{}
	if err := json.Unmarshal(resBody, &resRepo); err != nil {
		return ResponseRepo{}, fmt.Errorf("creating repo %s failed. %w: %s", repo.Name, ErrUnmarshaling, err)
	}

	return resRepo, nil
}

// ReadRepo method retrieves a repo from the MSR endpoint.
func (c *Client) ReadRepo(ctx context.Context, orgName string, repoName string) (ResponseRepo, error) {
	url := fmt.Sprintf("%s/%s/%s", c.createMsrUrl("repositories"), orgName, repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("reading repo %s failed in MSR client: %w: %s", repoName, ErrRequestCreation, err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("reading repo %s failed in MSR client: %w", repoName, err)
	}

	repo := ResponseRepo{}
	if err := json.Unmarshal(body, &repo); err != nil {
		return ResponseRepo{}, fmt.Errorf("reading repo %s failed in MSR client: %w: %s", repoName, ErrUnmarshaling, err)
	}

	return repo, nil
}

// UpdateRepo updates a repo in the MSR endpoint.
func (c *Client) UpdateRepo(ctx context.Context, orgName string, repoName string, repo UpdateRepo) (ResponseRepo, error) {
	if (repo == UpdateRepo{}) {
		return ResponseRepo{}, fmt.Errorf("updating repo failed. %w: %s", ErrEmptyStruct, repoName)
	}
	url := fmt.Sprintf("%s/%s/%s", c.createMsrUrl("repositories"), orgName, repoName)

	body, err := json.Marshal(repo)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("update repo %s failed in MSR client: %w: %s", repoName, ErrMarshaling, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))

	if err != nil {
		return ResponseRepo{}, fmt.Errorf("updating repo %s failed. %w: %s", repoName, ErrRequestCreation, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponseRepo{}, fmt.Errorf("updating repo %s failed. %w", repoName, err)
	}

	rRepo := ResponseRepo{}
	if json.Unmarshal(resBody, &rRepo) != nil {
		return ResponseRepo{}, fmt.Errorf("update repo %s failed in MSR client: %w: %s", repoName, ErrUnmarshaling, err)
	}
	return rRepo, nil
}

// DeleteRepo deletes a repo from MSR.
func (c *Client) DeleteRepo(ctx context.Context, orgName string, repoName string) error {
	url := fmt.Sprintf("%s/%s/%s", c.createMsrUrl("repositories"), orgName, repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("delete repo %s failed in MSR client: %w", repoName, err)
	}

	if _, err = c.doRequest(req); err != nil {
		return fmt.Errorf("deleting repo %s failed. %w", repoName, err)
	}

	return err
}
