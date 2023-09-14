package client

import (
	"context"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GeneratePass creates a random password.
func GeneratePass() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789!@#$")
	rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// PruningPolicyRulesToAPI converts slice of PruningPolicyRuleTFSDK rules to PruningPolicyRuleAPI.
func PruningPolicyRulesToAPI(ctx context.Context, rules []PruningPolicyRuleTFSDK) []PruningPolicyRuleAPI {
	// Convert the TFSDK rules to API rules
	var apiRules []PruningPolicyRuleAPI

	for _, r := range rules {
		apiRules = append(apiRules, r.PruningPolicyRuleToAPI(ctx))
	}
	return apiRules
}

// PruningPolicyRulesToTFSDK converts slice of PruningPolicyRuleAPI rules to PruningPolicyRuleTFSDK.
func PruningPolicyRulesToTFSDK(ctx context.Context, rules []PruningPolicyRuleAPI) []PruningPolicyRuleTFSDK {
	// Convert the API rules to TFSDK rules
	var apiRules []PruningPolicyRuleTFSDK

	for _, r := range rules {
		apiRules = append(apiRules, r.PruningPolicyRuleToTFSDK(ctx))
	}
	return apiRules
}

// PruningPolicyRuleToAPI converts a single PruningPolicyRuleTFSDK rule to PruningPolicyRuleAPI.
func (r *PruningPolicyRuleTFSDK) PruningPolicyRuleToAPI(ctx context.Context) PruningPolicyRuleAPI {
	var values []string
	for _, v := range r.Values {
		values = append(values, v.ValueString())
	}
	return PruningPolicyRuleAPI{
		Field:    r.Field.ValueString(),
		Operator: r.Operator.ValueString(),
		Values:   values,
	}
}

// PruningPolicyRuleToTFSDK converts a single PruningPolicyRuleAPI rule to PruningPolicyRuleToTFSDK.
func (r *PruningPolicyRuleAPI) PruningPolicyRuleToTFSDK(ctx context.Context) PruningPolicyRuleTFSDK {
	var values []types.String
	for _, v := range r.Values {
		values = append(values, basetypes.NewStringValue(v))
	}
	return PruningPolicyRuleTFSDK{
		Field:    basetypes.NewStringValue(r.Field),
		Operator: basetypes.NewStringValue(r.Operator),
		Values:   values,
	}
}
