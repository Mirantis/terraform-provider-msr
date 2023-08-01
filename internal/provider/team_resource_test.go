package provider

// import (
// 	"testing"

// 	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
// )

// func TestTeamResourceDefault(t *testing.T) {
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
// 		Steps: []resource.TestStep{
// 			// Create and Read testing
// 			{
// 				Config: providerConfig + testTeamResourceDefault(),
// 				Check: resource.ComposeAggregateTestCheckFunc(
// 					resource.TestCheckResourceAttr("msr_team.blah", "name", TestingVersion),
// 					// resource.TestCheckResourceAttr("msr_team.test", "org_id", TestingVersion),
// 					// resource.TestCheckResourceAttr("msr_team.test", "description", TestingVersion),
// 				),
// 			},
// 		},
// 	})
// }

// func testTeamResourceDefault() string {
// 	return `
// 	resource "msr_team" "blah" {
// 		name = "test"
// 		org_id = "test"
// 		description = "test"
// 	}`
// }
