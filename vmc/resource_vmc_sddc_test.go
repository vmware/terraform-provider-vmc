/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
)

func TestAccResourceVmcSddc_basic(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigBasic(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_1", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_1", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "cloud_password"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_1", "nsxt_reverse_proxy_url"),
				),
			},
		},
	})
}

func TestAccResourceVmcSddcZerocloud(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_sddc_test_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigZerocloud(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_zerocloud", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_zerocloud", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "cloud_password"),
				),
			},
		},
	})
}

func TestAccResourceVmcSddcRequiredFieldsOnlyZerocloud(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_sddc_test_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigRequiredFieldsZerocloud(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_zerocloud", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_zerocloud", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_zerocloud", "cloud_password"),
				),
			},
		},
	})
}

func TestAccResourceVmcSddcC6iMetal(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_sddc_c6i_test_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigDiskless(sddcName, constants.HostInstancetypeC6I),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_c6i", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_c6i", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_c6i", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_c6i", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_c6i", "cloud_password"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_c6i", "nsxt_reverse_proxy_url"),
				),
			},
		},
	})
}

func TestAccResourceVmcSddcM7i24xlMetal(t *testing.T) {
	var sddcResource model.Sddc
	sddcName := "terraform_sddc_m7i_24_xl_test_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcSddcDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcSddcConfigDiskless(sddcName, constants.HostInstancetypeM7i24xl),
				Check: resource.ComposeTestCheckFunc(
					testCheckVmcSddcExists("vmc_sddc.sddc_c6i", &sddcResource),
					testCheckSddcAttributes(&sddcResource),
					resource.TestCheckResourceAttr("vmc_sddc.sddc_sddc_m7i_24xl", "sddc_state", "READY"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_sddc_m7i_24xl", "vc_url"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_sddc_m7i_24xl", "cloud_username"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_sddc_m7i_24xl", "cloud_password"),
					resource.TestCheckResourceAttrSet("vmc_sddc.sddc_sddc_m7i_24xl", "nsxt_reverse_proxy_url"),
				),
			},
		},
	})
}

func testCheckVmcSddcExists(name string, sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["id"]
		sddcName := rs.Primary.Attributes["sddc_name"]

		connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
		orgID := connectorWrapper.OrgID

		sddcClient := orgs.NewSddcsClient(connectorWrapper)
		var err error
		*sddcResource, err = sddcClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("error retrieving SDDC : %s", err)
		}

		if sddcResource.Id != sddcID {
			return fmt.Errorf("error : SDDC %q does not exist", sddcName)
		}

		fmt.Printf("SDDC %s created successfully with id %s \n", sddcName, sddcID)
		return nil
	}
}

func testCheckSddcAttributes(sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		sddcState := sddcResource.SddcState
		if *sddcState != "READY" {
			return fmt.Errorf(" SDDC %s with ID %s is not in ready state", *sddcResource.Name, sddcResource.Id)
		}
		return nil
	}
}

func testCheckVmcSddcDestroy(s *terraform.State) error {

	connectorWrapper := testAccProvider.Meta().(*connector.Wrapper)
	sddcClient := orgs.NewSddcsClient(connectorWrapper)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_sddc" {
			continue
		}

		sddcID := rs.Primary.Attributes["id"]
		orgID := rs.Primary.Attributes["org_id"]
		sddc, err := sddcClient.Get(orgID, sddcID)
		if err == nil {
			if *sddc.SddcState != "DELETED" {
				return fmt.Errorf("SDDC %s with ID %s still exits", *sddc.Name, sddc.Id)
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

func testAccVmcSddcConfigBasic(sddcName string) string {
	return fmt.Sprintf(`

data "vmc_connected_accounts" "my_accounts" {
      account_number = %q
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = "US_WEST_2"
  sddc_type = "SingleAZ"
  instance_type = "i3.metal"
}

resource "vmc_sddc" "sddc_1" {
	sddc_name = %q
	vpc_cidr      = "10.2.0.0/16"
	num_host      = 3
	provider_type = "AWS"

	region = "US_WEST_2"
	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"
	
	account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
	}
	microsoft_licensing_config {
        mssql_licensing = "DISABLED"
        windows_licensing = "ENABLED"
    }
    timeouts {
      create = "300m"
      update = "300m"
      delete = "180m"
  }
}
`,
		os.Getenv(constants.AwsAccountNumber),
		sddcName,
	)
}

func testAccVmcSddcConfigZerocloud(sddcName string) string {
	return fmt.Sprintf(`

resource "vmc_sddc" "sddc_zerocloud" {
	sddc_name = %q
	vpc_cidr      = "10.40.0.0/16"
	num_host      = 2
	provider_type = "ZEROCLOUD"
	host_instance_type = "I3_METAL"
	region = "US_WEST_2"
	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"

    timeouts {
      create = "300m"
      update = "300m"
      delete = "180m"
  	}

	microsoft_licensing_config {
		mssql_licensing = "ENABLED"
		windows_licensing = "DISABLED"
	}
}
`,
		sddcName,
	)
}
func testAccVmcSddcConfigDiskless(sddcName string, hostInstanceType string) string {

	sddcResourceName := "sddc_" + strings.Replace(strings.ToLower(hostInstanceType), "_metal", "", 1)
	return fmt.Sprintf(`
data "vmc_connected_accounts" "my_accounts" {
  account_number = %q
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = "US_WEST_2"
}

resource "vmc_sddc" %q {
	sddc_name = %q
	vpc_cidr      = "10.40.0.0/16"
	num_host      = 3
	provider_type = "ZEROCLOUD"
	host_instance_type = %q
	region = "US_WEST_2"
	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"
	account_link_sddc_config {
		customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
		connected_account_id = data.vmc_connected_accounts.my_accounts.id
	  }

    timeouts {
      create = "300m"
      update = "300m"
      delete = "180m"
  	}

	microsoft_licensing_config {
		mssql_licensing = "ENABLED"
		windows_licensing = "DISABLED"
	}
}
`,
		os.Getenv(constants.AwsAccountNumber),
		sddcResourceName,
		sddcName,
		hostInstanceType,
	)
}

func testAccVmcSddcConfigRequiredFieldsZerocloud(sddcName string) string {
	return fmt.Sprintf(`

resource "vmc_sddc" "sddc_zerocloud" {
	sddc_name = %q
	num_host  = 2
	provider_type = "ZEROCLOUD"
	region = "US_WEST_2"
}
`,
		sddcName,
	)
}

func TestBuildAwsSddcConfigHostInstanceType(t *testing.T) {
	type test struct {
		input    map[string]interface{}
		expected string
		err      error
	}

	tests := []test{
		{input: map[string]interface{}{
			"host_instance_type": constants.HostInstancetypeI3,
		},
			expected: model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL,
			err:      nil,
		},
		{input: map[string]interface{}{
			"host_instance_type": constants.HostInstancetypeI3EN,
		},
			expected: model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL,
			err:      nil,
		},
		{input: map[string]interface{}{
			"host_instance_type": constants.HostInstancetypeI4I,
		},
			expected: model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL,
			err:      nil,
		},
		{input: map[string]interface{}{
			"host_instance_type": "RandomString",
		},
			expected: "",
			err:      fmt.Errorf("unknown host instance type: RandomString")},
	}

	for _, testCase := range tests {
		var testResourceSchema = schema.TestResourceDataRaw(t, sddcSchema(), testCase.input)
		got, err := buildAwsSddcConfig(testResourceSchema)
		assert.Equal(t, testCase.err, err)
		if err == nil {
			assert.Equal(t, testCase.expected, *got.HostInstanceType)
		}
	}
}
