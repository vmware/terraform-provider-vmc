/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/api"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/model"
)

func TestAccResourceVmcPublicIp_basic(t *testing.T) {
	var publicIpResource model.PublicIp
	displayName := "test-public-ip"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcPublicIpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcPublicIpConfigBasic(displayName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcPublicIpExists("vmc_publicip", &publicIpResource),
					resource.TestCheckResourceAttrSet("vmc_publicip.public_ip_1", "display_name"),
				),
			},
		},
	})
}

func testCheckVmcPublicIpExists(name string, publicIpResource *model.PublicIp) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		uuid := rs.Primary.Attributes["id"]
		displayName := rs.Primary.Attributes["display_name"]
		connector, err := getNSXTConnector()
		if err != nil {
			return fmt.Errorf("Bad: creating connector : %v ", err)
		}

		nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return fmt.Errorf("Bad: getting public IP with ID %s : %v", uuid, err)
		}

		if *publicIp.Id != uuid {
			return fmt.Errorf("Bad: Public Ip %q does not exist", displayName)
		}
		return nil
	}
}

func testCheckVmcPublicIpDestroy(s *terraform.State) error {
	connector, err := getNSXTConnector()
	if err != nil {
		return fmt.Errorf("Bad: creating connector : %v ", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_publicip" {
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
resource "vmc_publicip" "public_ip_1" {
	display_name = %q

}
`,
		displayName,
	)
}

func getNSXTConnector() (client.Connector, error) {
	apiToken := os.Getenv(APIToken)
	nsxtReverseProxyURL := os.Getenv(NSXTReverseProxyUrl)
	if len(nsxtReverseProxyURL) == 0 {
		return nil, fmt.Errorf(NSXTReverseProxyUrl + " is a required parameter.")
	}
	if strings.Contains(nsxtReverseProxyURL, SksNSXTManager) {
		nsxtReverseProxyURL = strings.Replace(nsxtReverseProxyURL, SksNSXTManager, "", -1)
	}
	httpClient := http.Client{}
	connector, err := NewClientConnectorByRefreshToken(apiToken, nsxtReverseProxyURL, DefaultCSPUrl, httpClient)
	if err != nil {
		return nil, fmt.Errorf("Error creating connector : %v ", err)
	}
	return connector, nil
}
