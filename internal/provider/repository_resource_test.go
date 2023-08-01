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
