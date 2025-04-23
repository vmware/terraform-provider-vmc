// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"github.com/vmware/terraform-provider-vmc/vmc/sddcgroup"
)

func TestAccResourceSddcGroupZerocloud(t *testing.T) {
	randStrTokenForTestRun := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	sddcGroupName := "terraform_test_sddc_group_" + randStrTokenForTestRun
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSddcGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcGroupConfigZerocloud(sddcGroupName),
				Check: resource.ComposeTestCheckFunc(
					testSddcGroupExists,
					resource.TestCheckResourceAttrSet("vmc_sddc_group.sddc_group", "id"),
				),
			},
			{
				ResourceName:      "vmc_sddc_group.sddc_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testCheckSddcGroupDestroyed(s *terraform.State) error {
	if sddcGroupExists(s) {
		return fmt.Errorf("sddc group still exists")
	}
	return nil
}

func testSddcGroupExists(s *terraform.State) error {
	if sddcGroupExists(s) {
		return nil
	}
	return fmt.Errorf("sddc group does not exist")
}

func sddcGroupExists(s *terraform.State) bool {
	connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
	sddcGroupClient := sddcgroup.NewSddcGroupClient(*connectorWrapper)
	err := sddcGroupClient.Authenticate()
	if err != nil {
		return false
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc_group" {
			continue
		}
		sddcGroupID := rs.Primary.ID
		sddcGroup, _, err := sddcGroupClient.GetSddcGroup(sddcGroupID)
		if sddcGroup.Deleted == false && err == nil {
			return true
		}
	}
	return false
}

func testAccVmcSddcGroupConfigZerocloud(sddcGroupName string) string {
	return fmt.Sprintf(`
resource "vmc_sddc_group" "sddc_group" {
	name = %q
	description = "bla bla"
	sddc_member_ids = [ %q, %q ]
}
`,
		sddcGroupName,
		os.Getenv(constants.SddcGroupTestSddc1Id),
		os.Getenv(constants.SddcGroupTestSddc2Id),
	)
}
