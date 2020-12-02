/* Copyright 2019 VMware, Inc.
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

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(300 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"storage_capacity": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"15TB", "20TB", "25TB", "30TB", "35TB"}, false),
			},
			"sddc_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"size": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  MediumSDDCSize,
				ValidateFunc: validation.StringInSlice([]string{
					MediumSDDCSize, CapitalMediumSDDCSize, LargeSDDCSize, CapitalLargeSDDCSize}, false),
				Description: "The size of the vCenter and NSX appliances. 'large' or 'LARGE' SDDC size corresponds to a large vCenter appliance and large NSX appliance. 'medium' or 'MEDIUM' SDDC size corresponds to medium vCenter appliance and medium NSX appliance. Default : 'medium'.",
			},
			"account_link_sddc_config": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_subnet_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
						"connected_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
				ForceNew: true,
			},
			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"num_host": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"sddc_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vxlan_subnet": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"delay_account_link": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "AWS",
				ValidateFunc: validation.StringInSlice([]string{
					"AWS", "ZEROCLOUD"}, false),
			},
			"skip_creating_vxlan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"sso_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "vmc.local",
			},
			"sddc_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  SingleAvailabilityZone,
				ValidateFunc: validation.StringInSlice([]string{
					SingleAvailabilityZone, MultiAvailabilityZone,
				}, false),
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.NoZeroValues,
				),
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{HostInstancetypeI3, HostInstancetypeR5, HostInstancetypeI3EN}, false),
			},
			"edrs_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice(
					[]string{StorageScaleUpPolicyType, CostPolicyType, PerformancePolicyType, RapidScaleUpPolicyType}, false),
				Description: "The EDRS policy type. This can either be 'cost', 'performance', 'storage-scaleup' or 'rapid-scaleup'. Default : storage-scaleup. ",
			},
			"enable_edrs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "True if EDRS is enabled",
			},
			"min_hosts": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 16),
				Description:  "The minimum number of hosts that the cluster can scale in to.",
			},
			"max_hosts": {
				Type:         schema.TypeInt,
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
			"sddc_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vc_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cloud_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_reverse_proxy_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_info": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"sddc_size": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
		CustomizeDiff: func(d *schema.ResourceDiff, meta interface{}) error {

			deploymentType := d.Get("deployment_type").(string)
			numHosts := d.Get("num_host").(int)

			if deploymentType == MultiAvailabilityZone && numHosts < MinMultiAZHosts {
				return fmt.Errorf("for MulitAZ deployment type number of hosts must be atleast %d ", MinMultiAZHosts)
			}

			var c map[string]interface{}
			accountLinkSddcConfig := d.Get("account_link_sddc_config").([]interface{})
			for _, config := range accountLinkSddcConfig {
				c = config.(map[string]interface{})
			}

			if deploymentType == MultiAvailabilityZone && c != nil && len(c["customer_subnet_ids"].([]interface{})) != 2 {
				return fmt.Errorf("deployment type %s requires 2 subnet IDs, one in each availability zone ", deploymentType)
			}

			if deploymentType == SingleAvailabilityZone && c != nil && len(c["customer_subnet_ids"].([]interface{})) != 1 {
				return fmt.Errorf("deployment type %s requires 1 subnet ID ", deploymentType)
			}

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

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	var storageCapacityConverted int64
	connectorWrapper := m.(*ConnectorWrapper)
	sddcClient := orgs.NewDefaultSddcsClient(connectorWrapper)
	orgID := connectorWrapper.OrgID

	storageCapacity := d.Get("storage_capacity").(string)
	if len(strings.TrimSpace(storageCapacity)) > 0 {
		storageCapacityConverted = ConvertStorageCapacitytoInt(storageCapacity)
	}

	sddcName := d.Get("sddc_name").(string)
	vpcCidr := d.Get("vpc_cidr").(string)
	numHost := d.Get("num_host").(int)
	sddcType := d.Get("sddc_type").(string)
	sddcSize := d.Get("size").(string)
	var sddcTypePtr *string
	if sddcType != "" {
		sddcTypePtr = &sddcType
	}
	vxlanSubnet := d.Get("vxlan_subnet").(string)
	delayAccountLink := d.Get("delay_account_link").(bool)
	accountLinkConfig := &model.AccountLinkConfig{
		DelayAccountLink: &delayAccountLink,
	}
	providerType := d.Get("provider_type").(string)
	skipCreatingVxlan := d.Get("skip_creating_vxlan").(bool)
	ssoDomain := d.Get("sso_domain").(string)
	sddcTemplateID := d.Get("sddc_template_id").(string)
	deploymentType := d.Get("deployment_type").(string)
	region := d.Get("region").(string)
	accountLinkSddcConfig := expandAccountLinkSddcConfig(d.Get("account_link_sddc_config").([]interface{}))
	hostInstanceType := model.HostInstanceTypes(d.Get("host_instance_type").(string))
	msftLicensingConfig := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))

	var awsSddcConfig = &model.AwsSddcConfig{
		StorageCapacity:       &storageCapacityConverted,
		Name:                  sddcName,
		VpcCidr:               &vpcCidr,
		NumHosts:              int64(numHost),
		SddcType:              sddcTypePtr,
		VxlanSubnet:           &vxlanSubnet,
		AccountLinkConfig:     accountLinkConfig,
		Provider:              providerType,
		SkipCreatingVxlan:     &skipCreatingVxlan,
		AccountLinkSddcConfig: accountLinkSddcConfig,
		SsoDomain:             &ssoDomain,
		SddcTemplateId:        &sddcTemplateID,
		DeploymentType:        &deploymentType,
		Region:                region,
		HostInstanceType:      &hostInstanceType,
		Size:                  &sddcSize,
		MsftLicenseConfig:     msftLicensingConfig,
	}

	// Create a Sddc
	task, err := sddcClient.Create(orgID, *awsSddcConfig, nil)
	if err != nil {
		return HandleCreateError("SDDC", err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	d.SetId(*sddcID)
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Printf("Authentication error : %v", errors.Unauthenticated{}.Error())
				err = connectorWrapper.authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %v", err))
				}
				return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("error creating SDDC : %v", err))

		}
		if *task.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("task failed to create SDDC"))
		} else if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected SDDC to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSddcRead(d, m))
	})
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcID := d.Id()
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddc, err := GetSDDC(connector, orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "SDDC", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Unable to retrieve SDDC with ID %s", sddc.Id)
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)

	d.Set("sddc_name", sddc.Name)
	d.Set("updated", sddc.Updated)
	d.Set("user_id", sddc.UserId)
	d.Set("updated_by_user_id", sddc.UpdatedByUserId)
	d.Set("created", sddc.Created)
	d.Set("version", sddc.Version)
	d.Set("updated_by_user_name", sddc.UpdatedByUserName)
	d.Set("user_name", sddc.UserName)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("provider", sddc.Provider)
	d.Set("account_link_state", sddc.AccountLinkState)
	d.Set("sddc_access_state", sddc.SddcAccessState)
	d.Set("sddc_state", sddc.SddcState)
	if sddc.ResourceConfig != nil {
		d.Set("vc_url", sddc.ResourceConfig.VcUrl)
		d.Set("cloud_username", sddc.ResourceConfig.CloudUsername)
		d.Set("cloud_password", sddc.ResourceConfig.CloudPassword)
		d.Set("nsxt_reverse_proxy_url", sddc.ResourceConfig.NsxApiPublicEndpointUrl)
		d.Set("region", *sddc.ResourceConfig.Region)
		d.Set("availability_zones", sddc.ResourceConfig.AvailabilityZones)
		d.Set("deployment_type", ConvertDeployType(*sddc.ResourceConfig.DeploymentType))
		d.Set("sso_domain", *sddc.ResourceConfig.SsoDomain)
		d.Set("skip_creating_vxlan", *sddc.ResourceConfig.SkipCreatingVxlan)
		d.Set("provider_type", sddc.ResourceConfig.Provider)
		d.Set("num_host", len(sddc.ResourceConfig.EsxHosts))
		d.Set("vpc_cidr", *sddc.ResourceConfig.VpcInfo.VpcCidr)
		d.Set("vxlan_subnet", sddc.ResourceConfig.VxlanSubnet)

		sddcSizeInfo := map[string]string{}
		sddcSizeInfo["vc_size"] = *sddc.ResourceConfig.SddcSize.VcSize
		sddcSizeInfo["nsx_size"] = *sddc.ResourceConfig.SddcSize.NsxSize
		d.Set("sddc_size", sddcSizeInfo)
	}
	sddcClient := sddcs.NewDefaultPrimaryclusterClient(connector)
	primaryCluster, err := sddcClient.Get(orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "Primary Cluster", sddcID, err)
	}
	d.Set("cluster_id", primaryCluster.ClusterId)
	cluster := map[string]string{}
	cluster["cluster_name"] = *primaryCluster.ClusterName
	cluster["cluster_state"] = *primaryCluster.ClusterState
	cluster["host_instance_type"] = *primaryCluster.EsxHostInfo.InstanceType
	if primaryCluster.MsftLicenseConfig != nil {
		cluster["mssql_licensing"] = *primaryCluster.MsftLicenseConfig.MssqlLicensing
		cluster["windows_licensing"] = *primaryCluster.MsftLicenseConfig.WindowsLicensing
	}
	d.Set("cluster_info", cluster)

	edrsPolicyClient := autoscalercluster.NewDefaultEdrsPolicyClient(connector)
	edrsPolicy, err := edrsPolicyClient.Get(orgID, sddcID, primaryCluster.ClusterId)
	if err != nil {
		return HandleReadError(d, "SDDC", sddcID, err)
	}
	d.Set("edrs_policy_type", *edrsPolicy.PolicyType)
	d.Set("enable_edrs", edrsPolicy.EnableEdrs)
	d.Set("max_hosts", *edrsPolicy.MaxHosts)
	d.Set("min_hosts", *edrsPolicy.MinHosts)

	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	sddcID := d.Id()
	orgID := (m.(*ConnectorWrapper)).OrgID

	task, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
	if err != nil {
		return HandleDeleteError("SDDC", sddcID, err)
	}
	tasksClient := orgs.NewDefaultTasksClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error deleting SDDC %s: %v", sddcID, err))
		}
		if *task.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("task failed to delete SDDC"))
		} else if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected SDDC to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*ConnectorWrapper)
	esxsClient := sddcs.NewDefaultEsxsClient(connectorWrapper)
	sddcClient := orgs.NewDefaultSddcsClient(connectorWrapper)
	sddcID := d.Id()
	orgID := (m.(*ConnectorWrapper)).OrgID

	// Convert SDDC from 1NODE to DEFAULT
	if d.HasChange("sddc_type") {
		oldTmp, newTmp := d.GetChange("sddc_type")
		oldType := oldTmp.(string)
		newType := newTmp.(string)

		// Validate for convert type params
		if oldType == "1NODE" && (newType == "" || newType == "DEFAULT") {
			_, newTmp := d.GetChange("num_host")
			newNum := newTmp.(int)

			if newNum == 2 { // 2node SDDC creation
				err := resourceSddcDelete(d, m)
				if err != nil {
					return err
				}
				return resourceSddcCreate(d, m)
			} else if newNum == 3 { // 3node SDDC scale up
				convertClient := sddcs.NewDefaultConvertClient(connectorWrapper)
				task, err := convertClient.Create(orgID, sddcID)

				if err != nil {
					return HandleUpdateError("SDDC", err)
				}
				err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
					tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
					task, err := tasksClient.Get(orgID, task.Id)

					if err != nil {
						if err.Error() == (errors.Unauthenticated{}.Error()) {
							log.Printf("Authentication error : %v", errors.Unauthenticated{}.Error())
							err = connectorWrapper.authenticate()
							if err != nil {
								return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
							}
							return resource.RetryableError(fmt.Errorf("SDDC scaling still in progress"))
						}
						return resource.NonRetryableError(fmt.Errorf("error scaling SDDC : %v", err))
					}
					if *task.Status == "FAILED" {
						return resource.NonRetryableError(fmt.Errorf("task failed to scale SDDC"))
					} else if *task.Status != "FINISHED" {
						return resource.RetryableError(fmt.Errorf("expected SDDC type to be updated but was in state %s", *task.Status))
					}
					return resource.NonRetryableError(resourceSddcRead(d, m))
				})
				if err != nil {
					return err
				}
			} else {
				return fmt.Errorf("scaling SDDC is not supported. Please check sddc_type and num_host")
			}
		}
		return resourceSddcRead(d, m)
	}

	// Add,remove hosts
	if d.HasChange("num_host") {
		oldTmp, newTmp := d.GetChange("num_host")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}
		if d.Get("deployment_type").(string) == MultiAvailabilityZone && diffNum%2 != 0 {

			return fmt.Errorf("for multiAZ deployment type, SDDC hosts must be added in pairs across availability zones")
		}
		esxConfig := model.EsxConfig{
			NumHosts: int64(diffNum),
		}

		task, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)

		if err != nil {
			return HandleUpdateError("SDDC", err)
		}
		tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error updating hosts : %v", err))
			}
			if *task.Status == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("task failed to update hosts"))
			} else if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected hosts to be updated but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceSddcRead(d, m))
		})
		if err != nil {
			return err
		}
		return resourceSddcRead(d, m)
	}

	// Update sddc name
	if d.HasChange("sddc_name") {
		newSDDCName := d.Get("sddc_name").(string)
		sddcPatchRequest := model.SddcPatchRequest{
			Name: &newSDDCName,
		}
		sddc, err := sddcClient.Patch(orgID, sddcID, sddcPatchRequest)

		if err != nil {
			return HandleUpdateError("SDDC", err)
		}
		d.Set("sddc_name", sddc.Name)
	}

	if d.HasChange("edrs_policy_type") || d.HasChange("enable_edrs") || d.HasChange("min_hosts") || d.HasChange("max_hosts") {

		SDDCType := d.Get("sddc_type").(string)
		if SDDCType == OneNodeSDDCType {
			return fmt.Errorf("EDRS policy cannot be updated for SDDC with type %s", OneNodeSDDCType)
		}
		clusterID := d.Get("cluster_id").(string)
		minHosts := int64(d.Get("min_hosts").(int))
		maxHosts := int64(d.Get("max_hosts").(int))
		policyType := d.Get("edrs_policy_type").(string)
		enableEDRS := d.Get("enable_edrs").(bool)
		if policyType == StorageScaleUpPolicyType && !enableEDRS {
			return fmt.Errorf("EDRS policy %s is the default and cannot be disabled", StorageScaleUpPolicyType)
		}
		edrsPolicy := &autoscalermodel.EdrsPolicy{
			EnableEdrs: enableEDRS,
			PolicyType: &policyType,
			MinHosts:   &minHosts,
			MaxHosts:   &maxHosts,
		}
		edrsPolicyClient := autoscalercluster.NewDefaultEdrsPolicyClient(connectorWrapper)
		task, err := edrsPolicyClient.Post(orgID, sddcID, clusterID, *edrsPolicy)
		if err != nil {
			return HandleUpdateError("EDRS Policy", err)
		}

		return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			taskClient := autoscalerapi.NewDefaultAutoscalerClient(connectorWrapper)
			task, err := taskClient.Get(orgID, task.Id)
			if err != nil {
				if err.Error() == (errors.Unauthenticated{}.Error()) {
					log.Printf("Authentication error : %v", errors.Unauthenticated{}.Error())
					err = connectorWrapper.authenticate()
					if err != nil {
						return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
					}
					return resource.RetryableError(fmt.Errorf("instance update still in progress"))
				}
				return resource.NonRetryableError(fmt.Errorf("error updating EDRS policy configuration : %v", err))
			}
			if *task.Status == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("task failed to update EDRS policy configuration"))
			} else if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected EDRS policy configuration to be updated but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceSddcRead(d, m))
		})
	}
	if d.HasChange("size") {
		return fmt.Errorf("SDDC size update operation is not supported")
	}

	// Update microsoft licensing config
	if d.HasChange("microsoft_licensing_config") {
		configChangeParam := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
		primaryClusterClient := sddcs.NewDefaultPrimaryclusterClient(connectorWrapper)
		primaryCluster, err := primaryClusterClient.Get(orgID, sddcID)
		if err != nil {
			return HandleReadError(d, "Primary Cluster", sddcID, err)
		}
		publishClient := msft_licensing.NewDefaultPublishClient(connectorWrapper)
		task, err := publishClient.Post(orgID, sddcID, primaryCluster.ClusterId, *configChangeParam)
		if err != nil {
			return fmt.Errorf("Error updating license : %s", err)
		}
		return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			tasksClient := orgs.NewDefaultTasksClient(connectorWrapper)
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil || *task.Status == "FAILED" {
				if err.Error() == (errors.Unauthenticated{}.Error()) {
					log.Printf("Authentication error : %v", errors.Unauthenticated{}.Error())
					err = connectorWrapper.authenticate()
					if err != nil {
						return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %s", err))
					}
					return resource.RetryableError(fmt.Errorf("instance update still in progress"))
				}
				return resource.NonRetryableError(fmt.Errorf("error updating microsoft licensing configuration : %v", err))
			}
			if *task.Status == "FAILED" {
				return resource.NonRetryableError(fmt.Errorf("task failed to update microsoft licensing configuration"))
			} else if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected microsoft licensing configuration to be updated but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourceSddcRead(d, m))
		})

	}
	return resourceSddcRead(d, m)
}

func expandAccountLinkSddcConfig(l []interface{}) []model.AccountLinkSddcConfig {

	if len(l) == 0 {
		return nil
	}

	var configs []model.AccountLinkSddcConfig

	for _, config := range l {
		c := config.(map[string]interface{})
		var subnetIds []string
		for _, subnetID := range c["customer_subnet_ids"].([]interface{}) {
			subnetIds = append(subnetIds, subnetID.(string))
		}
		var connectedAccId = c["connected_account_id"].(string)
		con := model.AccountLinkSddcConfig{
			CustomerSubnetIds:  subnetIds,
			ConnectedAccountId: &connectedAccId,
		}

		configs = append(configs, con)
	}
	return configs
}
