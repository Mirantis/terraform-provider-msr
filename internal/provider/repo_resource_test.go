package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestRepoResourceDefault(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testRepoResourceDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_repo.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("msr_repo.test", "org_name", TestingVersion),
					resource.TestCheckResourceAttr("msr_repo.test", "visibility", "private"),
					resource.TestCheckResourceAttr("msr_repo.test", "scan_on_push", "false"),
					resource.TestCheckResourceAttrSet("msr_repo.test", "id"),
				),
			},
			// Create and Read testing
			{
				Config: providerConfig + testRepoResourceValuesSet(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_repo.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("msr_repo.test", "org_name", TestingVersion),
					resource.TestCheckResourceAttr("msr_repo.test", "visibility", "public"),
					resource.TestCheckResourceAttr("msr_repo.test", "scan_on_push", "true"),
					resource.TestCheckResourceAttrSet("msr_repo.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:  "msr_repo.test",
				ImportStateId: "test,test",
				ImportState:   true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
				resource "msr_repo" "test" {
				name = "blah"
				org_name = "blah"
				visibility = "blah"
				scan_on_push = "false"
			}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_repo.test", "name", "blah"),
					resource.TestCheckResourceAttr("msr_repo.test", "org_name", "blah"),
					resource.TestCheckResourceAttr("msr_repo.test", "visibility", "blah"),
					resource.TestCheckResourceAttr("msr_repo.test", "scan_on_push", "false"),
				),
			},
			// Delete is called implicitly
		},
	})
}

func testRepoResourceDefault() string {
	return `
	resource "msr_repo" "test" {
		name = "test"
		org_name = "test"
	}`
}

func testRepoResourceValuesSet() string {
	return `
	resource "msr_repo" "test" {
		name = "test"
		org_name = "test"
		visibility = "public"
		scan_on_push = "true"
	}`
}
