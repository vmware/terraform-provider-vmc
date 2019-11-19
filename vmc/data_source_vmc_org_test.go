package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcOrg_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcOrgConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_org.my_org", "display_name", "VMC Org"),
					resource.TestCheckResourceAttr("data.vmc_org.my_org", "name", "j4acl4e3"),
				),
			},
		},
	})
}

func testAccDataSourceVmcOrgConfig() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
    csp_url       = "https://console-stg.cloud.vmware.com"
    vmc_url = "https://stg.skyscraper.vmware.com"
}
	
data "vmc_org" "my_org" {
	id = %q

}
`,
		os.Getenv("REFRESH_TOKEN"),
		os.Getenv("ORG_ID"),
	)
}
