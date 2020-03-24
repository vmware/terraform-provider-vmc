/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/validation"

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
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 13),
				Description:  "Custom extension key suffix for SRM. If not specified, default extension key will be used. The custom extension suffix must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The extension suffix must begin and end with a letter or number. The suffix is appended to com.vmware.vcDr- to form the full extension key",
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
			"srm_nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeMap},
			},
			"vr_node": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceSiteRecoveryCreate(d *schema.ResourceData, m interface{}) error {

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
		return fmt.Errorf("Error while activating site recovery for sddc %s: %v", sddcID, err)
	}

	// Wait until site recover is activated
	taskID := task.ResourceId
	d.SetId(*taskID)
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := draas.NewDefaultTaskClient(connector)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
				err = (m.(*ConnectorWrapper)).authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("Error authenticating in CSP: %s", err))
				}
				return resource.RetryableError(fmt.Errorf("Instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("Error describing instance: %s", err))

		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSiteRecoveryRead(d, m))
	})
	return nil

}

func resourceSiteRecoveryRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	sddcID := d.Id()
	orgID := (m.(*ConnectorWrapper)).OrgID
	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)
	siteRecovery, err := siteRecoveryClient.Get(orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting the SDDC with ID  : %v", err)
	}
	d.SetId(siteRecovery.Id)
	d.Set("site_recovery_state", siteRecovery.SiteRecoveryState)
	d.Set("draas_h5_url", siteRecovery.DraasH5Url)
	d.Set("user_id", siteRecovery.UserId)
	d.Set("user_name", siteRecovery.UserName)

	srm_nodes := []map[string]string{}
	for _, srmNode := range siteRecovery.SrmNodes {
		m := map[string]string{}
		m["id"] = *srmNode.Id
		m["ip_address"] = *srmNode.IpAddress
		m["host_name"] = *srmNode.Hostname
		m["state"] = *srmNode.State
		m["type"] = *srmNode.Type_
		m["vm_moref_id"] = *srmNode.VmMorefId
		srm_nodes = append(srm_nodes, m)
	}

	vr_node := map[string]string{}
	vr_node["vm_moref_id"] = *siteRecovery.VrNode.VmMorefId
	vr_node["id"] = *siteRecovery.VrNode.Id
	vr_node["hostname"] = *siteRecovery.VrNode.Hostname
	vr_node["type"] = *siteRecovery.VrNode.Type_
	vr_node["state"] = *siteRecovery.VrNode.State
	vr_node["ip_address"] = *siteRecovery.VrNode.IpAddress

	d.Set("srm_nodes", srm_nodes)
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
		return fmt.Errorf("Error while deactivating site recovery for sddc %s: %v", sddcID, err)
	}
	tasksClient := draas.NewDefaultTaskClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deactivating site recovery for sddc %s : %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}
