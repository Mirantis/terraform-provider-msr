package provider_test

import (
	"testing"

	provider "github.com/Mirantis/terraform-provider-msr/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = provider.Provider()
}
