/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVmcConnectedAccounts_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcConnectedAccountsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.vmc_connected_accounts.my_accounts", "id"),
				),
			},
		},
	})
}

func testAccDataSourceVmcConnectedAccountsConfig() string {
	return fmt.Sprintf(`
data "vmc_connected_accounts" "my_accounts" {
	account_number = %q
}
`,
		os.Getenv(AWSAccountNumber),
	)
}
