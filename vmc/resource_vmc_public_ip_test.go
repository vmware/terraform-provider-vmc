/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/api"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/model"
)

func TestAccResourceVmcPublicIp_basic(t *testing.T) {
	var publicIpResource model.PublicIp
	displayName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resourceName := "vmc_public_ip.public_ip_1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcPublicIpConfigBasic(displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVmcPublicIpExists("vmc_public_ip.public_ip_1", &publicIpResource),
					resource.TestCheckResourceAttrSet("vmc_public_ip.public_ip_1", "display_name"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccVmcPublicIPResourceImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckVmcPublicIpExists(name string, publicIpResource *model.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		uuid := rs.Primary.Attributes["id"]
		displayName := rs.Primary.Attributes["display_name"]
		connector, err := getNSXTReverseProxyUrlConnector(os.Getenv(NSXTReverseProxyUrl))
		if err != nil {
			return fmt.Errorf("error creating client connector : %v ", err)
		}

		nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return fmt.Errorf("error getting public IP with ID %s : %v", uuid, err)
		}

		if *publicIp.Id != uuid {
			return fmt.Errorf("error public IP %q does not exist", displayName)
		}
		return nil
	}
}

func testCheckVmcPublicIpDestroy(s *terraform.State) error {
	fmt.Printf("Reverse proxy : %s", os.Getenv(NSXTReverseProxyUrl))
	connector, err := getNSXTReverseProxyUrlConnector(os.Getenv(NSXTReverseProxyUrl))
	if err != nil {
		return fmt.Errorf("error creating client connector : %v ", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_public_ip" {
			continue
		}

		uuid := rs.Primary.Attributes["id"]
		fmt.Printf("UUID : %s ", uuid)
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		fmt.Printf("publicIP : %v", publicIp.Id)
		if err == nil {
			if *publicIp.Id == uuid {
				return fmt.Errorf("public IP %s with ID %s still exits", *publicIp.DisplayName, uuid)
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

func testAccVmcPublicIpConfigBasic(displayName string) string {
	return fmt.Sprintf(`
resource "vmc_public_ip" "public_ip_1" {
	display_name = %q
	nsxt_reverse_proxy_url = %q

}
`,
		displayName,
		os.Getenv(NSXTReverseProxyUrl),
	)
}

func testAccVmcPublicIPResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["nsxt_reverse_proxy_url"]), nil
	}
}
