/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"log"
	"strings"
	"time"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Delete: resourceClusterDelete,
		Update: resourceClusterUpdate,
		Read:   resourceClusterRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(35 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(40 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "SDDC identifier",
			},
			"num_hosts": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(3, 16),
				Description:  "The number of hosts.",
			},
			"host_cpu_cores_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Customize CPU cores on hosts in a cluster. Specify number of cores to be enabled on hosts in a cluster.",
			},
			"host_instance_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The instance type for the esx hosts added to this cluster.",
				ValidateFunc: validation.StringInSlice(
					[]string{HostInstancetypeI3, HostInstancetypeR5}, false),
			},
			"storage_capacity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "For EBS-backed instances only, the requested storage capacity in GiB.",
			},
			"cluster_info": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	sddcID := d.Get("sddc_id").(string)
	numHosts := int64(d.Get("num_hosts").(int))
	hostCPUCoresCount := int64(d.Get("host_cpu_cores_count").(int))

	connector := m.(*ConnectorWrapper)
	orgID := m.(*ConnectorWrapper).OrgID
	clusterClient := sddcs.NewDefaultClustersClient(connector)

	clusterConfig := &model.ClusterConfig{
		NumHosts:          numHosts,
		HostCpuCoresCount: &hostCPUCoresCount,
		// To be added : support for other host instance types
		//HostInstanceType:  &hostInstanceType,
		//StorageCapacity:   &storageCapacityConverted,
	}

	task, err := clusterClient.Create(orgID, sddcID, *clusterConfig)
	if err != nil {
		return fmt.Errorf("error while creating cluster for SDDC with id %s: %v", sddcID, err)
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := orgs.NewDefaultTasksClient(connector)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
				err = connector.authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
				}
				return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))

		}
		if task.Params.HasField(ClusterIdFieldName) {
			clusterID, err := task.Params.String(ClusterIdFieldName)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error getting clusterId : %s", err))

			}
			d.SetId(clusterID)
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceClusterRead(d, m))
	})
}

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	clusterID := d.Id()
	sddcID := d.Get("sddc_id").(string)
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddc, err := GetSDDC(connector, orgID, sddcID)
	log.Printf("SDDC ID : %s", sddcID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("SDDC with ID %s not found", sddcID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("error while getting the SDDC with ID %s,%v", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Unable to retrieve SDDC with ID %s", sddc.Id)
		d.SetId("")
		return nil
	}

	d.SetId(clusterID)
	cluster := map[string]string{}
	for i := 0; i < len(sddc.ResourceConfig.Clusters); i++ {
		currentResourceConfig := sddc.ResourceConfig.Clusters[i]
		if strings.Contains(currentResourceConfig.ClusterId, clusterID) {
			cluster["cluster_name"] = *currentResourceConfig.ClusterName
			cluster["cluster_state"] = *currentResourceConfig.ClusterState
			cluster["host_instance_type"] = *currentResourceConfig.EsxHostInfo.InstanceType
			d.Set("cluster_info", cluster)
			break
		}
	}
	return nil
}

func resourceClusterDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	clusterID := d.Id()

	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	clusterClient := sddcs.NewDefaultClustersClient(connector)
	task, err := clusterClient.Delete(orgID, sddcID, clusterID)
	if err != nil {
		fmt.Printf("error : %v", err)
	}
	tasksClient := orgs.NewDefaultTasksClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error while deleting SDDC %s: %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	esxsClient := sddcs.NewDefaultEsxsClient(connector)
	sddcID := d.Get("sddc_id").(string)
	orgID := (m.(*ConnectorWrapper)).OrgID
	clusterID := d.Id()

	// Add or remove hosts from a cluster
	if d.HasChange("num_hosts") {
		oldTmp, newTmp := d.GetChange("num_hosts")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}

		esxConfig := model.EsxConfig{
			NumHosts:  int64(diffNum),
			ClusterId: &clusterID,
		}

		task, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)

		if err != nil {
			return fmt.Errorf("error while updating number of host for cluster %s: %v", clusterID, err)
		}
		tasksClient := orgs.NewDefaultTasksClient(connector)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while waiting for task %s: %v", task.Id, err))
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected host to be updated but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceClusterRead(d, m))
		})
		if err != nil {
			return err
		}
	}
	return nil
}
