package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOrgResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + testOrgResource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_org.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("msr_org.test", "id", TestingVersion),
				),
			},
			// ImportState testing
			{
				ResourceName:      "msr_org.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// No Update for the org resource
			// Delete is called implicitly
		},
	})
}

func testOrgResource() string {
	return `
	resource "msr_org" "test" {
		name = "test"
	}`
}
