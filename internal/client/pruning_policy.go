package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PruningPolicyRuleAPI struct {
	Field    string   `tfsdk:"field" json:"field"`
	Operator string   `tfsdk:"operator" json:"operator"`
	Values   []string `tfsdk:"values" json:"values"`
}

type PruningPolicyRuleTFSDK struct {
	Field    types.String   `tfsdk:"field" json:"field"`
	Operator types.String   `tfsdk:"operator" json:"operator"`
	Values   []types.String `tfsdk:"values" json:"values"`
}

type CreatePruningPolicy struct {
	Enabled bool                   `json:"enabled"`
	Rules   []PruningPolicyRuleAPI `json:"rules"`
}

type ResponsePruningPolicy struct {
	ID      string                 `json:"id"`
	Enabled bool                   `json:"enabled"`
	Rules   []PruningPolicyRuleAPI `json:"rules"`
}

// CreatePruningPolicy creates a pruning policy in MSR.
func (c *Client) CreatePruningPolicy(ctx context.Context, orgName string, repoName string, policy CreatePruningPolicy) (ResponsePruningPolicy, error) {
	body, err := json.Marshal(policy)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("creating pruning policy %+v failed. %w: %s", policy, ErrMarshaling, err)
	}
	url := fmt.Sprintf("%s/%s/%s/pruningPolicies?initialEvaluation=true", c.createMsrUrl("repositories"), orgName, repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("creating pruning policy %+v failed. %w: %s", policy, ErrRequestCreation, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("creating pruning policy %+v failed. %w", policy, err)
	}

	resPolicy := ResponsePruningPolicy{}
	if err := json.Unmarshal(resBody, &resPolicy); err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("creating pruning policy %+v failed. %w: %s", policy, ErrUnmarshaling, err)
	}

	return resPolicy, nil
}

// ReadPruningPolicy reads specific pruning policy of a repo in MSR.
func (c *Client) ReadPruningPolicy(ctx context.Context, orgName string, repoName string, policyId string) (ResponsePruningPolicy, error) {
	url := fmt.Sprintf("%s/%s/%s/pruningPolicies/%s", c.createMsrUrl("repositories"), orgName, repoName, policyId)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("reading pruning policy for %s/%s failed. %w: %s", orgName, repoName, ErrRequestCreation, err)
	}
	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("reading pruning policy for %s/%s failed. %w", orgName, repoName, err)
	}

	resPolicy := ResponsePruningPolicy{}
	if err := json.Unmarshal(resBody, &resPolicy); err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("reading pruning policy for %s/%s failed. %w: %s", orgName, repoName, ErrUnmarshaling, err)
	}

	return resPolicy, nil
}

// ReadPruningPolicy reads the pruning policies of a repo in MSR.
func (c *Client) ReadPruningPolicies(ctx context.Context, orgName string, repoName string) ([]ResponsePruningPolicy, error) {
	url := fmt.Sprintf("%s/%s/%s/pruningPolicies", c.createMsrUrl("repositories"), orgName, repoName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []ResponsePruningPolicy{}, fmt.Errorf("reading pruning policies for %s/%s failed. %w: %s", orgName, repoName, ErrRequestCreation, err)
	}
	resBody, err := c.doRequest(req)
	if err != nil {
		return []ResponsePruningPolicy{}, fmt.Errorf("reading pruning policies for %s/%s failed. %w", orgName, repoName, err)
	}

	resPolicies := []ResponsePruningPolicy{}
	if err := json.Unmarshal(resBody, &resPolicies); err != nil {
		return []ResponsePruningPolicy{}, fmt.Errorf("reading pruning policies for %s/%s failed. %w: %s", orgName, repoName, ErrUnmarshaling, err)
	}

	return resPolicies, nil
}

// DeletePruningPolicy deletes a pruning policy in MSR.
func (c *Client) DeletePruningPolicy(ctx context.Context, orgName string, repoName string, policyId string) error {
	url := fmt.Sprintf("%s/%s/%s/pruningPolicies/%s", c.createMsrUrl("repositories"), orgName, repoName, policyId)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("deleting pruning policy for %s/%s failed. %w: %s", orgName, repoName, ErrRequestCreation, err)
	}
	if _, err = c.doRequest(req); err != nil {
		return fmt.Errorf("deleting pruning policy for %s/%s failed. %w", orgName, repoName, err)
	}

	return err
}

// UpdatePruningPolicy updates a pruning policy in MSR.
func (c *Client) UpdatePruningPolicy(ctx context.Context, orgName string, repoName string, policy CreatePruningPolicy, policyId string) (ResponsePruningPolicy, error) {
	body, err := json.Marshal(policy)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("creating pruning policy %+v failed. %w: %s", policy, ErrMarshaling, err)
	}
	url := fmt.Sprintf("%s/%s/%s/pruningPolicies/%s?initialEvaluation=true", c.createMsrUrl("repositories"), orgName, repoName, policyId)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("updating pruning policy for %s/%s failed. %w: %s", orgName, repoName, ErrRequestCreation, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resBody, err := c.doRequest(req)
	if err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("updating pruning policy for %s/%s failed. %w", orgName, repoName, err)
	}

	resPolicy := ResponsePruningPolicy{}
	if err := json.Unmarshal(resBody, &resPolicy); err != nil {
		return ResponsePruningPolicy{}, fmt.Errorf("updating pruning policy for %s/%s failed. %w: %s", orgName, repoName, ErrUnmarshaling, err)
	}

	return resPolicy, nil
}

// PruningPolicyExists compares if a pruning policy exists within a existing pruning policies
func (c *Client) PruningPolicyExists(ctx context.Context, newPolicy CreatePruningPolicy, existingPolicies []ResponsePruningPolicy) ResponsePruningPolicy {

	for _, policy := range existingPolicies {

		if len(policy.Rules) != len(newPolicy.Rules) {
			return ResponsePruningPolicy{}
		}

		var ruleMatch []bool
		for _, existRule := range policy.Rules {
			for _, newRule := range newPolicy.Rules {
				// If the rule has Field and Operator matching compare Values
				if existRule.Field == newRule.Field && existRule.Operator == newRule.Operator {
					for _, existValue := range existRule.Values {
						for _, newValue := range newRule.Values {
							if existValue == newValue {
								ruleMatch = append(ruleMatch, true)
							}
						}
					}
				}
			}
		}
		// we have a policy match
		if len(ruleMatch) == len(policy.Rules) {
			return policy
		}
	}

	// there is no matching existing policy in the MSR endpoint
	return ResponsePruningPolicy{}
}
