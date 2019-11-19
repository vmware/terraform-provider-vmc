package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/sddcs"
	"gitlab.eng.vmware.com/het/vmware-vmc-sdk/vapi/bindings/vmc/orgs/tasks"
	"os"
	"strings"
	"testing"
	"time"
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
		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		connector := connectorWrapper.Connector
		sddcClient := sddcs.NewSddcsClientImpl(connector)

		sddc, err := sddcClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: Get on sddcApi: %s", err)
		}

		if sddc.Id != sddcID {
			return fmt.Errorf("Bad: Sddc %q does not exist", sddcName)
		}

		fmt.Printf("SDDC %s created successfully with id %s ", sddcName, sddcID)
		return nil
	}
}

func testCheckVmcSddcDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	sddcClient := sddcs.NewSddcsClientImpl(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("Error while deleting sddc %s, %s", sddcID, err)
		}
		tasksClient := tasks.NewTasksClientImpl(connector)
		err = resource.Retry(180*time.Minute, func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err))
			}
			if task.ErrorMessage != nil && strings.Contains(*task.ErrorMessage, "Entity is not found") {
				fmt.Print("Resource already deleted")
				return resource.NonRetryableError(nil)
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(nil)
		})
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
	vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
  	csp_url       = "https://console-stg.cloud.vmware.com"
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
	num_host      = 3
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
