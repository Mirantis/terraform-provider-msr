package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateAccount struct
type CreateAccount struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Password   string `json:"password"`
	FullName   string `json:"fullName,omitempty"`
	IsActive   bool   `json:"isActive,omitempty"`
	IsAdmin    bool   `json:"isAdmin,omitempty"`
	IsOrg      bool   `json:"isOrg,omitempty"`
	SearchLDAP bool   `json:"searchLDAP,omitempty"`
}

// UpdateAccount struct
type UpdateAccount struct {
	FullName string `json:"fullName,omitempty"`
	IsActive bool   `json:"isActive,omitempty"`
	IsAdmin  bool   `json:"isAdmin,omitempty"`
}

// ResponseAccount struct
type ResponseAccount struct {
	Name         string `json:"name"`
	ID           string `json:"id"`
	FullName     string `json:"fullName,omitempty"`
	IsActive     bool   `json:"isActive"`
	IsAdmin      bool   `json:"isAdmin"`
	IsOrg        bool   `json:"isOrg"`
	IsImported   bool   `json:"isImported"`
	OnDemand     bool   `json:"onDemand"`
	OtpEnabled   bool   `json:"otpEnabled"`
	MembersCount int    `json:"membersCount"`
	TeamsCount   int    `json:"teamsCount"`
}

// Account filters enum
type AccountFilter string

const (
	Users         AccountFilter = "user"
	Orgs          AccountFilter = "orgs"
	Admins        AccountFilter = "admins"
	NonAdmins     AccountFilter = "non-admins"
	ActiveUsers   AccountFilter = "active-users"
	InactiveUsers AccountFilter = "inactive-users"
)

// APIFormOfFilter is a string readable form of the AccountFilters enum
func (accF AccountFilter) APIFormOfFilter() string {
	filters := [...]string{"users", "orgs", "admins", "non-admins", "active-users"}

	x := string(accF)
	for _, v := range filters {
		if v == x {
			return x
		}
	}

	return "all"
}

// CreateAccount method - checking the MSR health endpoint
func (c *Client) CreateAccount(ctx context.Context, acc CreateAccount) (ResponseAccount, error) {
	if (acc == CreateAccount{}) {
		return ResponseAccount{}, fmt.Errorf("creating account failed. %w: %+v", ErrEmptyStruct, acc)
	}
	body, err := json.Marshal(acc)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("creating account %s failed. %w: %s", acc.Name, ErrMarshaling, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.createEnziUrl("accounts"), bytes.NewBuffer(body))
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("creating account %s failed. %w: %s", acc.Name, ErrRequestCreation, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("creating account %s failed. %w", acc.Name, err)
	}

	resAcc := ResponseAccount{}
	if err := json.Unmarshal(resBody, &resAcc); err != nil {
		return ResponseAccount{}, fmt.Errorf("creating account %s failed. %w: %s", acc.Name, ErrUnmarshaling, err)
	}

	return resAcc, nil
}

// DeleteAccount deletes a user from in Enzi
func (c *Client) DeleteAccount(ctx context.Context, id string) error {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("deleting account %s failed. %w: %s", id, ErrRequestCreation, err)
	}

	if _, err = c.doRequest(req); err != nil {
		return fmt.Errorf("deleting account %s failed. %w", id, err)
	}
	return nil
}

// ReadAccount method retrieves a user from the enzi endpoint
func (c *Client) ReadAccount(ctx context.Context, id string) (ResponseAccount, error) {
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("reading account %s failed. %w: %s", id, ErrRequestCreation, err)
	}

	body, err := c.doRequest(req)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("reading account %s failed. %w", id, err)
	}

	resAcc := ResponseAccount{}
	if err := json.Unmarshal(body, &resAcc); err != nil {
		return ResponseAccount{}, fmt.Errorf("reading account %s failed. %w: %s", id, ErrUnmarshaling, err)
	}
	return resAcc, nil
}

// UpdateAccount updates a user in the enzi endpoint
func (c *Client) UpdateAccount(ctx context.Context, id string, acc UpdateAccount) (ResponseAccount, error) {
	// if (acc == UpdateAccount{}) {
	// 	return ResponseAccount{}, fmt.Errorf("updating account failed. %w: %+v", ErrEmptyStruct, acc)
	// }
	url := fmt.Sprintf("%s/%s", c.createEnziUrl("accounts"), id)
	body, err := json.Marshal(acc)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("updating account %s failed. %w: %s", id, ErrMarshaling, err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("updating account %s failed. %w: %s", id, ErrRequestCreation, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponseAccount{}, fmt.Errorf("updating account %s failed. %w", id, err)
	}

	resAcc := ResponseAccount{}
	if err := json.Unmarshal(resBody, &resAcc); err != nil {
		return ResponseAccount{}, fmt.Errorf("updating account %s failed. %w: %s", id, ErrUnmarshaling, err)
	}
	return resAcc, nil
}

// ReadAccounts method retrieves all accounts depending on the filter passed from the enzi endpoint
func (c *Client) ReadAccounts(ctx context.Context, accFilter AccountFilter) ([]ResponseAccount, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.createEnziUrl("accounts"), nil)
	if err != nil {
		return []ResponseAccount{}, fmt.Errorf("reading accounts in bulk '%s' failed. %w: %s",
			accFilter.APIFormOfFilter(), ErrRequestCreation, err)
	}

	q := req.URL.Query()
	q.Add("filter", accFilter.APIFormOfFilter())
	req.URL.RawQuery = q.Encode()

	body, err := c.doRequest(req)
	if err != nil {
		return []ResponseAccount{}, fmt.Errorf("reading accounts in bulk '%s' failed. %w",
			accFilter.APIFormOfFilter(), err)
	}

	accs := struct {
		UsersCount    int    `json:"usersCount"`
		OrgsCount     int    `json:"orgsCount"`
		ResourceCount int    `json:"resourceCount"`
		NextPageStart string `json:"nextPageStart"`

		Accounts []ResponseAccount `json:"accounts"`
	}{}

	if err := json.Unmarshal(body, &accs); err != nil {
		return []ResponseAccount{}, fmt.Errorf("reading accounts in bulk '%s' failed. %w: %s",
			accFilter.APIFormOfFilter(), ErrUnmarshaling, err)
	}

	return accs.Accounts, nil
}
