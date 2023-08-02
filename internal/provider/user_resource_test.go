package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestUserResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testUserResourceDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_user.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("msr_user.test", "password", TestingVersion),
					resource.TestCheckResourceAttr("msr_user.test", "full_name", TestingVersion),
					resource.TestCheckResourceAttr("msr_user.test", "is_admin", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName: "msr_user.test",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "msr_user" "test" {
				name = "blah"
				password = "blah"
				full_name = "blah"
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_user.test", "name", "blah"),
					resource.TestCheckResourceAttr("msr_user.test", "password", "blah"),
					resource.TestCheckResourceAttr("msr_user.test", "full_name", "blah"),
					resource.TestCheckResourceAttr("msr_user.test", "is_admin", "false"),
				),
			},
			// Delete is called implicitly
		},
	})
}

func testUserResourceDefault() string {
	return `
	resource "msr_user" "test" {
		name = "test"
		password = "test"
		full_name = "test"
	}`
}
