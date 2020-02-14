/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcOrg_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcOrgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_org.my_org", "display_name", os.Getenv("ORG_DISPLAY_NAME")),
				),
			},
		},
	})
}

func testAccDataSourceVmcOrgConfig() string {
	return fmt.Sprintf(`
data "vmc_org" "my_org" {
	id = %q

}
`,
		os.Getenv("ORG_ID"),
	)
}
