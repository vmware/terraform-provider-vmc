/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	autoscalercluster "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api/orgs/sddcs/clusters"
	autoscalermodel "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs/clusters/msft_licensing"
)

// clusterMutationKeyedMutex a mutex that allows only a single cluster per sddc to be mutated.
var clusterMutationKeyedMutex = KeyedMutex{}

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
		Schema: clusterSchema(),
		CustomizeDiff: func(c context.Context, d *schema.ResourceDiff, meta interface{}) error {
			newInstanceType := d.Get("host_instance_type").(string)

			switch newInstanceType {
			case HostInstancetypeI3, HostInstancetypeI3EN, HostInstancetypeI4I:
				if d.Get("storage_capacity").(string) != "" {
					return fmt.Errorf("storage_capacity is not supported for host_instance_type %q", newInstanceType)
				}
			case HostInstancetypeR5:
				if d.Get("storage_capacity").(string) == "" {
					return fmt.Errorf("storage_capacity is required for host_instance_type %q "+
						"Possible values are 15TB, 20TB, 25TB, 30TB, 35TB per host", newInstanceType)
				}
			}
			return nil
		},
	}
}

// clusterSchema this helper function extracts the creation of the Cluster schema, so that
// it's made available for mocking in tests.
func clusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sddc_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "SDDC identifier",
		},
		"num_hosts": {
			Type:         schema.TypeInt,
			Required:     true,
			ValidateFunc: validation.IntBetween(MinHosts, MaxHosts),
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
				[]string{HostInstancetypeI3, HostInstancetypeR5, HostInstancetypeI3EN, HostInstancetypeI4I}, false),
		},
		"storage_capacity": {
			Type:     schema.TypeString,
			Optional: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"15TB", "20TB", "25TB", "30TB", "35TB"}, false),
		},
		"edrs_policy_type": {
			Type: schema.TypeString,
			// Exact value known after create
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				[]string{StorageScaleUpPolicyType, CostPolicyType, PerformancePolicyType, RapidScaleUpPolicyType}, false),
			Description: "The EDRS policy type. This can either be 'cost', 'performance', 'storage-scaleup' or 'rapid-scaleup'. Default : storage-scaleup. ",
		},
		"enable_edrs": {
			Type: schema.TypeBool,
			// Value can be changed after create
			Optional:    true,
			Computed:    true,
			Description: "True if EDRS is enabled",
		},
		"min_hosts": {
			Type: schema.TypeInt,
			// Exact value known after create
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(MinHosts, MaxHosts),
			Description:  "The minimum number of hosts that the cluster can scale in to.",
		},
		"max_hosts": {
			Type: schema.TypeInt,
			// Exact value known after create
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(MinHosts, MaxHosts),
			Description:  "The maximum number of hosts that the cluster can scale out to.",
		},
		"microsoft_licensing_config": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"mssql_licensing": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The status of MSSQL licensing for this SDDC’s clusters. Possible values : enabled, ENABLED, disabled, DISABLED.",
						ValidateFunc: validation.StringInSlice([]string{
							LicenseConfigEnabled, LicenseConfigDisabled, CapitalLicenseConfigEnabled, CapitalLicenseConfigDisabled}, false),
					},
					"windows_licensing": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The status of Windows licensing for this SDDC's clusters. Possible values : enabled, ENABLED, disabled, DISABLED.",
						ValidateFunc: validation.StringInSlice([]string{
							LicenseConfigEnabled, LicenseConfigDisabled, CapitalLicenseConfigEnabled, CapitalLicenseConfigDisabled}, false),
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
	}
}

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	sddcID := d.Get("sddc_id").(string)
	clusterConfig, err := buildClusterConfig(d)
	if err != nil {
		return HandleCreateError("Cluster", err)
	}
	// Obtain a lock to allow only a single cluster creation at a time for a specific SDDC.
	var unlockFunction = clusterMutationKeyedMutex.lock(sddcID)
	connectorWrapper := m.(*ConnectorWrapper)
	orgID := m.(*ConnectorWrapper).OrgID
	clusterClient := sddcs.NewClustersClient(connectorWrapper)
	task, err := clusterClient.Create(orgID, sddcID, *clusterConfig)
	if err != nil {
		return HandleCreateError("Cluster", err)
	}
	var clusterID = ""
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := retryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return getTask(connectorWrapper, task.Id)
			},
			"error creating cluster ",
			func(task model.Task) {
				unlockFunction()
				// Obtain the ID of the newly created cluster
				if task.Params.HasField(ClusterIdFieldName) {
					clusterID, err = task.Params.String(ClusterIdFieldName)
					d.SetId(clusterID)
				}
			})
		if taskErr != nil {
			return taskErr
		}
		if clusterID == "" {
			return resource.NonRetryableError(fmt.Errorf("error getting clusterID"))
		}
		err = resourceClusterRead(d, m)
		if err == nil {
			return nil
		}
		return resource.NonRetryableError(err)
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
		if clusterConfig.ClusterId == clusterID {
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
		if clusterConfig.ClusterId == clusterID {
			cluster["cluster_name"] = *clusterConfig.ClusterName
			cluster["cluster_state"] = *clusterConfig.ClusterState
			cluster["host_instance_type"] = *clusterConfig.EsxHostInfo.InstanceType
			if clusterConfig.MsftLicenseConfig != nil {
				if clusterConfig.MsftLicenseConfig.MssqlLicensing != nil {
					cluster["mssql_licensing"] = *clusterConfig.MsftLicenseConfig.MssqlLicensing
				}
				if clusterConfig.MsftLicenseConfig.WindowsLicensing != nil {
					cluster["windows_licensing"] = *clusterConfig.MsftLicenseConfig.WindowsLicensing
				}
			}
			d.Set("cluster_info", cluster)
			d.Set("num_hosts", len(clusterConfig.EsxHostList))
			break
		}
	}

	edrsPolicyClient := autoscalercluster.NewEdrsPolicyClient(connector)
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
	connectorWrapper := m.(*ConnectorWrapper)
	clusterID := d.Id()

	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	var unlockFunction = clusterMutationKeyedMutex.lock(sddcID)
	clusterClient := sddcs.NewClustersClient(connectorWrapper)
	task, err := clusterClient.Delete(orgID, sddcID, clusterID)
	if err != nil {
		return HandleDeleteError("Cluster", clusterID, err)
	}
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		taskErr := retryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return getTask(connectorWrapper, task.Id)
			},
			"error deleting cluster "+clusterID,
			func(task model.Task) {
				unlockFunction()
			})
		if taskErr != nil {
			return taskErr
		}
		d.SetId("")
		return nil
	})
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*ConnectorWrapper)
	esxsClient := sddcs.NewEsxsClient(connectorWrapper)
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

		var unlockFunction = clusterMutationKeyedMutex.lock(sddcID)
		task, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)
		if err != nil {
			return HandleUpdateError("Cluster", err)
		}
		err = resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			taskErr := retryTaskUntilFinished(connectorWrapper,
				func() (model.Task, error) {
					return getTask(connectorWrapper, task.Id)
				},
				"error updating hosts for cluster "+clusterID,
				func(task model.Task) {
					unlockFunction()
				})
			if taskErr != nil {
				return taskErr
			}
			err = resourceClusterRead(d, m)
			if err == nil {
				return nil
			}
			return resource.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
	}
	if d.HasChange("edrs_policy_type") || d.HasChange("enable_edrs") || d.HasChange("min_hosts") || d.HasChange("max_hosts") {
		edrsPolicyClient := autoscalercluster.NewEdrsPolicyClient(connectorWrapper)
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
		if policyType == StorageScaleUpPolicyType && !enableEDRS {
			return fmt.Errorf("EDRS policy %s is the default and cannot be disabled", StorageScaleUpPolicyType)
		}
		var unlockFunction = clusterMutationKeyedMutex.lock(sddcID)
		task, err := edrsPolicyClient.Post(orgID, sddcID, clusterID, *edrsPolicy)
		if err != nil {
			return HandleUpdateError("EDRS Policy", err)
		}
		// The taskClient here has the same methods, but they return the autoscalerapi.Task instead of model.Task,
		// so the retryTaskUntilFinished method cannot be used
		return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			taskErr := retryTaskUntilFinished(connectorWrapper,
				func() (model.Task, error) {
					return getAutoscalerTask(connectorWrapper, task.Id)
				},
				"error updating EDRS policy configuration "+clusterID,
				func(task model.Task) {
					unlockFunction()
				})
			if taskErr != nil {
				return taskErr
			}
			err = resourceClusterRead(d, m)
			if err == nil {
				return nil
			}
			return resource.NonRetryableError(err)
		})
	}
	// Update microsoft licensing config
	if d.HasChange("microsoft_licensing_config") {
		configChangeParam := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
		publishClient := msft_licensing.NewPublishClient(connectorWrapper)
		var unlockFunction = clusterMutationKeyedMutex.lock(sddcID)
		task, err := publishClient.Post(orgID, sddcID, clusterID, *configChangeParam)
		if err != nil {
			return HandleUpdateError("Microsoft Licensing Config", err)
		}
		return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			taskErr := retryTaskUntilFinished(connectorWrapper,
				func() (model.Task, error) {
					return getTask(connectorWrapper, task.Id)
				},
				"error updating microsoft licensing configuration "+clusterID,
				func(task model.Task) {
					unlockFunction()
				})
			if taskErr != nil {
				return taskErr
			}
			err = resourceClusterRead(d, m)
			if err == nil {
				return nil
			}
			return resource.NonRetryableError(err)
		})

	}
	return nil
}

// buildClusterConfig extracts the creation of the model.ClusterConfig, so that it's
// available for testing
func buildClusterConfig(d *schema.ResourceData) (*model.ClusterConfig, error) {
	numHosts := int64(d.Get("num_hosts").(int))
	hostCPUCoresCount := int64(d.Get("host_cpu_cores_count").(int))

	hostInstanceType, err := toHostInstanceType(d.Get("host_instance_type").(string))
	if err != nil {
		return nil, err
	}
	var storageCapacityConverted int64
	storageCapacity := d.Get("storage_capacity").(string)
	if len(strings.TrimSpace(storageCapacity)) > 0 {
		storageCapacityConverted = ConvertStorageCapacitytoInt(storageCapacity)
	}
	msftLicensingConfig := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
	return &model.ClusterConfig{
		NumHosts:          numHosts,
		HostCpuCoresCount: &hostCPUCoresCount,
		HostInstanceType:  &hostInstanceType,
		StorageCapacity:   &storageCapacityConverted,
		MsftLicenseConfig: msftLicensingConfig,
	}, nil
}
