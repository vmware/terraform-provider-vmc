/* Copyright 2021 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
)

func dataSourceVmcSddc() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcSddcRead,

		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				Description: "Sddc ID.",
				Required:    true,
			},
			"sddc_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"num_host": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sddc_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provider_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"deployment_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sddc_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skip_creating_vxlan": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"sso_domain": {
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
			"nsxt_reverse_proxy_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nsxt_ui": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nsxt_cloudadmin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_cloudadmin_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_cloudaudit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_cloudaudit_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nsxt_private_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Below are added as part of the schema as they are set in the
			// dataSourceVmcSddcRead method
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
		},
	}
}

func dataSourceVmcSddcRead(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := (m.(*connector.Wrapper)).Connector
	sddcClient := orgs.NewSddcsClient(connectorWrapper)
	sddcID := d.Get("sddc_id").(string)
	orgID := (m.(*connector.Wrapper)).OrgID
	sddc, err := sddcClient.Get(orgID, sddcID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("SDDC with ID %s not found", sddcID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error while getting the SDDC with ID %s,%v", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Can't get, SDDC with ID %s is already deleted", sddc.Id)
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)

	d.Set("sddc_name", sddc.Name)
	d.Set("updated", sddc.Updated.String())
	d.Set("user_id", sddc.UserId)
	d.Set("updated_by_user_id", sddc.UpdatedByUserId)
	d.Set("created", sddc.Created.String())
	d.Set("version", sddc.Version)
	d.Set("updated_by_user_name", sddc.UpdatedByUserName)
	d.Set("user_name", sddc.UserName)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_type", sddc.SddcType)
	// the key "provider" is reserved by the Terraform SDK, however the same information
	// is carried by the sddc.ResourceConfig.Provider variable
	//d.Set("provider", sddc.Provider)
	d.Set("account_link_state", sddc.AccountLinkState)
	d.Set("sddc_access_state", sddc.SddcAccessState)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("sddc_state", sddc.SddcState)
	if sddc.ResourceConfig != nil {
		d.Set("vc_url", sddc.ResourceConfig.VcUrl)
		d.Set("cloud_username", sddc.ResourceConfig.CloudUsername)
		d.Set("nsxt_reverse_proxy_url", sddc.ResourceConfig.NsxApiPublicEndpointUrl)
		d.Set("region", sddc.ResourceConfig.Region)
		// Query the API for primary Cluster ID so only it's hosts can be added to the
		// sddc host
		primaryClusterClient := sddcs.NewPrimaryclusterClient(connectorWrapper)
		primaryCluster, err := primaryClusterClient.Get(orgID, sddcID)
		if err != nil {
			return HandleReadError(d, "Primary Cluster", sddcID, err)
		}
		d.Set("num_host", getHostCountCluster(&sddc, primaryCluster.ClusterId))
		d.Set("provider_type", sddc.ResourceConfig.Provider)
		d.Set("availability_zones", sddc.ResourceConfig.AvailabilityZones)
		d.Set("deployment_type", ConvertDeployType(*sddc.ResourceConfig.DeploymentType))
		d.Set("sso_domain", *sddc.ResourceConfig.SsoDomain)
		d.Set("skip_creating_vxlan", *sddc.ResourceConfig.SkipCreatingVxlan)
		d.Set("nsxt_ui", *sddc.ResourceConfig.Nsxt)
		if sddc.ResourceConfig.NsxCloudAdmin != nil {
			d.Set("nsxt_cloudadmin", *sddc.ResourceConfig.NsxCloudAdmin)
			d.Set("nsxt_cloudadmin_password", *sddc.ResourceConfig.NsxCloudAdminPassword)
			d.Set("nsxt_cloudaudit", *sddc.ResourceConfig.NsxCloudAudit)
			d.Set("nsxt_cloudaudit_password", *sddc.ResourceConfig.NsxCloudAuditPassword)
			d.Set("nsxt_private_ip", *sddc.ResourceConfig.NsxMgrManagementIp)
			d.Set("nsxt_private_url", *sddc.ResourceConfig.NsxMgrLoginUrl)
		}
	}

	return nil
}
