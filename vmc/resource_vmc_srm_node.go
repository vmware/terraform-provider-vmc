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

func resourceSRMNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceSRMNodeCreate,
		Read:   resourceSRMNodeRead,
		Delete: resourceSRMNodeDelete,
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
				Required:    true,
				ForceNew:    true,
				Description: "SDDC identifier",
			},
			"srm_node_extension_key_suffix": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 13),
				Description:  "The custom extension suffix for SRM must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The suffix is appended to com.vmware.vcDr- to form the full extension key. ",
			},
			"srm_instance": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceSRMNodeCreate(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector

	siteRecoverySrmNodesClient := draas.NewDefaultSiteRecoverySrmNodesClient(connector)

	srmExtensionKeySuffix := d.Get("srm_node_extension_key_suffix").(string)
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	provisionSrmConfigParam := &model.ProvisionSrmConfig{
		SrmExtensionKeySuffix: &srmExtensionKeySuffix,
	}

	task, err := siteRecoverySrmNodesClient.Post(orgID, sddcID, provisionSrmConfigParam)

	if err != nil {
		return fmt.Errorf("Error while activating site recovery instance for sddc %s: %v", sddcID, err)
	}

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
		return resource.NonRetryableError(resourceSRMNodeRead(d, m))
	})
	return nil

}

func resourceSRMNodeRead(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)
	siteRecovery, err := siteRecoveryClient.Get(orgID, sddcID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("SRM information for SDDC with ID %s not found", sddcID)
			d.SetId("")
			return fmt.Errorf("SRM information for SDDC with ID %s not found", sddcID)
		}
		return fmt.Errorf("Error while getting SRM instance information for SDDC with ID %s : %v", sddcID, err)
	}
	srmExtensionKey := d.Get("srm_node_extension_key_suffix").(string)
	srm_node := map[string]string{}

	for i := 0; i < len(siteRecovery.SrmNodes); i++ {
		currentSRMNode := siteRecovery.SrmNodes[i]
		if strings.Contains(*currentSRMNode.Hostname, srmExtensionKey) {
			srm_node["id"] = *currentSRMNode.Id
			srm_node["ip_address"] = *currentSRMNode.IpAddress
			srm_node["host_name"] = *currentSRMNode.Hostname
			srm_node["state"] = *currentSRMNode.State
			srm_node["type"] = *currentSRMNode.Type_
			srm_node["vm_moref_id"] = *currentSRMNode.VmMorefId
			d.Set("srm_instance", srm_node)
			break
		}
	}
	return nil
}

func resourceSRMNodeDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	siteRecoverySrmNodesClient := draas.NewDefaultSiteRecoverySrmNodesClient(connector)

	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	srmNodeID := d.Id()
	task, err := siteRecoverySrmNodesClient.Delete(orgID, sddcID, srmNodeID)
	if err != nil {
		return fmt.Errorf("Error while deactivating site recovery instance for sddc %s: %v", sddcID, err)
	}
	tasksClient := draas.NewDefaultTaskClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deactivating site recovery instance for sddc %s : %v", sddcID, err))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}
