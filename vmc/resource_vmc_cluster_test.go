/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
)

func TestAccResourceVmcCluster_basic(t *testing.T) {
	var sddcResource model.Sddc
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckVmcClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVmcClusterConfigBasic(),
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
			if strings.Contains(*currentResourceConfig.ClusterName, "Cluster-2") {
				clusterExists = true
				break
			}
		}

		if clusterExists != true {
			return fmt.Errorf("error retrieving Cluster : %v", err)
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
			if strings.Contains(*currentResourceConfig.ClusterName, "Cluster-2") {
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

func testAccVmcClusterConfigBasic() string {
	return fmt.Sprintf(`

resource "vmc_cluster" "cluster_1" {
	sddc_id = %q
	num_hosts      = 3
    }
`,
		os.Getenv(TestSDDCId),
	)
}
