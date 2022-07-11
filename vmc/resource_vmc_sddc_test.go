/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"os"
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

func testCheckVmcSddcExists(name string, sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["id"]
		sddcName := rs.Primary.Attributes["sddc_name"]

		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		orgID := connectorWrapper.OrgID
		connector := connectorWrapper.Connector

		sddcClient := orgs.NewSddcsClient(connector)
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

	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	sddcClient := orgs.NewSddcsClient(connector)

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
	size = "large"

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
		os.Getenv(AWSAccountNumber),
		sddcName,
	)
}

func TestBuildAwsSddcConfig(t *testing.T) {
	type test struct {
		input    map[string]interface{}
		expected model.AwsSddcConfig
	}

	tests := []test{
		{input: map[string]interface{}{
			"sddc_name":          "testName1",
			"region":             "us-east-1",
			"provider_type":      ZeroCloudProviderType,
			"num_host":           MinHosts,
			"host_instance_type": HostInstancetypeI3,
		},
			expected: model.AwsSddcConfig{
				Name:             "testName1",
				Region:           "us-east-1",
				Provider:         ZeroCloudProviderType,
				NumHosts:         MinHosts,
				HostInstanceType: String(model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL),
				Size:             String(MediumSDDCSize),
				DeploymentType:   String(SingleAvailabilityZone),
			}},
		{input: map[string]interface{}{
			"sddc_name":          "testName2",
			"region":             "us-east-2",
			"provider_type":      AWSProviderType,
			"num_host":           MaxHosts,
			"host_instance_type": HostInstancetypeI3EN,
		},
			expected: model.AwsSddcConfig{
				Name:             "testName2",
				Region:           "us-east-2",
				Provider:         AWSProviderType,
				NumHosts:         MaxHosts,
				HostInstanceType: String(model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL),
				Size:             String(MediumSDDCSize),
				DeploymentType:   String(SingleAvailabilityZone),
			}},
		{input: map[string]interface{}{
			"sddc_name":          "testName3",
			"region":             "us-west-1",
			"provider_type":      AWSProviderType,
			"num_host":           7,
			"host_instance_type": HostInstancetypeI4I,
		},
			expected: model.AwsSddcConfig{
				Name:             "testName3",
				Region:           "us-west-1",
				Provider:         AWSProviderType,
				NumHosts:         7,
				HostInstanceType: String(model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL),
				Size:             String(MediumSDDCSize),
				DeploymentType:   String(SingleAvailabilityZone),
			}},
		{input: map[string]interface{}{
			"sddc_name":          "testName4",
			"region":             "us-west-1",
			"provider_type":      AWSProviderType,
			"num_host":           MaxHosts,
			"host_instance_type": HostInstancetypeR5,
		},
			expected: model.AwsSddcConfig{
				Name:             "testName4",
				Region:           "us-west-1",
				Provider:         AWSProviderType,
				NumHosts:         MaxHosts,
				HostInstanceType: String(model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL),
				Size:             String(MediumSDDCSize),
				DeploymentType:   String(SingleAvailabilityZone),
			}},
	}

	for _, testCase := range tests {
		var mockResourceSchema = schema.TestResourceDataRaw(t, sddcSchema(), testCase.input)
		got, _ := buildAwsSddcConfig(mockResourceSchema)
		assert.Equal(t, got.NumHosts, testCase.expected.NumHosts)
		assert.Equal(t, got.SddcType, testCase.expected.SddcType)
		assert.Equal(t, got.SddcId, testCase.expected.SddcId)
		assert.Equal(t, got.Region, testCase.expected.Region)
		assert.Equal(t, got.Provider, testCase.expected.Provider)
		assert.Equal(t, *got.HostInstanceType, *testCase.expected.HostInstanceType)
		assert.Equal(t, got.Size, testCase.expected.Size)
		assert.Equal(t, got.DeploymentType, testCase.expected.DeploymentType)
	}
}

// String returns a pointer to the string value passed in.
func String(v string) *string {
	return &v
}
