package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"os"
	"testing"
)

func TestAccDataSourceVmcCustomerSubnets_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceVmcCustomerSubnetsConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.#", "3"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.0", "subnet-13a0c249"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.1", "subnet-14a42d6d"),
					resource.TestCheckResourceAttr("data.vmc_customer_subnets.my_subnets", "ids.2", "subnet-2170db6a"),
				),
			},
		},
	})
}

func testAccDataSourceVmcCustomerSubnetsConfig() string {
	return fmt.Sprintf(`
provider "vmc" {
	refresh_token = %q
}
	
data "vmc_org" "my_org" {
	id = "058f47c4-92aa-417f-8747-87f3ed61cb45"
}
	
data "vmc_connected_accounts" "my_accounts" {
	org_id = "${data.vmc_org.my_org.id}"
}

data "vmc_customer_subnets" "my_subnets" {
	org_id = "${data.vmc_org.my_org.id}"
	connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
	region = "us-west-2"
}
`,
		os.Getenv("REFRESH_TOKEN"),
	)
}
