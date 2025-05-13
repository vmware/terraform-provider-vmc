// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	draasmodel "github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	task "github.com/vmware/terraform-provider-vmc/vmc/task"
)

// srmNodeCreationLockMutex a mutex that allows only a single srm node per sddc to be created.
var srmNodeCreationLockMutex = task.KeyedMutex{}

func resourceSrmNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceSrmNodeCreate,
		Read:   resourceSrmNodeRead,
		Delete: resourceSrmNodeDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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
				if err := d.Set("sddc_id", idParts[1]); err != nil {
					return nil, err
				}
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

func resourceSrmNodeCreate(d *schema.ResourceData, m interface{}) error {
	err := (m.(*connector.Wrapper)).Authenticate()
	if err != nil {
		return fmt.Errorf("authentication error from Cloud Service Provider: %s", err)
	}
	connectorWrapper := m.(*connector.Wrapper)

	siteRecoverySrmNodesClient := draas.NewSiteRecoverySrmNodesClient(connectorWrapper)

	srmExtensionKeySuffix := d.Get("srm_node_extension_key_suffix").(string)
	orgID := (m.(*connector.Wrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)

	unlockFn := srmNodeCreationLockMutex.Lock(sddcID)
	provisionSrmConfigParam := &draasmodel.ProvisionSrmConfig{
		SrmExtensionKeySuffix: &srmExtensionKeySuffix,
	}

	srmNodeCreateTask, err := siteRecoverySrmNodesClient.Post(orgID, sddcID, provisionSrmConfigParam)

	if err != nil {
		return HandleCreateError("SRM Node", err)
	}

	d.SetId(*srmNodeCreateTask.ResourceId)
	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return task.GetDraasTask(connectorWrapper, srmNodeCreateTask.Id)
			},
			"error creating SRM node",
			func(_ model.Task) {
				unlockFn()
			})
		if taskErr != nil {
			return taskErr
		}
		err = resourceSrmNodeRead(d, m)
		if err == nil {
			return nil
		}
		return retry.NonRetryableError(err)
	})
}

func resourceSrmNodeRead(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := (m.(*connector.Wrapper)).Connector
	orgID := (m.(*connector.Wrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	srmNodeID := d.Id()
	siteRecoveryClient := draas.NewSiteRecoveryClient(connectorWrapper)
	siteRecovery, err := siteRecoveryClient.Get(orgID, sddcID)
	if err != nil {
		return HandleReadError(d, "SRM Node", sddcID, err)
	}
	srmNodeMap := map[string]string{}
	if err := d.Set("sddc_id", *siteRecovery.SddcId); err != nil {
		return err
	}
	for _, SRMNode := range siteRecovery.SrmNodes {
		if *SRMNode.Id == srmNodeID {
			srmNodeMap["id"] = *SRMNode.Id
			srmNodeMap["ip_address"] = *SRMNode.IpAddress
			srmNodeMap["host_name"] = *SRMNode.Hostname
			srmNodeMap["state"] = *SRMNode.State
			srmNodeMap["type"] = *SRMNode.Type_
			// During tests VmMorefId might be nil
			if SRMNode.VmMorefId != nil {
				srmNodeMap["vm_moref_id"] = *SRMNode.VmMorefId
			}
			hostName := strings.TrimPrefix(*SRMNode.Hostname, constants.SrmPrefix)
			partStr := strings.Split(hostName, constants.SddcSuffix)
			if err := d.Set("srm_node_extension_key_suffix", partStr[0]); err != nil {
				return err
			}
			break
		}
	}
	if err := d.Set("srm_instance", srmNodeMap); err != nil {
		return err
	}
	return nil
}

func resourceSrmNodeDelete(d *schema.ResourceData, m interface{}) error {
	connectorWrapper := m.(*connector.Wrapper)
	siteRecoverySrmNodesClient := draas.NewSiteRecoverySrmNodesClient(connectorWrapper)

	orgID := (m.(*connector.Wrapper)).OrgID
	sddcID := d.Get("sddc_id").(string)
	unlockFn := srmNodeCreationLockMutex.Lock(sddcID)
	srmNodeID := d.Id()
	srmNodeDeleteTask, err := siteRecoverySrmNodesClient.Delete(orgID, sddcID, srmNodeID)
	if err != nil {
		return HandleDeleteError("SRM Node", sddcID, err)
	}
	return retry.RetryContext(context.Background(), d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper,
			func() (model.Task, error) {
				return task.GetDraasTask(connectorWrapper, srmNodeDeleteTask.Id)
			},
			"failed to delete SRM node",
			func(_ model.Task) {
				unlockFn()
			})
		if taskErr != nil {
			return taskErr
		}
		d.SetId("")
		return nil
	})
}
