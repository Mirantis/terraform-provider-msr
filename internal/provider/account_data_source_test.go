package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccountDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testMSRaccountDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the item to ensure all attributes are set
					resource.TestCheckResourceAttr("data.msr_account.test", "name_or_id", TestingVersion),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.msr_account.test", "id"),
				),
			},
			{
				Config: providerConfig + testMSRaccountSetValues(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify the item to ensure all attributes are set
					resource.TestCheckResourceAttr("data.msr_account.test", "name_or_id", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_account.test", "name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_account.test", "full_name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_account.test", "is_org", "true"),
					resource.TestCheckResourceAttr("data.msr_account.test", "is_active", "true"),
					resource.TestCheckResourceAttr("data.msr_account.test", "is_admin", "true"),
					resource.TestCheckResourceAttr("data.msr_account.test", "is_imported", "true"),
					resource.TestCheckResourceAttr("data.msr_account.test", "otp_enabled", "true"),
					resource.TestCheckResourceAttr("data.msr_account.test", "members_count", "1"),
					resource.TestCheckResourceAttr("data.msr_account.test", "teams_count", "1"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.msr_account.test", "id"),
				),
			},
		},
	})
}

func testMSRaccountDefault() string {
	return `
	data "msr_account" "test" {
		name_or_id = "test"
	}
	`
}

func testMSRaccountSetValues() string {
	return `
	data "msr_account" "test" {
		name_or_id = "test"
		name = "test"
		full_name = "test"
		is_org = "true"
		is_active = "true"
		is_admin = "true"
		is_imported = "true"
		on_demand = "true"
		otp_enabled = "true"
		members_count = 1
		teams_count = 1
	}
	`
}
