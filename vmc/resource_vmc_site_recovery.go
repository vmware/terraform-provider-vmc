/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas/model"
)

func resourceSiteRecovery() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteRecoveryCreate,
		Read:   resourceSiteRecoveryRead,
		Update: resourceSiteRecoveryUpdate,
		Delete: resourceSiteRecoveryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
		},
	}
}

func resourceSiteRecoveryCreate(d *schema.ResourceData, m interface{}) error {

	err := (m.(*ConnectorWrapper)).authenticate()
	if err != nil {
		return fmt.Errorf("authentication error from Cloud Service Provider: %s", err)
	}
	connector := (m.(*ConnectorWrapper)).Connector

	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)

	srmExtensionKeySuffix := d.Get("srm_extension_key_suffix").(string)
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	activateSiteRecoveryConfigParam := &model.ActivateSiteRecoveryConfig{
		SrmExtensionKeySuffix: &srmExtensionKeySuffix,
	}

	task, err := siteRecoveryClient.Post(orgID, sddcID, activateSiteRecoveryConfigParam)

	if err != nil {
		return HandleCreateError("Site recovery", err)
	}

	// Wait until site recovery is activated
	taskID := task.ResourceId
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := draas.NewDefaultTaskClient(connector)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
				err = (m.(*ConnectorWrapper)).authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider:: %s", err))
				}
				return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))

		}

		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		}
		d.SetId(*taskID)
		return resource.NonRetryableError(resourceSiteRecoveryRead(d, m))
	})
}

func resourceSiteRecoveryRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcID := d.Id()
	orgID := (m.(*ConnectorWrapper)).OrgID
	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)
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
	srm_node := map[string]string{}
	for i := 0; i < len(siteRecovery.SrmNodes); i++ {
		currentSRMNode := siteRecovery.SrmNodes[i]
		if len(strings.TrimSpace(srmExtensionKey)) == 0 {
			tempStr := strings.Trim(*currentSRMNode.Hostname, ".")
			if strings.Contains(tempStr, "-") {
				srm_node["id"] = *currentSRMNode.Id
				srm_node["ip_address"] = *currentSRMNode.IpAddress
				srm_node["host_name"] = *currentSRMNode.Hostname
				srm_node["state"] = *currentSRMNode.State
				srm_node["type"] = *currentSRMNode.Type_
				srm_node["vm_moref_id"] = *currentSRMNode.VmMorefId
				break
			}
		} else if strings.Contains(*currentSRMNode.Hostname, srmExtensionKey) {
			srm_node["id"] = *currentSRMNode.Id
			srm_node["ip_address"] = *currentSRMNode.IpAddress
			srm_node["host_name"] = *currentSRMNode.Hostname
			srm_node["state"] = *currentSRMNode.State
			srm_node["type"] = *currentSRMNode.Type_
			srm_node["vm_moref_id"] = *currentSRMNode.VmMorefId
			break
		}
	}

	vr_node := map[string]string{}
	vr_node["vm_moref_id"] = *siteRecovery.VrNode.VmMorefId
	vr_node["id"] = *siteRecovery.VrNode.Id
	vr_node["hostname"] = *siteRecovery.VrNode.Hostname
	vr_node["type"] = *siteRecovery.VrNode.Type_
	vr_node["state"] = *siteRecovery.VrNode.State
	vr_node["ip_address"] = *siteRecovery.VrNode.IpAddress

	d.Set("srm_node", srm_node)
	d.Set("vr_node", vr_node)
	return nil
}

func resourceSiteRecoveryDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)

	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	task, err := siteRecoveryClient.Delete(orgID, sddcID, nil, nil)
	if err != nil {
		return HandleDeleteError("Site recovery", sddcID, err)
	}
	tasksClient := draas.NewDefaultTaskClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error deactivating site recovery for SDDC %s : %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}

func resourceSiteRecoveryUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChange("srm_extension_key_suffix") {
		err := resource.NonRetryableError(resourceSiteRecoveryDelete(d, m))
		if err != nil {
			return HandleDeleteError("Site Recovery", d.Get("sddc_id").(string), err.Err)
		}

		// This wait is required after deactivation before activation
		time.Sleep(15 * time.Minute)

		err = resource.NonRetryableError(resourceSiteRecoveryCreate(d, m))
		if err != nil {
			return HandleCreateError("Site Recovery", err.Err)
		}
	}
	return nil
}
