/* Copyright 2021 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVmcSddc_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcSddcConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_sddc.sddc_imported", "sddc_name", os.Getenv(TestSDDCName)),
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
