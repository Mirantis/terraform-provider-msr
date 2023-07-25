package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestOrgResourceNoID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testOrgResource_noID(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("msr_org", "name", "test"),
				),
			},
		},
	})
}

func testOrgResource_noID() string {
	return `
	provider "msr" {
		host = "test"
		username = "test"
		password = "test"
	}
	resource "msr_org" "test" {
		name = "test"
	}`
}

// func TestOrgResourceWithID(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: testOrgResource_withID(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("msr_org", "name", "test"),
// 					resource.TestCheckResourceAttr("msr_org", "id", "test"),
// 				),
// 			},
// 		},
// 	})
// }

// func testOrgResource_withID() string {
// 	return `
// 	resource "msr_org" "test" {
// 		name = "test"
// 		id = "test"
// 	}`
// }
