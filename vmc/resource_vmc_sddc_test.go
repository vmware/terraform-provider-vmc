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
	"gitlab.eng.vmware.com/het/vmc-go-sdk/vmc"
)

func TestAccResourceVmcSddc_basic(t *testing.T) {
	sddcName := "test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccVmcSddcConfigBasic(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists(sddcName),
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

		client := testAccProvider.Meta().(*vmc.APIClient)
		_, resp, err := client.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
		if err != nil {
			return fmt.Errorf("Bad: Get on sddcApi: %s", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("Bad: Sddc %q does not exist", sddcName)
		}
		return nil
	}
}

func testCheckVmcSddcDestroy(s *terraform.State) error {

	client := testAccProvider.Meta().(*vmc.APIClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		task, _, err := client.SddcApi.OrgsOrgSddcsSddcDelete(context.Background(), orgID, sddcID, nil)
		if err != nil {
			return fmt.Errorf("Error while deleting sddc %q: %v", sddcID, err)
		}
		err = waitForTask(client, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %q: %v", task.Id, err)
		}
	}

	return nil
}

func testAccVmcSddcConfigBasic(sddcName string) string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = "340354ab-92ca-4477-acfb-62aeaf6bd74f"
}


data "vmc_org" "test_org" {
	id = %q
}

resource "vmc_sddc" "test_sddc" {
	org_id        = "${data.vmc_org.test_org.id}"
	sddc_name     = %q
	num_host      = 1
	provider_type = "ZEROCLOUD"
	region        = "US_WEST_2"
}	
`,
		os.Getenv("VMC_ORG_ID"),
		sddcName,
	)
}
