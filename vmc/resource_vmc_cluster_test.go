/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
)

func TestAccResourceVmcCluster_basic(t *testing.T) {
	var sddcResource model.Sddc
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
		},
	})
}

func testAccCheckVmcClusterExists(name string, sddcResource *model.Sddc) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		sddcID := rs.Primary.Attributes["sddc_id"]

		connectorWrapper := testAccProvider.Meta().(*ConnectorWrapper)
		orgID := connectorWrapper.OrgID
		connector := connectorWrapper.Connector

		sddcClient := orgs.NewDefaultSddcsClient(connector)
		var err error
		fmt.Printf("SDDC ID : %s", sddcID)
		fmt.Printf("Org ID : %s", orgID)
		*sddcResource, err = sddcClient.Get(orgID, sddcID)
		if err != nil {
			return fmt.Errorf("error retrieving SDDC : %s", err)
		}
		clusterExists := false
		for i := 0; i < len(sddcResource.ResourceConfig.Clusters); i++ {
			currentResourceConfig := sddcResource.ResourceConfig.Clusters[i]
			if strings.Contains(*currentResourceConfig.ClusterName, "Cluster-1") {
				fmt.Printf("Cluster Name : %s", *currentResourceConfig.ClusterName)
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
	sddcClient := orgs.NewDefaultSddcsClient(connector)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "vmc_cluster" {
			continue
		}

		sddcID := rs.Primary.Attributes["sddc_id"]
		orgID := connectorWrapper.OrgID
		sddcResource, err := sddcClient.Get(orgID, sddcID)

		for i := 0; i < len(sddcResource.ResourceConfig.Clusters); i++ {
			currentResourceConfig := sddcResource.ResourceConfig.Clusters[i]
			if !strings.Contains(*currentResourceConfig.ClusterName, "Cluster-1") {
				fmt.Printf("Cluster Name : %s", *currentResourceConfig.ClusterName)
				return fmt.Errorf("cluster still exists : %v", err)

			}
		}

		//check if error type if not_found
		if err.Error() != (errors.NotFound{}.Error()) {
			return err
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
	num_hosts      = 3
    }

`,
		os.Getenv(AWSAccountNumber),
		sddcName,
	)
}
