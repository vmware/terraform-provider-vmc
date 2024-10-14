// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/nsx_vmc_app/infra"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
)

func TestAccResourceVmcPublicIp_basic(t *testing.T) {
	displayName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "vmc_public_ip.public_ip_1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcPublicIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcPublicIPConfigBasic(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVmcPublicIPExists("vmc_public_ip.public_ip_1"),
					resource.TestCheckResourceAttrSet("vmc_public_ip.public_ip_1", "display_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccVmcPublicIPResourceImportStateIDFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVmcPublicIPExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		uuid := rs.Primary.Attributes["id"]
		displayName := rs.Primary.Attributes["display_name"]
		connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
		nsxConnector, err := getNsxtReverseProxyURLConnector(os.Getenv(constants.NsxtReverseProxyURL), connectorWrapper)
		if err != nil {
			return fmt.Errorf("error creating client nsxConnector : %v ", err)
		}

		publicIpsClient := infra.NewPublicIpsClient(nsxConnector)
		publicIP, err := publicIpsClient.Get(uuid)
		if err != nil {
			return fmt.Errorf("error getting public IP with ID %s : %v", uuid, err)
		}

		if *publicIP.Id != uuid {
			return fmt.Errorf("error public IP %q does not exist", displayName)
		}
		return nil
	}
}

func testCheckVmcPublicIPDestroy(s *terraform.State) error {
	fmt.Printf("Reverse proxy : %s", os.Getenv(constants.NsxtReverseProxyURL))
	connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
	nsxConnector, err := getNsxtReverseProxyURLConnector(os.Getenv(constants.NsxtReverseProxyURL), connectorWrapper)
	if err != nil {
		return fmt.Errorf("error creating client nsxConnector : %v ", err)
	}
	publicIpsClient := infra.NewPublicIpsClient(nsxConnector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_public_ip" {
			continue
		}

		uuid := rs.Primary.Attributes["id"]
		fmt.Printf("UUID : %s ", uuid)
		publicIP, err := publicIpsClient.Get(uuid)
		fmt.Printf("publicIP : %v", publicIP.Id)
		if err == nil {
			if *publicIP.Id == uuid {
				return fmt.Errorf("public IP %s with ID %s still exits", *publicIP.DisplayName, uuid)
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

func testAccVmcPublicIPConfigBasic(displayName string) string {
	return fmt.Sprintf(`
resource "vmc_public_ip" "public_ip_1" {
	display_name = %q
	nsxt_reverse_proxy_url = %q

}
`,
		displayName,
		os.Getenv(constants.NsxtReverseProxyURL),
	)
}

func testAccVmcPublicIPResourceImportStateIDFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["nsxt_reverse_proxy_url"]), nil
	}
}
