/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/sddcgroup"
	"github.com/vmware/terraform-provider-vmc/vmc/task"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"strings"
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
			Create: schema.DefaultTimeout(90 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
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
		"creator": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"timestamp": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"tgw_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"tgw_region": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"vpc_aws_account": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"vpc_ram_share_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"vpc_attachment_status": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"vpc_attachments": {
			Type: schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"vpc_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"state": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"attach_id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"configured_prefixes": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
				},
			},
			Computed: true,
			Optional: true,
		},
		"dxgw_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"dxgw_owner": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"dxgw_status": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"dxgw_allowed_prefixes": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"external_tgw_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"external_tgw_owner": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"external_tgw_region": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"external_tgw_configured_prefixes": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
	}
}

func resourceSddcGroupCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	connectorWrapper := i.(*connector.Wrapper)
	sddcGroupsClient := sddcgroup.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcMemberIDs := getCurrentSddcMemberIDs(data)
	err = sddcGroupsClient.ValidateCreateSddcGroup(sddcMemberIDs)
	if err != nil {
		return diag.FromErr(err)
	}
	sddcGroupName := data.Get("name").(string)
	sddcGroupDescription := data.Get("description").(string)
	sddcGroupID, taskID, err := sddcGroupsClient.CreateSddcGroup(sddcGroupName, sddcGroupDescription, sddcMemberIDs)
	if err != nil {
		return diag.FromErr(err)
	}
	data.SetId(sddcGroupID)
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, taskID)
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
	connectorWrapper := i.(*connector.Wrapper)
	sddcGroupsClient := sddcgroup.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcGroupID := data.Id()
	sddcGroup, networkConnectivityConfig, err := sddcGroupsClient.GetSddcGroup(sddcGroupID)
	if err != nil {
		return diag.FromErr(err)
	}
	if sddcGroup == nil {
		return diag.FromErr(fmt.Errorf("sddcGroup %s is nil after trying to fetch it", sddcGroupID))
	}
	_ = data.Set("name", sddcGroup.Name)
	_ = data.Set("description", sddcGroup.Description)
	_ = data.Set("org_id", sddcGroup.OrgID)
	_ = data.Set("deleted", sddcGroup.Deleted)
	_ = data.Set("creator", sddcGroup.Creator.UserName)
	_ = data.Set("timestamp", sddcGroup.Creator.Timestamp)
	var sddcMemberIDs []string
	for _, groupMember := range sddcGroup.Membership.Included {
		sddcMemberIDs = append(sddcMemberIDs, groupMember.ID)
	}
	_ = data.Set("sddc_member_ids", sddcMemberIDs)
	if networkConnectivityConfig == nil || networkConnectivityConfig.Traits == nil {
		// below data cannot be read, so skip
		return nil
	}
	if networkConnectivityConfig.Traits.TransitGateway != nil &&
		len(networkConnectivityConfig.Traits.TransitGateway.L3Connectors) > 0 {
		_ = data.Set("tgw_id", networkConnectivityConfig.Traits.TransitGateway.L3Connectors[0].ID)
		_ = data.Set("tgw_region", networkConnectivityConfig.Traits.TransitGateway.L3Connectors[0].Region)
	}
	if networkConnectivityConfig.Traits.AwsInfo != nil &&
		len(networkConnectivityConfig.Traits.AwsInfo.Accounts) > 0 {
		_ = data.Set("vpc_aws_account", networkConnectivityConfig.Traits.AwsInfo.Accounts[0].AccountNumber)
		_ = data.Set("vpc_ram_share_id", networkConnectivityConfig.Traits.AwsInfo.Accounts[0].RAMShareID)
		_ = data.Set("vpc_attachment_status", networkConnectivityConfig.Traits.AwsInfo.Accounts[0].Status)
		var vpcAttachments []map[string]string
		for _, vpcAttachment := range networkConnectivityConfig.Traits.AwsInfo.Accounts[0].AccountAttachments {
			vpcAttachments = append(vpcAttachments, map[string]string{
				"vpc_id":              vpcAttachment.VpcID,
				"state":               vpcAttachment.State,
				"attach_id":           vpcAttachment.AttachmentID,
				"configured_prefixes": strings.Join(vpcAttachment.StaticRoutes, " "),
			})
		}
		_ = data.Set("vpc_attachments", vpcAttachments)
	}
	if networkConnectivityConfig.Traits.DxGateway != nil &&
		len(networkConnectivityConfig.Traits.DxGateway.DirectConnectGatewayAssociations) > 0 {
		_ = data.Set("dxgw_id", networkConnectivityConfig.Traits.DxGateway.DirectConnectGatewayAssociations[0].DxgwID)
		_ = data.Set("dxgw_owner", networkConnectivityConfig.Traits.DxGateway.DirectConnectGatewayAssociations[0].DxgwOwner)
		_ = data.Set("dxgw_status", networkConnectivityConfig.Traits.DxGateway.DirectConnectGatewayAssociations[0].Status)
		_ = data.Set("dxgw_allowed_prefixes", strings.Join(networkConnectivityConfig.Traits.DxGateway.
			DirectConnectGatewayAssociations[0].PeeringRegions[0].AllowedPrefixes, " "))
	}
	if networkConnectivityConfig.Traits.ExternalTgw != nil &&
		len(networkConnectivityConfig.Traits.ExternalTgw.CustomerTransitGatewayAssociations) > 0 {
		_ = data.Set("external_tgw_id", networkConnectivityConfig.Traits.ExternalTgw.CustomerTransitGatewayAssociations[0].TgwID)
		_ = data.Set("external_tgw_owner", networkConnectivityConfig.Traits.ExternalTgw.CustomerTransitGatewayAssociations[0].TgwOwner)
		_ = data.Set("external_tgw_region", networkConnectivityConfig.Traits.ExternalTgw.CustomerTransitGatewayAssociations[0].TgwRegion)
		_ = data.Set("external_tgw_configured_prefixes", strings.Join(networkConnectivityConfig.Traits.ExternalTgw.
			CustomerTransitGatewayAssociations[0].PeeringRegions[0].ConfiguredPrefixes, " "))
	}
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
	connectorWrapper := i.(*connector.Wrapper)
	sddcGroupsClient := sddcgroup.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}
	sddcMemberIds := getCurrentSddcMemberIDs(data)
	// Removal of all sddc members from the group is required prior to deletion
	diags := updateSddcGroupMembers(data, i, new([]string), sddcMemberIds)
	if diags != nil {
		return diags
	}

	deleteSddcTaskID, err := sddcGroupsClient.DeleteSddcGroup(data.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, deleteSddcTaskID)
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
	connectorWrapper := i.(*connector.Wrapper)
	sddcGroupsClient := sddcgroup.NewSddcGroupClient(connectorWrapper.VmcURL, connectorWrapper.CspURL,
		connectorWrapper.RefreshToken, connectorWrapper.OrgID)
	err := sddcGroupsClient.Authenticate()
	if err != nil {
		return diag.FromErr(err)
	}

	updateMembersTaskID, err := sddcGroupsClient.UpdateSddcGroupMembers(data.Id(), addedIds, removedIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = resource.RetryContext(context.Background(), data.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		taskErr := task.RetryTaskUntilFinished(connectorWrapper, func() (model.Task, error) {
			return task.GetV2Task(connectorWrapper, updateMembersTaskID)
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

func getCurrentSddcMemberIDs(data *schema.ResourceData) *[]string {
	sddcMemberIdsSet := data.Get("sddc_member_ids").(*schema.Set)
	var sddcMemberIDs []string
	for _, sddcMemberID := range sddcMemberIdsSet.List() {
		sddcMemberIDs = append(sddcMemberIDs, sddcMemberID.(string))
	}
	return &sddcMemberIDs
}

func getAddedIds(oldIDs *schema.Set, newIDs *schema.Set) *[]string {
	var addedIds []string
	for _, newID := range newIDs.List() {
		if !oldIDs.Contains(newID) {
			addedIds = append(addedIds, newID.(string))
		}
	}
	return &addedIds
}

func getRemovedIds(oldIDs *schema.Set, newIDs *schema.Set) *[]string {
	var removedIds []string
	for _, oldID := range oldIDs.List() {
		if !newIDs.Contains(oldID) {
			removedIds = append(removedIds, oldID.(string))
		}
	}
	return &removedIds
}
