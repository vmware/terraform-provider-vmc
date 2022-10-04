/* Copyright 2021 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVmcSddcBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcSddcConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_sddc.sddc_imported", "sddc_name", os.Getenv(TestSDDCName)),
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
`,
		os.Getenv(TestSDDCId),
	)
}
