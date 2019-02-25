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
					testCheckVmcSddcExists("vmc_sddc.test_sddc"),
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

		fmt.Print("SDDC created successfully")
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
	csp_url       = "https://console-stg.cloud.vmware.com"
  }

data "vmc_org" "test_org" {
	id = "05e0a625-3293-41bb-a01f-35e762781c2a"
}

resource "vmc_sddc" "test_sddc" {
	org_id        = "${data.vmc_org.test_org.id}"
	sddc_name     = %q
	num_host      = 4
	provider_type = "ZEROCLOUD"
	region        = "US_WEST_1"
}	
`,
		os.Getenv("REFRESH_TOKEN"),
		sddcName,
	)
}
