/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVmcOrg_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcOrgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_org.my_org", "display_name", os.Getenv(OrgDisplayName)),
				),
			},
		},
	})
}

func testAccDataSourceVmcOrgConfig() string {
	return fmt.Sprintf(`
data "vmc_org" "my_org" {
}
`,
	)
}
