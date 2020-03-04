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
		},
	})
}

func testAccCheckVmcPublicIpExists(name string, publicIpResource *model.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		uuid := rs.Primary.Attributes["id"]
		displayName := rs.Primary.Attributes["display_name"]
		connector, err := getNsxApiPublicEndpointConnector(os.Getenv(NSXApiPublicEndpointUrl))
		if err != nil {
			return fmt.Errorf("Bad: creating connector : %v ", err)
		}

		nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return fmt.Errorf("Bad: getting public IP with ID %s : %v", uuid, err)
		}

		if *publicIp.Id != uuid {
			return fmt.Errorf("Bad: Public IP %q does not exist", displayName)
		}
		return nil
	}
}

func testCheckVmcPublicIpDestroy(s *terraform.State) error {
	connector, err := getNsxApiPublicEndpointConnector(os.Getenv(NSXApiPublicEndpointUrl))
	if err != nil {
		return fmt.Errorf("Bad: creating connector : %v ", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_public_ip" {
			continue
		}

		uuid := rs.Primary.Attributes["id"]
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err == nil {
			if *publicIp.Id == uuid {
				return fmt.Errorf("Public IP %s with ID %s still exits", *publicIp.DisplayName, uuid)
			}
			return nil
		}
		// check if error type if not_found
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
	nsx_api_public_endpoint_url = %q

}
`,
		displayName,
		os.Getenv(NSXApiPublicEndpointUrl),
	)
}
