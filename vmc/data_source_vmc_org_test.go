/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVmcOrgBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckZerocloud(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcOrgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_org.my_org", "display_name", os.Getenv(constants.OrgDisplayName)),
				),
			},
		},
	})
}

func testAccDataSourceVmcOrgConfig() string {
	return `data "vmc_org" "my_org" {}`
}
