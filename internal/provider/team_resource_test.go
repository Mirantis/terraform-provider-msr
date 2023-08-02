package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestTeamResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testTeamResourceDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_team.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("msr_team.test", "org_id", TestingVersion),
					resource.TestCheckResourceAttr("msr_team.test", "description", TestingVersion),
				),
			},
			// ImportState testing
			{
				ResourceName:  "msr_team.test",
				ImportStateId: "test,test",
				ImportState:   true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "msr_team" "test" {
				name = "blah"
				org_id = "blah"
				description = "blah"
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_team.test", "name", "blah"),
					resource.TestCheckResourceAttr("msr_team.test", "org_id", "blah"),
					resource.TestCheckResourceAttr("msr_team.test", "description", "blah"),
				),
			},
		},
	})
}

func testTeamResourceDefault() string {
	return `
	resource "msr_team" "test" {
		name = "test"
		org_id = "test"
		description = "test"
	}`
}
