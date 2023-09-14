package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPruningPolicyResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testPruningPolicyResourceDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "enabled", "true"),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "org_name", TestingVersion),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "repo_name", TestingVersion),
					// first rule
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.field", TestingVersion),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.operator", TestingVersion),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.values.0", TestingVersion),
					// second rule
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.1.field", TestingVersion),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.1.operator", TestingVersion),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.1.values.0", TestingVersion),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("msr_pruning_policy.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:  "msr_pruning_policy.test",
				ImportState:   true,
				ImportStateId: "test,test,test",
			},
			// Update and Read testing
			{
				Config: providerConfig + `
					resource "msr_pruning_policy" "test" {
						enabled = "false"
						org_name = "blah"
						repo_name = "blah"
						rule {
							field = "blah"
							operator = "blah"
							values = ["blah"]
						}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "enabled", "false"),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "org_name", "blah"),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "repo_name", "blah"),
					// first rule
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.field", "blah"),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.operator", "blah"),
					resource.TestCheckResourceAttr("msr_pruning_policy.test", "rule.0.values.0", "blah"),
				),
			},
			// Delete is called implicitly
		},
	})
}

func testPruningPolicyResourceDefault() string {
	return `
	resource "msr_pruning_policy" "test" {
		enabled = "true"
		org_name = "test"
		repo_name = "test"
		rule {
			field = "test"
			operator = "test"
			values = ["test"]
		}
		rule {
			field = "test"
			operator = "test"
			values = ["test"]
		}
	}`
}
