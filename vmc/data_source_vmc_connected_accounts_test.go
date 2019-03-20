package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcConnectedAccounts_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcConnectedAccountsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.accounts", "ids.#", "1"),
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.accounts", "ids.0", "2d5259d2-e85d-3d02-a511-9045094b4b10"),
				),
			},
		},
	})
}

func testAccDataSourceVmcConnectedAccountsConfig() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
}
	
data "vmc_org" "my_org" {
	id = "058f47c4-92aa-417f-8747-87f3ed61cb45"
}
	
data "vmc_connected_accounts" "accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}
`,
		os.Getenv("REFRESH_TOKEN"),
	)
}
