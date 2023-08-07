package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Team struct {
	Description  string `json:"description"`
	ID           string `json:"id"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
	OrgID        string `json:"orgID"`
}

type teamUsers struct {
	Members []struct {
		IsAdmin bool `json:"isAdmin"`
		// There is aditional fields available no present in Account
		Member ResponseAccount `json:"member"`
	} `json:"members"`
}

// CreateTeam creates a team in Enzin.
func (c *Client) CreateTeam(ctx context.Context, orgID string, team Team) (Team, error) {
	body, err := json.Marshal(team)
	if err != nil {
		return Team{}, fmt.Errorf("create team failed in MSR client: %w", err)
	}
	url := fmt.Sprintf("%s/%s/teams", c.createEnziUrl("accounts"), orgID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return Team{}, fmt.Errorf("request creation failed in MSR client: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return Team{}, err
	}

	if err := json.Unmarshal(resBody, &team); err != nil {
		return Team{}, fmt.Errorf("create a team failed in MSR client: %w", err)
	}

	return team, nil
}

// ReadTeam method retrieves a team from the enzi endpoint.
func (c *Client) ReadTeam(ctx context.Context, orgID string, teamID string) (Team, error) {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, teamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return Team{}, err
	}

	body, err := c.doRequest(req)

	if err != nil {
		return Team{}, fmt.Errorf("org_id: %s, team_id: %s - %w", orgID, teamID, err)
	}

	team := Team{}
	if err := json.Unmarshal(body, &team); err != nil {
		return Team{}, fmt.Errorf("read team failed in MSR client: %w", err)
	}

	return team, nil
}

// UpdateTeam updates a team in the enzi endpoint.
func (c *Client) UpdateTeam(ctx context.Context, orgID string, team Team) (Team, error) {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, team.ID)

	body, err := json.Marshal(team)
	if err != nil {
		return Team{}, fmt.Errorf("update team failed in MSR client: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))

	if err != nil {
		return Team{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)

	if err != nil {
		return Team{}, err
	}

	if json.Unmarshal(resBody, &team) != nil {
		return Team{}, fmt.Errorf("update team failed in MSR client: %w", err)
	}
	return team, nil
}

// DeleteTeam deletes a team from Enzi.
func (c *Client) DeleteTeam(ctx context.Context, orgID string, teamID string) error {
	url := fmt.Sprintf("%s/%s/teams/%s", c.createEnziUrl("accounts"), orgID, teamID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)

	if err != nil {
		return fmt.Errorf("delete team failed in MSR client: %w", err)
	}

	_, err = c.doRequest(req)

	return err
}

// AddUserToTeam adds user to a team.
func (c *Client) AddUserToTeam(ctx context.Context, orgID string, teamID string, user ResponseAccount) error {
	body, err := json.Marshal(map[string]bool{"isAdmin": user.IsAdmin})
	if err != nil {
		return fmt.Errorf("add user to team failed in MSR client: %w", err)
	}
	endpoint := c.createEnziUrl(fmt.Sprintf("accounts/%s/teams/%s/members/%s", orgID, teamID, user.ID))
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("adding user to team failed in MSR client: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return fmt.Errorf("the REST call returned error: %w", err)
	}

	return nil
}

// GetTeamUsers retrieves the users of a given team.
func (c *Client) GetTeamUsers(ctx context.Context, orgID string, teamID string) (teamUsers, error) {
	endpoint := c.createEnziUrl(fmt.Sprintf("accounts/%s/teams/%s/members", orgID, teamID))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return teamUsers{}, fmt.Errorf("retrieving user of team failed in MSR client: %w", err)
	}

	resBody, err := c.doRequest(req)
	fmt.Printf("Body %+v", string(resBody))
	if err != nil {
		return teamUsers{}, fmt.Errorf("retrieving user of team failed in MSR client: %w", err)
	}

	tUsers := teamUsers{}
	if err := json.Unmarshal(resBody, &tUsers); err != nil {
		return teamUsers{}, fmt.Errorf("retrieving user of team failed in MSR client: %w", err)

	}

	return tUsers, nil
}

// DeleteUserFromTeam deletes a user from a given team.
func (c *Client) DeleteUserFromTeam(ctx context.Context, orgID string, teamID string, userID string) error {
	// Check if the user exists -> then proceed to delete it.
	if _, err := c.ReadAccount(ctx, userID); err != nil {
		return err
	}
	endpoint := c.createEnziUrl(fmt.Sprintf("accounts/%s/teams/%s/members/%s", orgID, teamID, userID))
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("deleting user from team failed in MSR client: %w", err)
	}

	if _, err := c.doRequest(req); err != nil {
		return fmt.Errorf("deleting user from team failed in MSR client: %w", err)
	}

	return nil
}

// UpdateTeamUsers updates a team user base to match the latest state defined by Terraform.
func (c *Client) UpdateTeamUsers(ctx context.Context, orgID string, teamID string, newUsers []string) error {

	tUsers, err := c.GetTeamUsers(ctx, orgID, teamID)
	if err != nil {
		return fmt.Errorf("updating team users failed: %w", err)
	}
	for _, user := range tUsers.Members {
		if err := c.DeleteUserFromTeam(ctx, orgID, teamID, user.Member.ID); err != nil {
			return fmt.Errorf("updating team users failed: %w", err)
		}
	}

	for _, u := range newUsers {
		a := ResponseAccount{
			ID: u,
		}
		// Log the failure in the future if you fail to add an user to a team
		if err := c.AddUserToTeam(ctx, orgID, teamID, a); err != nil {
			continue
		}
	}

	return nil
}
