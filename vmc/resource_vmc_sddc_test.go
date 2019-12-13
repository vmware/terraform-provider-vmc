/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"os"
	"testing"
)

func TestAccResourceVmcSddc_basic(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigBasic(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_1", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_1", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "cloud_password"),
				),
			},
		},
	})
}

func testCheckVmcSddcExists(name string, sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["id"]
		sddcName := rs.Primary.Attributes["sddc_name"]
		orgID := rs.Primary.Attributes["org_id"]
		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		connector := connectorWrapper.Connector

		sddcClient := orgs.NewDefaultSddcsClient(connector)
		var err error
		*sddcResource, err = sddcClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: Get on sddcApi: %s", err)
		}

		if sddcResource.Id != sddcID {
			return fmt.Errorf("Bad: Sddc %q does not exist", sddcName)
		}

		fmt.Printf("SDDC %s created successfully with id %s \n", sddcName, sddcID)
		return nil
	}
}

func testCheckSddcAttributes(sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		sddcState := sddcResource.SddcState
		if *sddcState != "READY" {
			return fmt.Errorf("The SDDC %s with ID %s is not in ready state", *sddcResource.Name, sddcResource.Id)
		}
		return nil
	}
}

func testCheckVmcSddcDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	sddcClient := orgs.NewDefaultSddcsClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		sddc, err := sddcClient.Get(orgID, sddcID)
		if err == nil {
			if *sddc.SddcState != "DELETED" {
				return fmt.Errorf("SDDC %s with ID %s still exits", *sddc.Name, sddc.Id)
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

func testAccVmcSddcConfigBasic(sddcName string) string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
}
	
data "vmc_org" "my_org" {
	id = %q
}

data "vmc_connected_accounts" "accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}

resource "vmc_sddc" "sddc_1" {
	org_id = "${data.vmc_org.my_org.id}"

	# storage_capacity    = 100
	sddc_name = %q

	vpc_cidr      = "10.2.0.0/16"
	num_host      = 1
	provider_type = "ZEROCLOUD"

	region = "US_WEST_2"

	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"

	# TODO raise exception here need to debug
	#account_link_sddc_config = [
	#	{
	#	  customer_subnet_ids  = ["subnet-13a0c249"]
	#	  connected_account_id = "${data.vmc_connected_accounts.accounts.ids.0}"
	#	},
	#  ]
}
`,
		os.Getenv("REFRESH_TOKEN"),
		os.Getenv("ORG_ID"),
		sddcName,
	)
}
