// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/nsx_vmc_app/infra/external"
	nsx_vmc_appModel "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/nsx_vmc_app/model"
	autoscalercluster "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api/orgs/sddcs/clusters"
	autoscalermodel "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs/clusters/msft_licensing"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	task "github.com/vmware/terraform-provider-vmc/vmc/task"
)

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(300 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},
		Schema: sddcSchema(),
	}
}

// sddcSchema this helper function extracts the creation of the SDDC schema, so that
// it's made available for mocking in tests.
func sddcSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"sddc_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
		},
		"size": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  constants.MediumSddcSize,
			ValidateFunc: validation.StringInSlice([]string{
				constants.MediumSddcSize, constants.CapitalMediumSddcSize, constants.LargeSddcSize, constants.CapitalLargeSddcSize}, false),
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
			Description:  "The amount of hosts in the primary cluster of the SDDC",
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
			Default:  constants.AwsProviderType,
			ValidateFunc: validation.StringInSlice([]string{
				constants.AwsProviderType, constants.ZeroCloudProviderType}, false),
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
			Default:  constants.SingleAvailabilityZone,
			ValidateFunc: validation.StringInSlice([]string{
				constants.SingleAvailabilityZone, constants.MultiAvailabilityZone,
			}, false),
		},
		"region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.All(
				validation.NoZeroValues,
			),
			DiffSuppressFunc: func(_, o, n string, _ *schema.ResourceData) bool {
				return o == strings.ReplaceAll(strings.ToUpper(n), "-", "_")
			},
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
				[]string{constants.HostInstancetypeI3, constants.HostInstancetypeI3EN, constants.HostInstancetypeI4I, constants.HostInstancetypeC6I, constants.HostInstancetypeM7i24xl, constants.HostInstancetypeM7i48xl}, false),
		},
		"edrs_policy_type": {
			Type: schema.TypeString,
			// Exact value known after create
			Optional: true,
			Computed: true,
			ValidateFunc: validation.StringInSlice(
				[]string{constants.StorageScaleUpPolicyType, constants.CostPolicyType, constants.PerformancePolicyType, constants.RapidScaleUpPolicyType}, false),
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
			ValidateFunc: validation.IntBetween(constants.MinHosts, constants.MaxHosts),
			Description:  "The minimum number of hosts that the cluster can scale in to.",
		},
		"max_hosts": {
			Type: schema.TypeInt,
			// Exact value known after create
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.IntBetween(constants.MinHosts, constants.MaxHosts),
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
							constants.LicenseConfigEnabled, constants.LicenseConfigDisabled, constants.CapitalLicenseConfigEnabled, constants.CapitalLicenseConfigDisabled}, false),
					},
					"windows_licensing": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "The status of Windows licensing for this SDDC's clusters. Possible values : enabled, ENABLED, disabled, DISABLED.",
						ValidateFunc: validation.StringInSlice([]string{
							constants.LicenseConfigEnabled, constants.LicenseConfigDisabled, constants.CapitalLicenseConfigEnabled, constants.CapitalLicenseConfigDisabled}, false),
					},
					"academic_license": {
						Type:        schema.TypeBool,
						Optional:    true,
						Description: "Flag to identify if it is Academic Standard or Commercial Standard License.",
					},
				},
			},
			Optional:    true,
			Description: "Indicates the desired licensing support, if any, of Microsoft software.",
		},
		"intranet_mtu_uplink": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      constants.MinIntranetMtuLink,
			Description:  "Uplink MTU of direct connect, SDDC-grouping and outposts traffic in edge tier-0 router port.",
			ValidateFunc: validation.IntBetween(constants.MinIntranetMtuLink, constants.MaxIntranetMtuLink),
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
		"nsxt_ui": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"nsxt_cloudadmin": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nsxt_cloudadmin_password": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nsxt_cloudaudit": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nsxt_cloudaudit_password": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nsxt_private_ip": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nsxt_private_url": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		// Following properties are set in the resourceSddcRead function and need to be
		// present in the schema during E2E test result validation phase
		"updated": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_by_user_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"version": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"updated_by_user_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"user_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"org_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"account_link_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"sddc_access_state": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	sddcClient := orgs.NewSddcsClient(connectorWrapper)
	orgID := connectorWrapper.OrgID

	var awsSddcConfig, err = buildAwsSddcConfig(d)
	if err != nil {
		return err
	}

	// Create a Sddc
	sddcCreateTask, err := sddcClient.Create(orgID, *awsSddcConfig, nil)
	if err != nil {
		return HandleCreateError("SDDC", err)
	}

	sddcID := sddcCreateTask.ResourceId
	d.SetId(*sddcID)
	msftLicensingConfig := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))

	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetTask(connectorWrapper, sddcCreateTask.Id)
		}, "error creating SDDC", nil)
		if taskErr != nil {
			return taskErr
		}
		err = resourceSddcRead(d, m)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		// Updating the microsoft_license_config after creation since
		// the backend API throws an error when non nil microsoft_licensing_config
		// is present in the sddc spec
		if msftLicensingConfig != nil {
			err = updateMsftLicenseConfig(d, m, msftLicensingConfig)
			if err != nil {
				return retry.NonRetryableError(err)
			}
		}

		return nil
	})
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	sddcID := d.Id()
	orgID := (m.(*connector.Wrapper)).OrgID
	sddc, err := GetSddc(connectorWrapper.Connector, orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "SDDC", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Unable to retrieve SDDC with ID %s", sddc.Id)
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)

	if err := d.Set("sddc_name", sddc.Name); err != nil {
		return err
	}
	if err := d.Set("updated", sddc.Updated.String()); err != nil {
		return err
	}
	if err := d.Set("user_id", sddc.UserId); err != nil {
		return err
	}
	if err := d.Set("updated_by_user_id", sddc.UpdatedByUserId); err != nil {
		return err
	}
	if err := d.Set("created", sddc.Created.String()); err != nil {
		return err
	}
	if err := d.Set("version", sddc.Version); err != nil {
		return err
	}
	if err := d.Set("updated_by_user_name", sddc.UpdatedByUserName); err != nil {
		return err
	}
	if err := d.Set("user_name", sddc.UserName); err != nil {
		return err
	}
	if err := d.Set("org_id", sddc.OrgId); err != nil {
		return err
	}
	if err := d.Set("sddc_type", sddc.SddcType); err != nil {
		return err
	}
	if err := d.Set("account_link_state", sddc.AccountLinkState); err != nil {
		return err
	}
	if err := d.Set("sddc_access_state", sddc.SddcAccessState); err != nil {
		return err
	}
	if err := d.Set("sddc_state", sddc.SddcState); err != nil {
		return err
	}
	primaryClusterClient := sddcs.NewPrimaryclusterClient(connectorWrapper.Connector)
	primaryCluster, err := primaryClusterClient.Get(orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "Primary Cluster", sddcID, err)
	}
	if err := d.Set("cluster_id", primaryCluster.ClusterId); err != nil {
		return err
	}
	cluster := map[string]string{}
	cluster["cluster_name"] = *primaryCluster.ClusterName
	cluster["cluster_state"] = *primaryCluster.ClusterState
	if primaryCluster.EsxHostInfo != nil {
		cluster["host_instance_type"] = *primaryCluster.EsxHostInfo.InstanceType
	}
	if primaryCluster.MsftLicenseConfig != nil {
		if primaryCluster.MsftLicenseConfig.MssqlLicensing != nil {
			cluster["mssql_licensing"] = *primaryCluster.MsftLicenseConfig.MssqlLicensing
		}
		if primaryCluster.MsftLicenseConfig.WindowsLicensing != nil {
			cluster["windows_licensing"] = *primaryCluster.MsftLicenseConfig.WindowsLicensing
		}
	}
	if err := d.Set("cluster_info", cluster); err != nil {
		return err
	}
	if sddc.ResourceConfig != nil {
		if err := d.Set("vc_url", sddc.ResourceConfig.VcUrl); err != nil {
			return err
		}
		if err := d.Set("cloud_username", sddc.ResourceConfig.CloudUsername); err != nil {
			return err
		}
		if err := d.Set("cloud_password", sddc.ResourceConfig.CloudPassword); err != nil {
			return err
		}
		if err := d.Set("nsxt_reverse_proxy_url", sddc.ResourceConfig.NsxApiPublicEndpointUrl); err != nil {
			return err
		}
		if err := d.Set("region", *sddc.ResourceConfig.Region); err != nil {
			return err
		}
		if err := d.Set("availability_zones", sddc.ResourceConfig.AvailabilityZones); err != nil {
			return err
		}
		if err := d.Set("deployment_type", ConvertDeployType(*sddc.ResourceConfig.DeploymentType)); err != nil {
			return err
		}
		if err := d.Set("sso_domain", *sddc.ResourceConfig.SsoDomain); err != nil {
			return err
		}
		if err := d.Set("skip_creating_vxlan", *sddc.ResourceConfig.SkipCreatingVxlan); err != nil {
			return err
		}
		if err := d.Set("provider_type", sddc.ResourceConfig.Provider); err != nil {
			return err
		}
		// SDDC's num_host should account for the amount of hosts on its primary cluster only.
		// Otherwise, there will be no way to scale up or down the primary cluster.
		if err := d.Set("num_host", getHostCountCluster(&sddc, primaryCluster.ClusterId)); err != nil {
			return err
		}
		if sddc.ResourceConfig.VpcInfo != nil && sddc.ResourceConfig.VpcInfo.VpcCidr != nil {
			if err := d.Set("vpc_cidr", *sddc.ResourceConfig.VpcInfo.VpcCidr); err != nil {
				return err
			}
		}
		skipCreatingVxLan := *sddc.ResourceConfig.SkipCreatingVxlan
		if !skipCreatingVxLan {
			if err := d.Set("vxlan_subnet", sddc.ResourceConfig.VxlanSubnet); err != nil {
				return err
			}
		}
		sddcSizeInfo := map[string]string{}
		sddcSizeInfo["vc_size"] = *sddc.ResourceConfig.SddcSize.VcSize
		sddcSizeInfo["nsx_size"] = *sddc.ResourceConfig.SddcSize.NsxSize
		if err := d.Set("sddc_size", sddcSizeInfo); err != nil {
			return err
		}
		if sddc.ResourceConfig.NsxCloudAdmin != nil {
			if err := d.Set("nsxt_cloudadmin", *sddc.ResourceConfig.NsxCloudAdmin); err != nil {
				return err
			}
			// Evade nil pointer dereference when user's access_token doesn't have NSX roles
			if sddc.ResourceConfig.NsxCloudAdminPassword != nil {
				if err := d.Set("nsxt_cloudadmin_password", *sddc.ResourceConfig.NsxCloudAdminPassword); err != nil {
					return err
				}
			}
			if sddc.ResourceConfig.NsxCloudAuditPassword != nil {
				if err := d.Set("nsxt_cloudaudit_password", *sddc.ResourceConfig.NsxCloudAuditPassword); err != nil {
					return err
				}
			}
			if err := d.Set("nsxt_cloudaudit", *sddc.ResourceConfig.NsxCloudAudit); err != nil {
				return err
			}
			if err := d.Set("nsxt_private_ip", *sddc.ResourceConfig.NsxMgrManagementIp); err != nil {
				return err
			}
			if err := d.Set("nsxt_private_url", *sddc.ResourceConfig.NsxMgrLoginUrl); err != nil {
				return err
			}
		}
	}
	edrsPolicyClient := autoscalercluster.NewEdrsPolicyClient(connectorWrapper.Connector)
	edrsPolicy, err := edrsPolicyClient.Get(orgID, sddcID, primaryCluster.ClusterId)
	if err != nil {
		return HandleReadError(d, "SDDC", sddcID, err)
	}
	if err := d.Set("edrs_policy_type", *edrsPolicy.PolicyType); err != nil {
		return err
	}
	if err := d.Set("enable_edrs", edrsPolicy.EnableEdrs); err != nil {
		return err
	}
	if err := d.Set("max_hosts", *edrsPolicy.MaxHosts); err != nil {
		return err
	}
	if err := d.Set("min_hosts", *edrsPolicy.MinHosts); err != nil {
		return err
	}

	if *sddc.Provider != constants.ZeroCloudProviderType {
		// store intranet_mtu_uplink only for non zerocloud provider types
		nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
		nsxtReverseProxyURLConnector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
		if err != nil {
			return HandleCreateError("NSXT reverse proxy URL connectorWrapper", err)
		}
		cloudServicesCommonClient := external.NewConfigClient(nsxtReverseProxyURLConnector)
		externalConnectivityConfig, err := cloudServicesCommonClient.Get()
		if err != nil {
			return HandleReadError(d, "External connectivity configuration", sddcID, err)
		}
		if err := d.Set("intranet_mtu_uplink", externalConnectivityConfig.IntranetMtu); err != nil {
			return err
		}
	}
	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	sddcClient := orgs.NewSddcsClient(connectorWrapper.Connector)
	sddcID := d.Id()
	orgID := (m.(*connector.Wrapper)).OrgID

	sddcDeleteTask, err := sddcClient.Delete(orgID, sddcID, nil, nil, nil)
	if err != nil {
		return HandleDeleteError("SDDC", sddcID, err)
	}
	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetTask(connectorWrapper, sddcDeleteTask.Id)
		}, "failed to delete SDDC", nil)
		if taskErr != nil {
			return taskErr
		}
		d.SetId("")
		return nil
	})
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	esxsClient := sddcs.NewEsxsClient(connectorWrapper)
	sddcClient := orgs.NewSddcsClient(connectorWrapper)
	sddcID := d.Id()
	orgID := (m.(*connector.Wrapper)).OrgID

	// Convert SDDC from 1NODE to DEFAULT
	if d.HasChange("sddc_type") {
		oldTmp, newTmp := d.GetChange("sddc_type")
		oldType := oldTmp.(string)
		newType := newTmp.(string)

		// Validate for convert type params
		if oldType == "1NODE" && (newType == "" || newType == "DEFAULT") {
			_, newTmp := d.GetChange("num_host")
			newNum := newTmp.(int)

			switch newNum {
			case 2: // 2node SDDC creation
				err := resourceSddcDelete(d, m)
				if err != nil {
					return err
				}
				return resourceSddcCreate(d, m)
			case 3: // 3node SDDC scale up
				convertClient := sddcs.NewConvertClient(connectorWrapper)
				sddcTypeUpdateTask, err := convertClient.Create(orgID, sddcID, nil)

				if err != nil {
					return HandleUpdateError("SDDC", err)
				}
				err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
					taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
						return task.GetTask(connectorWrapper, sddcTypeUpdateTask.Id)
					}, "error scaling SDDC", nil)
					if taskErr != nil {
						return taskErr
					}
					err = resourceSddcRead(d, m)
					if err == nil {
						return nil
					}
					return retry.NonRetryableError(err)
				})
				if err != nil {
					return err
				}
			default:
				return fmt.Errorf("scaling SDDC is not supported. Please check sddc_type and num_host")
			}
		}
	}

	// Add,remove hosts
	if d.HasChange("num_host") {
		primaryClusterID := d.Get("cluster_id").(string)
		oldTmp, newTmp := d.GetChange("num_host")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		if len(primaryClusterID) == 0 {
			return fmt.Errorf("cannot find primary cluster on SDDC %s", sddcID)
		}
		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}
		if d.Get("deployment_type").(string) == constants.MultiAvailabilityZone && diffNum%2 != 0 {

			return fmt.Errorf("for multiAZ deployment type, SDDC hosts must be added in pairs across availability zones")
		}
		esxConfig := model.EsxConfig{
			NumHosts:  int64(diffNum),
			ClusterId: &primaryClusterID,
		}

		hostUpdateTask, err := esxsClient.Create(orgID, sddcID, esxConfig, &action)

		if err != nil {
			return HandleUpdateError("SDDC", err)
		}
		err = retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
				return task.GetTask(connectorWrapper, hostUpdateTask.Id)
			}, "failed to update hosts", nil)
			if taskErr != nil {
				return taskErr
			}
			err = resourceSddcRead(d, m)
			if err == nil {
				return nil
			}
			return retry.NonRetryableError(err)
		})
		if err != nil {
			return err
		}
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
		if err := d.Set("sddc_name", sddc.Name); err != nil {
			return err
		}
	}

	if d.HasChange("intranet_mtu_uplink") {
		if d.Get("provider_type") == constants.ZeroCloudProviderType {
			return fmt.Errorf("intranet MTU uplink cannot be updated for %s provider type", constants.ZeroCloudProviderType)
		}
		intranetMTUUplink := d.Get("intranet_mtu_uplink").(int)
		intranetMTUUplinkPointer := int64(intranetMTUUplink)
		nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
		nxstReverseProxyURLConnector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
		if err != nil {
			return HandleCreateError("NSXT reverse proxy URL connector", err)
		}
		cloudServicesCommonClient := external.NewConfigClient(nxstReverseProxyURLConnector)
		externalConnectivityConfig := nsx_vmc_appModel.ExternalConnectivityConfig{IntranetMtu: &intranetMTUUplinkPointer}
		_, err = cloudServicesCommonClient.Update(externalConnectivityConfig)
		if err != nil {
			return HandleUpdateError("Intranet MTU Uplink", err)
		}
	}

	if d.HasChange("edrs_policy_type") || d.HasChange("enable_edrs") || d.HasChange("min_hosts") || d.HasChange("max_hosts") {
		sddcType := d.Get("sddc_type").(string)
		if sddcType == constants.OneNodeSddcType {
			return fmt.Errorf("EDRS policy cannot be updated for SDDC with type %s", constants.OneNodeSddcType)
		}
		clusterID := d.Get("cluster_id").(string)
		minHosts := int64(d.Get("min_hosts").(int))
		maxHosts := int64(d.Get("max_hosts").(int))
		policyType := d.Get("edrs_policy_type").(string)
		enableEDRS := d.Get("enable_edrs").(bool)
		if policyType == constants.StorageScaleUpPolicyType && !enableEDRS {
			return fmt.Errorf("EDRS policy %s is the default and cannot be disabled", constants.StorageScaleUpPolicyType)
		}
		edrsPolicy := &autoscalermodel.EdrsPolicy{
			EnableEdrs: enableEDRS,
			PolicyType: &policyType,
			MinHosts:   &minHosts,
			MaxHosts:   &maxHosts,
		}
		edrsPolicyClient := autoscalercluster.NewEdrsPolicyClient(connectorWrapper)
		edrsPolicyUpdateTask, err := edrsPolicyClient.Post(orgID, sddcID, clusterID, *edrsPolicy)
		if err != nil {
			return HandleUpdateError("EDRS Policy", err)
		}

		return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
			taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
				return task.GetTask(connectorWrapper, edrsPolicyUpdateTask.Id)
			}, "failed to update EDRS policy configuration", nil)
			if taskErr != nil {
				return taskErr
			}
			err = resourceSddcRead(d, m)
			if err == nil {
				return nil
			}
			return retry.NonRetryableError(err)
		})
	}

	// Update sddc_size is not supported
	if d.HasChange("size") {
		return fmt.Errorf("SDDC size update operation is not supported")
	}

	// Update Microsoft licensing config
	if d.HasChange("microsoft_licensing_config") {
		configChangeParam := expandMsftLicenseConfig(d.Get("microsoft_licensing_config").([]interface{}))
		return updateMsftLicenseConfig(d, m, configChangeParam)
	}
	return resourceSddcRead(d, m)
}

func updateMsftLicenseConfig(d *schema.ResourceData, m interface{}, msftLicenseConfig *model.MsftLicensingConfig) error {
	connectorWrapper := m.(*connector.Wrapper)
	sddcID := d.Id()
	orgID := (m.(*connector.Wrapper)).OrgID
	primaryClusterClient := sddcs.NewPrimaryclusterClient(connectorWrapper)
	primaryCluster, err := primaryClusterClient.Get(orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "Primary Cluster", sddcID, err)
	}
	publishClient := msft_licensing.NewPublishClient(connectorWrapper)
	microsoftLicensingUpdateTask, err := publishClient.Post(orgID, sddcID, primaryCluster.ClusterId, *msftLicenseConfig)
	if err != nil {
		return fmt.Errorf("error updating license : %s", err)
	}
	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetTask(connectorWrapper, microsoftLicensingUpdateTask.Id)
		}, "failed updating Microsoft licensing configuration", nil)
		if taskErr != nil {
			return taskErr
		}
		err = resourceSddcRead(d, m)
		if err == nil {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}

// buildAwsSddcConfig extracts the creation of the model.AwsSddcConfig, so that it's
// available for testing
func buildAwsSddcConfig(d *schema.ResourceData) (*model.AwsSddcConfig, error) {
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

	var c map[string]interface{}
	accountLinkSddcConfigVar := d.Get("account_link_sddc_config").([]interface{})
	for _, config := range accountLinkSddcConfigVar {
		c = config.(map[string]interface{})
	}

	if deploymentType == constants.MultiAvailabilityZone && c != nil && len(c["customer_subnet_ids"].([]interface{})) != 2 {
		return nil, fmt.Errorf("deployment type %s requires 2 subnet IDs, one in each availability zone ", deploymentType)
	}

	if deploymentType == constants.SingleAvailabilityZone && c != nil && len(c["customer_subnet_ids"].([]interface{})) != 1 {
		return nil, fmt.Errorf("deployment type %s requires 1 subnet ID ", deploymentType)
	}

	accountLinkSddcConfig := expandAccountLinkSddcConfig(accountLinkSddcConfigVar)
	dataHostInstanceType := d.Get("host_instance_type").(string)
	hostInstanceType, err := toHostInstanceType(dataHostInstanceType)
	if len(dataHostInstanceType) > 0 && err != nil {
		// return error only if nonempty host_instance_type is provided
		// since host_instance_type field is optional in schema
		return nil, err
	}

	model := model.AwsSddcConfig{
		Name:                  sddcName,
		VpcCidr:               &vpcCidr,
		NumHosts:              int64(numHost),
		SddcType:              sddcTypePtr,
		VxlanSubnet:           &vxlanSubnet,
		Provider:              providerType,
		SkipCreatingVxlan:     &skipCreatingVxlan,
		AccountLinkConfig:     accountLinkConfig,
		AccountLinkSddcConfig: accountLinkSddcConfig,
		SsoDomain:             &ssoDomain,
		SddcTemplateId:        &sddcTemplateID,
		DeploymentType:        &deploymentType,
		Region:                &region,
		HostInstanceType:      &hostInstanceType,
		Size:                  &sddcSize,
		MsftLicenseConfig:     nil,
	}

	return &model, nil
}

func expandAccountLinkSddcConfig(l []interface{}) []model.AccountLinkSddcConfig {

	if len(l) == 0 {
		return nil
	}

	var configs []model.AccountLinkSddcConfig

	for _, config := range l {
		c := config.(map[string]interface{})
		var subnetIDs []string
		for _, subnetID := range c["customer_subnet_ids"].([]interface{}) {
			subnetIDs = append(subnetIDs, subnetID.(string))
		}
		var connectedAccID = c["connected_account_id"].(string)
		con := model.AccountLinkSddcConfig{
			CustomerSubnetIds:  subnetIDs,
			ConnectedAccountId: &connectedAccID,
		}

		configs = append(configs, con)
	}
	return configs
}
