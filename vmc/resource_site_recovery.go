/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"log"

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

		Schema: map[string]*schema.Schema{
			"sddc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "SDDC identifier",
			},
			"srm_extension_key_suffix": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Custom extension key suffix for SRM. If not specified, default extension key will be used. The custom extension suffix must contain 13 characters or less, be composed of letters, numbers, ., -, and _ characters. The extension suffix must begin and end with a letter or number. The suffix is appended to com.vmware.vcDr- to form the full extension key",
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
		},
	}
}

func resourceSiteRecoveryCreate(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector

	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)

	srmExtensionKeySuffix := d.Get("srm_extension_key_suffix").(string)
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	if srmExtensionKeySuffix != "" {

	}
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
