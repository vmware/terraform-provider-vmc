/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"testing"

	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas/model"
)

func TestAccResourceVmcSiteRecoveryZerocloud(t *testing.T) {
	var siteRecovery model.SiteRecovery
	resourceName := "vmc_site_recovery.site_recovery_1"
	srmExtensionKeySuffix := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSiteRecoveryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSiteConfigBasic(srmExtensionKeySuffix),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSiteRecoveryExists("vmc_site_recovery.site_recovery_1", &siteRecovery),
					resource.TestCheckResourceAttr("vmc_site_recovery.site_recovery_1",
						"site_recovery_state", model.SiteRecovery_SITE_RECOVERY_STATE_ACTIVATED),
				),
			},
			{
				ResourceName:            resourceName,
				ImportStateIdFunc:       testAccVmcSiteRecoveryResourceImportStateIDFunc(resourceName),
				ImportStateVerifyIgnore: []string{"sddc_id", "srm_extension_key_suffix"},
				ImportState:             true,
				ImportStateVerify:       true,
			},
		},
	})
}

func testCheckVmcSiteRecoveryExists(name string, siteRecovery *model.SiteRecovery) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["sddc_id"]
		connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
		orgID := connectorWrapper.OrgID

		draasClient := draas.NewSiteRecoveryClient(connectorWrapper)
		var err error
		*siteRecovery, err = draasClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("error retrieving site recovery information for SDDC %s : %s", sddcID, err)
		}

		if *siteRecovery.SddcId != sddcID {
			return fmt.Errorf("error retrieving site recovery for SDDC with id %s ", sddcID)
		}

		return nil
	}
}

func testCheckVmcSiteRecoveryDestroy(s *terraform.State) error {
	connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
	draasClient := draas.NewSiteRecoveryClient(connectorWrapper)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_site_recovery" {
			continue
		}

		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := connectorWrapper.OrgID
		siteRecovery, err := draasClient.Get(orgID, sddcID)
		if err == nil {
			if *siteRecovery.SiteRecoveryState != model.SiteRecovery_SITE_RECOVERY_STATE_DEACTIVATED &&
				*siteRecovery.SiteRecoveryState != model.SiteRecovery_SITE_RECOVERY_STATE_DELETED {
				return fmt.Errorf("Site recovery activated for  SDDC with ID : %s ", sddcID)
			}
			return nil
		}
		// check if error type is not_found
		if err.Error() != (errors.NotFound{}.Error()) {
			return err
		}
	}

	return nil
}

func testAccVmcSiteConfigBasic(srmExtensionKeySuffix string) string {
	return fmt.Sprintf(`
resource "vmc_sddc" "srm_test_sddc" {
	sddc_name           = "terraform_srm_test"
	num_host            = 2
	provider_type       = "ZEROCLOUD"
	host_instance_type = "I3_METAL"
	region = "US_WEST_2"
	delay_account_link  = true
}

	resource "vmc_site_recovery" "site_recovery_1" {
		sddc_id = vmc_sddc.srm_test_sddc.id
	srm_extension_key_suffix = %q
	}`,
		srmExtensionKeySuffix,
	)
}

func testAccVmcSiteRecoveryResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return rs.Primary.Attributes["sddc_id"], nil
	}
}
