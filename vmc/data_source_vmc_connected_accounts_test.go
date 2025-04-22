// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/vmware/terraform-provider-vmc/vmc/constants"
)

func TestAccDataSourceVmcConnectedAccountsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
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
		os.Getenv(constants.AwsAccountNumber),
	)
}
