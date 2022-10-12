/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
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

func TestAccResourceVmcClusterBasic(t *testing.T) {
	var sddcResource model.Sddc
	resourceName := "vmc_cluster.cluster_1"
	sddcName := "terraform_test_sddc_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcClusterConfigBasic(sddcName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVmcClusterExists("vmc_cluster.cluster_1", &sddcResource),
					resource.TestCheckResourceAttrSet("vmc_cluster.cluster_1", "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccVmcClusterResourceImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccResourceVmcClusterZerocloud(t *testing.T) {
	var sddcResource model.Sddc
	clusterRef := "cluster_zerocloud"
	resourceName := "vmc_cluster." + clusterRef
	sddcName := "terraform_cluster_test_" + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckZerocloud(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcClusterConfigBasicZerocloud(sddcName, clusterRef),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVmcClusterExists(resourceName, &sddcResource),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportStateIdFunc: testAccVmcClusterResourceImportStateIdFunc(resourceName),
				ImportState:       true,
				ImportStateVerify: true,
				// "microsoft_licensing_config" and "host_instance_type" are set in the
				// cluster_info map, not on the cluster resource itself.
				ImportStateVerifyIgnore: []string{"microsoft_licensing_config", "host_instance_type"},
			},
		},
	})
}

func testAccCheckVmcClusterExists(clusterRef string, sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[clusterRef]
		if !ok {
			return fmt.Errorf("not found: %s", clusterRef)
		}
		sddcID := rs.Primary.Attributes["sddc_id"]

		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		orgID := connectorWrapper.OrgID
		connector := connectorWrapper.Connector

		sddcClient := orgs.NewSddcsClient(connector)
		var err error

		*sddcResource, err = sddcClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("error retrieving SDDC : %s", err)
		}
		clusterExists := false
		for i := 0; i < len(sddcResource.ResourceConfig.Clusters); i++ {
			currentResourceConfig := sddcResource.ResourceConfig.Clusters[i]
			if strings.HasSuffix(*currentResourceConfig.ClusterName, "-2") {
				clusterExists = true
				break
			}
		}

		if clusterExists != true {
			return fmt.Errorf("error retrieving cluster : %v", err)
		}

		fmt.Printf("Cluster for SDDC %s created successfully \n", sddcID)

		return nil
	}
}

func testCheckVmcClusterDestroy(s *terraform.State) error {
	connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
	connector := connectorWrapper.Connector
	sddcClient := orgs.NewSddcsClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_cluster" {
			continue
		}

		sddcID := rs.Primary.Attributes["sddc_id"]
		clusterId := rs.Primary.Attributes["id"]
		orgID := connectorWrapper.OrgID
		sddcResource, err := sddcClient.Get(orgID, sddcID)

		for i := 0; i < len(sddcResource.ResourceConfig.Clusters); i++ {
			currentResourceConfig := sddcResource.ResourceConfig.Clusters[i]
			if currentResourceConfig.ClusterId == clusterId {
				return fmt.Errorf("cluster still exists : %v", err)
			}
		}

		// check if error type is not_found
		if err != nil {
			if err.Error() != (errors.NotFound{}.Error()) {
				return err
			}

		}
	}

	return nil
}

func testAccVmcClusterConfigBasic(sddcName string) string {
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
	host_instance_type = "I3_METAL"
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
}

resource "vmc_cluster" "cluster_1" {
	sddc_id = vmc_sddc.sddc_1.id
	host_instance_type = "I3_METAL"
	num_hosts = 3
	microsoft_licensing_config {
        mssql_licensing = "DISABLED"
        windows_licensing = "ENABLED"
    }
}
`,
		os.Getenv(AWSAccountNumber),
		sddcName,
	)
}

func testAccVmcClusterConfigBasicZerocloud(sddcName string, clusterRef string) string {
	return fmt.Sprintf(`

data "vmc_connected_accounts" "my_accounts" {
      account_number = %q
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = "US_WEST_2"
}

resource "vmc_sddc" "sddc_zerocloud_cluster" {
	sddc_name = %q
	vpc_cidr      = "10.2.0.0/16"
	num_host      = 2
	provider_type = "ZEROCLOUD"
	host_instance_type = "I3_METAL"
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
}

resource "vmc_cluster" %q {
	sddc_id = vmc_sddc.sddc_zerocloud_cluster.id
	host_instance_type = "I3_METAL"
	num_hosts = 2
	microsoft_licensing_config {
        mssql_licensing = "DISABLED"
        windows_licensing = "ENABLED"
    }
}

resource "vmc_cluster" "cluster_2" {
	sddc_id = vmc_sddc.sddc_zerocloud_cluster.id
	host_instance_type = "I3_METAL"
	num_hosts = 2
	microsoft_licensing_config {
        mssql_licensing = "DISABLED"
        windows_licensing = "ENABLED"
    }
}
`,
		os.Getenv(AWSAccountNumber),
		sddcName,
		clusterRef,
	)
}

func testAccVmcClusterResourceImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}
		return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["sddc_id"]), nil
	}
}

func TestBuildClusterConfig(t *testing.T) {
	type test struct {
		instanceType             string
		expectedHostInstanceType string
		expectedErr              error
	}

	tests := []test{
		{instanceType: HostInstancetypeI3,
			expectedHostInstanceType: model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL,
			expectedErr:              nil},
		{instanceType: HostInstancetypeI3EN,
			expectedHostInstanceType: model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL,
			expectedErr:              nil},
		{instanceType: HostInstancetypeI4I,
			expectedHostInstanceType: model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL,
			expectedErr:              nil},
		{instanceType: HostInstancetypeR5,
			expectedHostInstanceType: model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL,
			expectedErr:              nil},
		{instanceType: "RandomString",
			expectedHostInstanceType: "",
			expectedErr:              fmt.Errorf("unknown host instance type: RandomString"),
		},
	}

	for _, testCase := range tests {
		config := map[string]interface{}{
			"num_hosts":          MinHosts,
			"host_instance_type": testCase.instanceType,
		}
		var testResourceSchema = schema.TestResourceDataRaw(t, clusterSchema(), config)
		got, err := buildClusterConfig(testResourceSchema)
		assert.Equal(t, err, testCase.expectedErr)
		if err == nil {
			assert.Equal(t, *got.HostInstanceType, testCase.expectedHostInstanceType)
		}
	}
}
