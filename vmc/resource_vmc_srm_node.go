/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/validation"

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

	err := (m.(*ConnectorWrapper)).authenticate()
	if err != nil {
		return fmt.Errorf("authentication error from Cloud Service Provider: %s", err)
	}
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
		return HandleCreateError("SRM Node", err)
	}

	d.SetId(*task.ResourceId)
	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		tasksClient := draas.NewDefaultTaskClient(connector)
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			if err.Error() == (errors.Unauthenticated{}.Error()) {
				log.Print("Auth error", err.Error(), errors.Unauthenticated{}.Error())
				err = (m.(*ConnectorWrapper)).authenticate()
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider: %s", err))
				}
				return resource.RetryableError(fmt.Errorf("instance creation still in progress"))
			}
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))
		}
		if *task.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(resourceSRMNodeRead(d, m))
	})
}

func resourceSRMNodeRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	orgID := (m.(*ConnectorWrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	srmNodeID := d.Id()
	siteRecoveryClient := draas.NewDefaultSiteRecoveryClient(connector)
	siteRecovery, err := siteRecoveryClient.Get(orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "SRM Node", sddcID, err)
	}
	srm_node := map[string]string{}
	d.Set("sddc_id", *siteRecovery.SddcId)
	for _, SRMNode := range siteRecovery.SrmNodes {
		if *SRMNode.Id == srmNodeID {
			srm_node["id"] = *SRMNode.Id
			srm_node["ip_address"] = *SRMNode.IpAddress
			srm_node["host_name"] = *SRMNode.Hostname
			srm_node["state"] = *SRMNode.State
			srm_node["type"] = *SRMNode.Type_
			srm_node["vm_moref_id"] = *SRMNode.VmMorefId
			hostName := strings.TrimPrefix(*SRMNode.Hostname, SRMPrefix)
			partStr := strings.Split(hostName, SDDCSuffix)
			d.Set("srm_node_extension_key_suffix", partStr[0])
			break
		}
	}
	d.Set("srm_instance", srm_node)
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
		return HandleDeleteError("SRM Node", sddcID, err)
	}
	tasksClient := draas.NewDefaultTaskClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error deactivating site recovery instance for SDDC %s : %v", sddcID, err))
		}
		if *task.Status == "FAILED" {
			return resource.NonRetryableError(fmt.Errorf("expected instance to be deleted but was in state %s", *task.Status))
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("expected instance to be deleted but was in state %s", *task.Status))
		}
		d.SetId("")
		return resource.NonRetryableError(nil)
	})
}
