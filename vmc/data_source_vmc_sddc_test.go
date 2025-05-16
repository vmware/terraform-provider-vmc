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

func TestAccDataSourceVmcSddcBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcSddcConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_sddc.sddc_imported", "sddc_name", os.Getenv(constants.TestSddcName)),
					// TODO: consider adding another env variable for the primary cluster host count
					resource.TestCheckResourceAttr("data.vmc_sddc.sddc_imported", "num_host", "2"),
				),
			},
		},
	})
}

func testAccDataSourceVmcSddcConfig() string {
	return fmt.Sprintf(`
data "vmc_sddc" "sddc_imported" {
  sddc_id = %q
}
`, os.Getenv(constants.TestSddcID),
	)
}
