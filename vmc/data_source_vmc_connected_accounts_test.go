/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

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
					resource.TestCheckResourceAttr("data.vmc_connected_accounts.my_accounts", "ids.#", "1"),
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
	id = %q
}
	
data "vmc_connected_accounts" "my_accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}
`,
		os.Getenv("REFRESH_TOKEN"),
		os.Getenv("ORG_ID"),
	)
}
