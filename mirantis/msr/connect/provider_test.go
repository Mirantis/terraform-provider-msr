package connect_test

import (
	"testing"

	connect "github.com/Mirantis/terraform-provider-msr/mirantis/msr/connect"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestProvider(t *testing.T) {
	if err := connect.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = connect.Provider()
}
