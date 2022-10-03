/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcCustomerSubnetsBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcCustomerSubnetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.#", "4"),
					// The following subnet IDs are tightly coupled to the AWS account number provided for testing.
					// Since the provisioning of an AWS account incurs costs the hardcoded subnets IDs will have to do.
					// If this tests fails it was probably ran with a different "AWSAccountNumber" than originally designed.
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.0", "subnet-01715c65359792049"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.1", "subnet-01d62fb7a6ef9ca1b"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.2", "subnet-0cd7c7fdd15b08b07"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.3", "subnet-08d5d9dc3aad0383a"),
				),
			},
		},
	})
}

func testAccDataSourceVmcCustomerSubnetsConfig() string {
	return fmt.Sprintf(`
	
data "vmc_connected_accounts" "my_accounts" {
    account_number = %q
}

data "vmc_customer_subnets" "my_subnets" {
	connected_account_id = data.vmc_connected_accounts.my_accounts.id
	region = "US_WEST_2"
}
`,
		os.Getenv(AWSAccountNumber))
}
