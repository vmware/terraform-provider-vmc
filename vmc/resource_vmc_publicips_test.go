/* Copyright © 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/model"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs/publicips"
	"net"
	"os"
	"testing"
)

func TestAccResourceVmcPublicIP_basic(t *testing.T) {
	var publicIPResource model.SddcPublicIp
	VMName := "terraform_test_vm_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcPublicIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVmcPublicIPConfigBasic(VMName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcPublicIPExists("vmc_publicips.publicip_1", &publicIPResource),
					testCheckPublicIPAttributes(&publicIPResource),
					resource.TestCheckResourceAttrSet("vmc_publicips.publicip_1", "public_ip"),
				),
			},
		},
	})
}

func testCheckVmcPublicIPExists(name string, publicIPResource *model.SddcPublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := rs.Primary.Attributes["org_id"]
		vmName := rs.Primary.Attributes["name"]
		allocationID := rs.Primary.Attributes["id"]
		if allocationID == "" {
			return fmt.Errorf("allocation ID of the Public IP Resource is not set")
		}
		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		connector := connectorWrapper.Connector
		publicIPClient := publicips.NewPublicipsClientImpl(connector)
		var err error
		*publicIPResource, err = publicIPClient.Get(orgID, sddcID, allocationID)
		if err != nil {
			return fmt.Errorf("Bad: Get on publicIP API: %s", err)
		}

		if *publicIPResource.Name != vmName {
			return fmt.Errorf("Bad: Public IP %q does not exist", allocationID)
		}
		fmt.Printf("Public IP created successfully with id %s \n", allocationID)
		return nil
	}
}

func testCheckPublicIPAttributes(publicIPResource *model.SddcPublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		addr := net.ParseIP(publicIPResource.PublicIp)
		if addr == nil {
			return fmt.Errorf("The alloted Public IP %s is not valid", publicIPResource.PublicIp)
		}
		return nil
	}
}

func testCheckVmcPublicIPDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	publicIPClient := publicips.NewPublicipsClientImpl(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_publicips" {
			continue
		}

		allocationID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		sddcID := rs.Primary.Attributes["sddc_id"]
		vmName := rs.Primary.Attributes["name"]
		publicIPs, err := publicIPClient.List(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Error while getting the public Ips %s", err)
		}
		for _, publicIp := range publicIPs {
			if *(publicIp.AllocationId) == allocationID {
				return fmt.Errorf("Entity PublicIP %s still exits with allocation ID %s", vmName, allocationID)
			}
		}
	}
	return nil
}

func testAccVmcPublicIPConfigBasic(name string) string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
}
	
data "vmc_org" "my_org" {
	id = %q
}

resource "vmc_publicips" "publicip_1" {
	org_id = "${data.vmc_org.my_org.id}"
	sddc_id = %q
	name     = %q
	private_ip = "10.105.167.133"
}
`,
		os.Getenv("REFRESH_TOKEN"),
		os.Getenv("ORG_ID"),
		os.Getenv("TEST_SDDC_ID"),
		name,
	)
}
