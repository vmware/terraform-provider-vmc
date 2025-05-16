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

func TestAccDataSourceVmcCustomerSubnetsBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcCustomerSubnetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.#", "4"),
					// The following subnet IDs are tightly coupled to the AWS account number provided for testing.
					// Since the provisioning of an AWS account incurs costs the hardcoded subnets IDs will have to do.
					// If this tests fails it was probably ran with a different "AwsAccountNumber" than originally designed.
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.0", "subnet-01715c65359792049"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.1", "subnet-01d62fb7a6ef9ca1b"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.2", "subnet-0cd7c7fdd15b08b07"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.3", "subnet-08d5d9dc3aad0383a"),
				),
			},
		},
	})
}

func TestAccDataSourceVmcCustomerSubnetsOnlyRequiredProperties(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcCustomerSubnetsOnlyRequiredProperties(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.#", "4"),
					// The following subnet IDs are tightly coupled to the AWS account number provided for testing.
					// Since the provisioning of an AWS account incurs costs the hardcoded subnets IDs will have to do.
					// If this tests fails it was probably ran with a different "AwsAccountNumber" than originally designed.
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
  region               = "US_WEST_2"
  sddc_type            = "SingleAZ"
  instance_type        = "i3.metal"
}
`, os.Getenv(constants.AwsAccountNumber))
}

func testAccDataSourceVmcCustomerSubnetsOnlyRequiredProperties() string {
	return fmt.Sprintf(`


data "vmc_connected_accounts" "my_accounts" {
  account_number = %q
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = "US_WEST_2"
}
`, os.Getenv(constants.AwsAccountNumber))
}
