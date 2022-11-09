/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"context"
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	task "github.com/vmware/terraform-provider-vmc/vmc/task"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	draasmodel "github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
)

func resourceSiteRecovery() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryCreate,
		Read:   resourceSiteRecoveryRead,
		Update: resourceSiteRecoveryUpdate,
		Delete: resourceSiteRecoveryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "SDDC identifier",
			},
			"srm_extension_key_suffix": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 13),
				Description:  "The custom extension suffix for SRM must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The suffix is appended to com.vmware.vcDr- to form the full extension key. ",
			},
			"site_recovery_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Site recovery state. Possible values are: ACTIVATED, ACTIVATING, CANCELED, DEACTIVATED, DEACTIVATING, DELETED, FAILED",
			},
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User id that last updated this record.",
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "User name that last updated this record.",
			},
			"srm_node": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"vr_node": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"draas_h5_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSiteRecoveryCreate(d *schema.ResourceData, m interface{}) error {

	err := (m.(*connector.Wrapper)).Authenticate()
	if err != nil {
		return fmt.Errorf("authentication error from Cloud Service Provider: %s", err)
	}
	connectorWrapper := (m.(*connector.Wrapper))

	siteRecoveryClient := draas.NewSiteRecoveryClient(connectorWrapper)

	srmExtensionKeySuffix := d.Get("srm_extension_key_suffix").(string)
	orgID := (m.(*connector.Wrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	activateSiteRecoveryConfigParam := &draasmodel.ActivateSiteRecoveryConfig{
		SrmExtensionKeySuffix: &srmExtensionKeySuffix,
	}

	siteRecoveryCreateTask, err := siteRecoveryClient.Post(orgID, sddcID, activateSiteRecoveryConfigParam)

	if err != nil {
		return HandleCreateError("Site recovery", err)
	}

	// Wait until site recovery is activated
	taskID := siteRecoveryCreateTask.ResourceId
	d.SetId(*taskID)
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return task.GetDraasTask(connectorWrapper, siteRecoveryCreateTask.Id)
			},
			"error activation site recovery ",
			nil)
		if taskErr != nil {
			return taskErr
		}
		err = resourceSiteRecoveryRead(d, m)
		if err == nil {
			return nil
		}
		return resource.NonRetryableError(err)
	})
}

func resourceSiteRecoveryRead(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := (m.(*connector.Wrapper)).Connector
	sddcID := d.Id()
	orgID := (m.(*connector.Wrapper)).OrgID
	siteRecoveryClient := draas.NewSiteRecoveryClient(connectorWrapper)
	siteRecovery, err := siteRecoveryClient.Get(orgID, sddcID)
	if err != nil {

		return HandleReadError(d, "Site recovery", sddcID, err)
	}
	d.SetId(siteRecovery.Id)
	d.Set("site_recovery_state", siteRecovery.SiteRecoveryState)
	d.Set("draas_h5_url", siteRecovery.DraasH5Url)
	d.Set("user_id", siteRecovery.UserId)
	d.Set("user_name", siteRecovery.UserName)

	srmExtensionKey := d.Get("srm_extension_key_suffix").(string)
	srmNodeMap := map[string]string{}
	for _, SRMNode := range siteRecovery.SrmNodes {
		if len(strings.TrimSpace(srmExtensionKey)) == 0 {
			tempStr := strings.Trim(*SRMNode.Hostname, ".")
			if strings.Contains(tempStr, "-") {
				srmNodeMap["id"] = *SRMNode.Id
				srmNodeMap["ip_address"] = *SRMNode.IpAddress
				srmNodeMap["host_name"] = *SRMNode.Hostname
				srmNodeMap["state"] = *SRMNode.State
				srmNodeMap["type"] = *SRMNode.Type_
				// During tests VmMorefId might be nil
				if SRMNode.VmMorefId != nil {
					srmNodeMap["vm_moref_id"] = *SRMNode.VmMorefId
				}
				break
			}
		} else if strings.Contains(*SRMNode.Hostname, strings.TrimSpace(srmExtensionKey)) {
			srmNodeMap["id"] = *SRMNode.Id
			srmNodeMap["ip_address"] = *SRMNode.IpAddress
			srmNodeMap["host_name"] = *SRMNode.Hostname
			srmNodeMap["state"] = *SRMNode.State
			srmNodeMap["type"] = *SRMNode.Type_
			// During tests VmMorefId might be nil
			if SRMNode.VmMorefId != nil {
				srmNodeMap["vm_moref_id"] = *SRMNode.VmMorefId
			}
			break
		}
	}

	vrNodeMap := map[string]string{}
	// During tests VmMorefId might be nil
	if siteRecovery.VrNode.VmMorefId != nil {
		vrNodeMap["vm_moref_id"] = *siteRecovery.VrNode.VmMorefId
	}
	vrNodeMap["id"] = *siteRecovery.VrNode.Id
	vrNodeMap["hostname"] = *siteRecovery.VrNode.Hostname
	vrNodeMap["type"] = *siteRecovery.VrNode.Type_
	vrNodeMap["state"] = *siteRecovery.VrNode.State
	vrNodeMap["ip_address"] = *siteRecovery.VrNode.IpAddress
	d.Set("sddc_id", *siteRecovery.SddcId)
	d.Set("srm_node", srmNodeMap)
	d.Set("vr_node", vrNodeMap)
	return nil
}

func resourceSiteRecoveryDelete(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	siteRecoveryClient := draas.NewSiteRecoveryClient(connectorWrapper)

	orgID := (m.(*connector.Wrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	siteRecoveryDeleteTask, err := siteRecoveryClient.Delete(orgID, sddcID, nil, nil)
	if err != nil {
		return HandleDeleteError("Site recovery", sddcID, err)
	}
	return resource.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return task.GetDraasTask(connectorWrapper, siteRecoveryDeleteTask.Id)
			},
			"error deactivating site recovery for SDDC ",
			nil)
		if taskErr != nil {
			return taskErr
		}
		d.SetId("")
		return nil
	})
}

func resourceSiteRecoveryUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("srm_extension_key_suffix") {
		err := resourceSiteRecoveryDelete(d, m)
		if err != nil {
			return HandleDeleteError("Site Recovery", d.Get("sddc_id").(string), err)
		}

		// This wait is required after deactivation before activation
		time.Sleep(15 * time.Minute)

		err = resourceSiteRecoveryCreate(d, m)
		if err != nil {
			return HandleCreateError("Site Recovery", err)
		}
	}
	return nil
}
