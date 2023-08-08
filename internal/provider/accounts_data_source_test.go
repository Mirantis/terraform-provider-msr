package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccountsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testMSRaccountsDefault(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.msr_accounts.test", "id"),
				),
			},
			{
				Config: providerConfig + testMSRaccountsSetValues(),
				Check: resource.ComposeAggregateTestCheckFunc(
					// first account
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.name_or_id", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.full_name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.is_org", "true"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.is_active", "true"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.is_admin", "true"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.is_imported", "true"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.otp_enabled", "true"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.members_count", "1"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.0.teams_count", "1"),
					// second account
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.name_or_id", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.full_name", TestingVersion),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.is_org", "false"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.is_active", "false"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.is_admin", "false"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.is_imported", "false"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.otp_enabled", "false"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.members_count", "5"),
					resource.TestCheckResourceAttr("data.msr_accounts.test", "accounts.1.teams_count", "5"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttrSet("data.msr_accounts.test", "id"),
				),
			},
		},
	})
}

func testMSRaccountsDefault() string {
	return `
	data "msr_accounts" "test" {
		filter = "all"
	}
	`
}

func testMSRaccountsSetValues() string {
	return `
	data "msr_accounts" "test" {
		filter = "all"

		accounts {
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

		accounts {
			name_or_id = "test"
			name = "test"
			full_name = "test"
			is_org = "false"
			is_active = "false"
			is_admin = "false"
			is_imported = "false"
			on_demand = "false"
			otp_enabled = "false"
			members_count = 5
			teams_count = 5
		}
	}
	`
}
