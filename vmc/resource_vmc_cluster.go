/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	autoscalerapi "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api"
	autoscalercluster "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api/orgs/sddcs/clusters"
	autoscalermodel "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs/clusters/msft_licensing"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceClusterCreate,
		Delete: resourceClusterDelete,
		Update: resourceClusterUpdate,
		Read:   resourceClusterRead,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected id,sddc_id", d.Id())
				}
				if err := IsValidUUID(idParts[0]); err != nil {
					return nil, fmt.Errorf("invalid format for id : %v", err)
				}
				if err := IsValidUUID(idParts[1]); err != nil {
					return nil, fmt.Errorf("invalid format for sddc_id : %v", err)
				}

				d.SetId(idParts[0])
				d.Set("sddc_id", idParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
					[]string{HostInstancetypeI3, HostInstancetypeR5, HostInstancetypeI3EN}, false),
			},
			"storage_capacity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"15TB", "20TB", "25TB", "30TB", "35TB"}, false),
			},
			"edrs_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  StorageScaleUpPolicyType,
				ValidateFunc: validation.StringInSlice(
					[]string{StorageScaleUpPolicyType, CostPolicyType, PerformancePolicyType, RapidScaleUpPolicyType}, false),
				Description: "The EDRS policy type. This can either be 'cost', 'performance', 'storage-scaleup' or 'rapid-scaleup'. Default : storage-scaleup. ",
			},
			"enable_edrs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "True if EDRS is enabled",
			},
			"min_hosts": {
				Type:         schema.TypeInt,
				Default:      3,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 16),
				Description:  "The minimum number of hosts that the cluster can scale in to.",
			},
			"max_hosts": {
				Type:         schema.TypeInt,
				Default:      16,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 16),
				Description:  "The maximum number of hosts that the cluster can scale out to.",
			},
			"microsoft_licensing_config": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mssql_licensing": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The status of MSSQL licensing for this SDDC’s clusters. Possible values : ENABLED or DISABLED.",
							ValidateFunc: validation.StringInSlice([]string{
								LicenseConfigEnabled, LicenseConfigDisabled,
							}, false),
						},
						"windows_licensing": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The status of Windows licensing for this SDDC's clusters. Possible values : ENABLED, DISABLED or CUSTOMER'S ",
							ValidateFunc: validation.StringInSlice([]string{
								LicenseConfigEnabled, LicenseConfigDisabled,
							}, false),
						},
					},
				},
				Optional:    true,
				Description: "Indicates the desired licensing support, if any, of Microsoft software.",
			},
			"cluster_info": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
		CustomizeDiff: func(d *schema.ResourceDiff, meta interface{}) error {
			newInstanceType := d.Get("host_instance_type").(string)

			switch newInstanceType {

			case HostInstancetypeI3, HostInstancetypeI3EN:

				if d.Get("storage_capacity").(string) != "" {

					return fmt.Errorf("storage_capacity is not supported for host_instance_type %q", newInstanceType)

				}
			case HostInstancetypeR5:

				if d.Get("storage_capacity").(string) == "" {

					return fmt.Errorf("storage_capacity is required for host_instance_type %q", newInstanceType)

				}

			}
			return nil
		},
	}
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	var storageCapacityConverted int64
	sddcID := d.Get("sddc_id").(string)
	numHosts := int64(d.Get("num_hosts").(int))
	hostCPUCoresCount := int64(d.Get("host_cpu_cores_count").(int))

	connector := m.(*ConnectorWrapper)
	orgID := m.(*ConnectorWrapper).OrgID
	clusterClient := sddcs.NewDefaultClustersClient(connector)
	hostInstanceType := model.HostInstanceTypes(d.Get("host_instance_type").(string))
	storageCapacity := d.Get("storage_capacity").(string)
	if len(strings.TrimSpace(storageCapacity)) > 0 {
		storageCapacityConverted = ConvertStorageCapacitytoInt(storageCapacity)
	}
	msftLicensingConfig := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
	clusterConfig := &model.ClusterConfig{
		NumHosts:          numHosts,
		HostCpuCoresCount: &hostCPUCoresCount,
		HostInstanceType:  &hostInstanceType,
		StorageCapacity:   &storageCapacityConverted,
		MsftLicenseConfig: msftLicensingConfig,
	}

	task, err := clusterClient.Create(orgID, sddcID, *clusterConfig)
	if err != nil {
		return HandleCreateError("Cluster", err)
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
	if err != nil {
		return HandleReadError(d, "Cluster", clusterID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Unable to retrieve SDDC with ID %s", sddc.Id)
		d.SetId("")
		return nil
	}
	clusterExists := false

	for _, clusterConfig := range sddc.ResourceConfig.Clusters {
		if strings.Contains(clusterConfig.ClusterId, clusterID) {
			clusterExists = true
		}
	}
	if !clusterExists {
		log.Printf("Unable to retrieve cluster with ID %s", clusterID)
		d.SetId("")
		return nil
	}
	d.SetId(clusterID)
	cluster := map[string]string{}
	for _, clusterConfig := range sddc.ResourceConfig.Clusters {
		if strings.Contains(clusterConfig.ClusterId, clusterID) {
			cluster["cluster_name"] = *clusterConfig.ClusterName
			cluster["cluster_state"] = *clusterConfig.ClusterState
			cluster["host_instance_type"] = *clusterConfig.EsxHostInfo.InstanceType
			cluster["mssql_licensing"] = *clusterConfig.MsftLicenseConfig.MssqlLicensing
			cluster["windows_licensing"] = *clusterConfig.MsftLicenseConfig.WindowsLicensing
			d.Set("cluster_info", cluster)
			d.Set("num_hosts", len(clusterConfig.EsxHostList))
			break
		}
	}

	edrsPolicyClient := autoscalercluster.NewDefaultEdrsPolicyClient(connector)
	edrsPolicy, err := edrsPolicyClient.Get(orgID, sddcID, clusterID)
	if err != nil {
		return HandleReadError(d, "Cluster", clusterID, err)
	}
	d.Set("edrs_policy_type", *edrsPolicy.PolicyType)
	d.Set("enable_edrs", edrsPolicy.EnableEdrs)
	d.Set("max_hosts", *edrsPolicy.MaxHosts)
	d.Set("min_hosts", *edrsPolicy.MinHosts)
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
		return HandleDeleteError("Cluster", clusterID, err)
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
	connectorWrapper := m.(*ConnectorWrapper)
	esxsClient := sddcs.NewDefaultEsxsClient(connectorWrapper)
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
			return HandleUpdateError("Cluster", err)
		}
		tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while waiting for task %s: %v", task.Id, err))
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected hosts to be updated but were in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceClusterRead(d, m))
		})
		if err != nil {
			return err
		}
	}
	if d.HasChange("edrs_policy_type") || d.HasChange("enable_edrs") || d.HasChange("min_hosts") || d.HasChange("max_hosts") {
		edrsPolicyClient := autoscalercluster.NewDefaultEdrsPolicyClient(connectorWrapper)
		minHosts := int64(d.Get("min_hosts").(int))
		maxHosts := int64(d.Get("max_hosts").(int))
		policyType := d.Get("edrs_policy_type").(string)
		enableEDRS := d.Get("enable_edrs").(bool)
		edrsPolicy := &autoscalermodel.EdrsPolicy{
			EnableEdrs: enableEDRS,
			PolicyType: &policyType,
			MinHosts:   &minHosts,
			MaxHosts:   &maxHosts,
		}
		task, err := edrsPolicyClient.Post(orgID, sddcID, clusterID, *edrsPolicy)
		if err != nil {
			return HandleUpdateError("EDRS Policy", err)
		}
		return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			taskClient := autoscalerapi.NewDefaultAutoscalerClient(connectorWrapper)
			task, err := taskClient.Get(orgID, task.Id)
			if err != nil {
				if err.Error() == (errors.Unauthenticated{}.Error()) {
					log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
					err = connectorWrapper.authenticate()
					if err != nil {
						return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
					}
					return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
				}
				return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))

			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceClusterRead(d, m))
		})
	}
	// Update microsoft licensing config
	if d.HasChange("microsoft_licensing_config") {
		configChangeParam := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
		publishClient := msft_licensing.NewDefaultPublishClient(connectorWrapper)
		task, err := publishClient.Post(orgID, sddcID, clusterID, *configChangeParam)
		if err != nil {
			return HandleUpdateError("Microsoft Licensing Config", err)
		}
		return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				if err.Error() == (errors.Unauthenticated{}.Error()) {
					log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
					err = connectorWrapper.authenticate()
					if err != nil {
						return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
					}
					return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
				}
				return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))

			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceClusterRead(d, m))
		})

	}
	return nil
}
