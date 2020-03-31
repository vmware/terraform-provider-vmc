/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
)

func TestAccResourceVmcSRMNodes_basic(t *testing.T) {
	srmExtensionKeySuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSRMNodesDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSRMNodeConfigBasic(srmExtensionKeySuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSRMNodesExists("vmc_srm_nodes.srm_nodes_1"),
					resource.TestCheckResourceAttrSet("vmc_srm_nodes.srm_nodes_1", "srm_nodes"),
				),
			},
		},
	})
}

func testCheckVmcSRMNodesExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["sddc_id"]
		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		connector := connectorWrapper.Connector
		orgID := connectorWrapper.OrgID

		draasClient := draas.NewDefaultSiteRecoveryClient(connector)
		var err error
		siteRecovery, err := draasClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: Get on DraaS API: %s", err)
		}

		if *siteRecovery.SddcId != sddcID {
			return fmt.Errorf("Bad: Site recovery: %s ", sddcID)
		}

		return nil
	}
}

func testCheckVmcSRMNodesDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	draasClient := draas.NewDefaultSiteRecoveryClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_srm_nodes" {
			continue
		}

		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := connectorWrapper.OrgID
		siteRecovery, err := draasClient.Get(orgID, sddcID)
		if err == nil {
			if *siteRecovery.SiteRecoveryState != "DEACTIVATED" {
				return fmt.Errorf("Site recovery activated for  SDDC with ID : %s ", sddcID)
			}
			return nil
		}
		//check if error type if not_found
		if err.Error() != (errors.NotFound{}.Error()) {
			return err
		}
	}

	return nil
}

func testAccVmcSRMNodeConfigBasic(srmExtensionKeySuffix string) string {
	return fmt.Sprintf(`
resource "vmc_site_recovery" "site_recovery_1" {
  sddc_id = %q
}

resource "vmc_srm_nodes" "srm_node_1"{
  sddc_id = %q
  srm_extension_key_suffix = %q
  depends_on = [vmc_site_recovery.site_recovery_1]
}`,
		os.Getenv(TestSDDCId),
		os.Getenv(TestSDDCId),
		srmExtensionKeySuffix,
	)
}
