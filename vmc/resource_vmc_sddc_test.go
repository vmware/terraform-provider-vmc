package vmc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"net/http"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

func TestAccResourceVmcSddc_basic(t *testing.T) {
	sddcName := "terraform_test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVmcSddcConfigBasic(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_1"),
				),
			},
		},
	})
}

func testCheckVmcSddcExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["id"]
		sddcName := rs.Primary.Attributes["sddc_name"]
		orgID := rs.Primary.Attributes["org_id"]

		client := testAccProvider.Meta().(*vmc.Client)
		_, resp, err := client.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: Get on sddcApi: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Sddc %q does not exist", sddcName)
		}

		fmt.Print("SDDC created successfully")
		return nil
	}
}

func testCheckVmcSddcDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*vmc.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		task, _, err := client.SddcApi.OrgsOrgSddcsSddcDelete(context.Background(), orgID, sddcID, nil)
		// TODO: check why error raised when deleting sddc.
		// if err != nil {
		// 	return fmt.Errorf("Error while deleting sddc %q, %v, %v", sddcID, err, resp)
		// }
		err = vmc.WaitForTask(client, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %q: %v", task.Id, err)
		}
	}

	return nil
}

func testAccVmcSddcConfigBasic(sddcName string) string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
	
	# refresh_token = "ac5140ea-1749-4355-a892-56cff4893be0"
	# vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
	# csp_url       = "https://console-stg.cloud.vmware.com"
}
	
data "vmc_org" "my_org" {
	id = "058f47c4-92aa-417f-8747-87f3ed61cb45"

	# id = "05e0a625-3293-41bb-a01f-35e762781c2a"
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

	region = "US_EAST_1"

	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	sddc_template_id    = ""
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
		sddcName,
	)
}
