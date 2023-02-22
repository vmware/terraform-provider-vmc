/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/sddcgroup"
	"testing"
)

func TestAccResourceSddcGroupZerocloud(t *testing.T) {
	randStrTokenForTestRun := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	sddc1Name := "terraform_sddc_group_1_test_" + randStrTokenForTestRun
	sddc2Name := "terraform_sddc_group_2_test_" + randStrTokenForTestRun
	sddcGroupName := "terraform_test_sddc_group_" + randStrTokenForTestRun
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckSddcGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcGroupConfigZerocloud(sddc1Name, sddc2Name, sddcGroupName),
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

func testAccVmcSddcGroupConfigZerocloud(sddc1Name string, sddc2Name string, sddcGroupName string) string {
	return fmt.Sprintf(`
resource "vmc_sddc" "sddc_zerocloud_1" {
	sddc_name = %q
	vpc_cidr      = "10.11.0.0/16"
	num_host      = 1
	provider_type = "ZEROCLOUD"
	host_instance_type = "I3_METAL"
	region = "US_WEST_2"
	vxlan_subnet = "192.168.1.0/24"
	delay_account_link  = true
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"
	deployment_type = "SingleAZ"
}

resource "vmc_sddc" "sddc_zerocloud_2" {
	sddc_name = %q
	vpc_cidr      = "10.12.0.0/16"
	num_host      = 1
	provider_type = "ZEROCLOUD"
	host_instance_type = "I3_METAL"
	region = "US_WEST_2"
	vxlan_subnet = "192.168.1.0/24"
	delay_account_link  = true
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"
	deployment_type = "SingleAZ"
}

resource "vmc_sddc_group" "sddc_group" {
	name = %q
	description = "bla bla"
	sddc_member_ids = [ vmc_sddc.sddc_zerocloud_1.id, vmc_sddc.sddc_zerocloud_2.id ]
}
`,
		sddc1Name,
		sddc2Name,
		sddcGroupName,
	)
}
