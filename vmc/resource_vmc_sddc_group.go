/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/sddc_group"
	"github.com/vmware/terraform-provider-vmc/vmc/task"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"time"
)

func resourceSddcGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSddcGroupCreate,
		ReadContext:   resourceSddcGroupRead,
		UpdateContext: resourceSddcGroupUpdate,
		DeleteContext: resourceSddcGroupDelete,
		Schema:        sddcGroupSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func sddcGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "The name of the SDDC group",
		},
		"description": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.NoZeroValues,
			Description:  "Short description of the SDDC Group",
		},
		"sddc_member_ids": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Required:    true,
			Description: "A set of the IDs of SDDC members of the SDDC Group",
		},
		"org_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deleted": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func resourceSddcGroupCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	connectorWrapper := i.(*connector.ConnectorWrapper)
	sddcGroupsClient := sddc_group.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcMemberIds := getCurrentSddcMemberIds(data)
	err = sddcGroupsClient.ValidateCreateSddcGroup(sddcMemberIds)
	if err != nil {
		return diag.FromErr(err)
	}
	sddcGroupName := data.Get("name").(string)
	sddcGroupDescription := data.Get("description").(string)
	sddcGroupId, taskId, err := sddcGroupsClient.CreateSddcGroup(sddcGroupName, sddcGroupDescription, sddcMemberIds)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(sddcGroupId)
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, taskId)
		}, "error creating SDDC group", nil)
		if taskErr != nil {
			return taskErr
		}
		diags := resourceSddcGroupRead(ctx, data, i)
		if !diags.HasError() {
			return nil
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceSddcGroupRead(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	connectorWrapper := i.(*connector.ConnectorWrapper)
	sddcGroupsClient := sddc_group.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcGroupId := data.Id()
	sddcGroup, err := sddcGroupsClient.GetSddcGroup(sddcGroupId)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = data.Set("name", sddcGroup.Name)
	_ = data.Set("description", sddcGroup.Description)
	_ = data.Set("org_id", sddcGroup.OrgId)
	_ = data.Set("deleted", sddcGroup.Deleted)
	var sddcMemberIds []string
	for _, groupMember := range sddcGroup.Membership.Included {
		sddcMemberIds = append(sddcMemberIds, groupMember.Id)
	}
	_ = data.Set("sddc_member_ids", sddcMemberIds)
	return nil
}

func resourceSddcGroupUpdate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	if data.HasChange("sddc_member_ids") {
		oldIdsRaw, newIdsRaw := data.GetChange("sddc_member_ids")
		oldIds := oldIdsRaw.(*schema.Set)
		newIds := newIdsRaw.(*schema.Set)
		addedIds := getAddedIds(oldIds, newIds)
		removedIds := getRemovedIds(oldIds, newIds)

		diags := updateSddcGroupMembers(data, i, addedIds, removedIds)
		if diags != nil {
			return diags
		}
	}
	return resourceSddcGroupRead(ctx, data, i)
}

func resourceSddcGroupDelete(_ context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	connectorWrapper := i.(*connector.ConnectorWrapper)
	sddcGroupsClient := sddc_group.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcMemberIds := getCurrentSddcMemberIds(data)
	// Removal of all sddc members from the group is required prior to deletion
	diags := updateSddcGroupMembers(data, i, new([]string), sddcMemberIds)
	if diags != nil {
		return diags
	}

	deleteSddcTaskId, err := sddcGroupsClient.DeleteSddcGroup(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, deleteSddcTaskId)
		}, "error deleting SDDC group", nil)
		if taskErr != nil {
			return taskErr
		}
		data.SetId("")
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func updateSddcGroupMembers(data *schema.ResourceData,
	i interface{}, addedIds *[]string, removedIds *[]string) diag.Diagnostics {
	connectorWrapper := i.(*connector.ConnectorWrapper)
	sddcGroupsClient := sddc_group.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}

	updateMembersTaskId, err := sddcGroupsClient.UpdateSddcGroupMembers(data.Id(), addedIds, removedIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, updateMembersTaskId)
		}, "error updating SDDC group members", nil)
		if taskErr != nil {
			return taskErr
		}
		return nil
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func getCurrentSddcMemberIds(data *schema.ResourceData) *[]string {
	sddcMemberIdsSet := data.Get("sddc_member_ids").(*schema.Set)
	var sddcMemberIds []string
	for _, sddcMemberId := range sddcMemberIdsSet.List() {
		sddcMemberIds = append(sddcMemberIds, sddcMemberId.(string))
	}
	return &sddcMemberIds
}

func getAddedIds(oldIds *schema.Set, newIds *schema.Set) *[]string {
	var addedIds []string
	for _, newId := range newIds.List() {
		if !oldIds.Contains(newId) {
			addedIds = append(addedIds, newId.(string))
		}
	}
	return &addedIds
}

func getRemovedIds(oldIds *schema.Set, newIds *schema.Set) *[]string {
	var removedIds []string
	for _, oldId := range oldIds.List() {
		if !newIds.Contains(oldId) {
			removedIds = append(removedIds, oldId.(string))
		}
	}
	return &removedIds
}
