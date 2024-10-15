// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package sddcgroup

type DeploymentGroupMember struct {
	ID string `json:"id"`
}

type L3Connector struct {
	ID     string `json:"id"`
	Region string `json:"region"`
}

type AwsNetworkConnectivityTrait struct {
	L3Connectors []L3Connector `json:"l3connectors,omitempty"`
}

type AccountAttachment struct {
	VpcID        string   `json:"vpc_id"`
	State        string   `json:"state"`
	AttachmentID string   `json:"attach_id"`
	StaticRoutes []string `json:"configured_prefixes,omitempty"`
}

type AwsAccount struct {
	AccountNumber      string              `json:"account_number"`
	RAMShareID         string              `json:"resource_share_name"`
	Status             string              `json:"state"`
	AccountAttachments []AccountAttachment `json:"attachments,omitempty"`
}

type AwsVpcAttachmentsTrait struct {
	Accounts []AwsAccount `json:"accounts"`
}

type PeeringRegions struct {
	AllowedPrefixes    []string `json:"allowed_prefixes,omitempty"`
	ConfiguredPrefixes []string `json:"configured_prefixes,omitempty"`
}

type DirectConnectGatewayAssociation struct {
	DxgwID         string           `json:"direct_connect_gateway_id"`
	DxgwOwner      string           `json:"direct_connect_gateway_owner"`
	Status         string           `json:"state"`
	PeeringRegions []PeeringRegions `json:"peering_regions"`
}

type AwsDirectConnectGatewayAssociationsTrait struct {
	DirectConnectGatewayAssociations []DirectConnectGatewayAssociation `json:"direct_connect_gateway_associations"`
}

type TgwRegion struct {
	Region string `json:"code"`
}

type CustomerTransitGatewayAssociation struct {
	TgwID          string           `json:"customer_transit_gateway_id"`
	TgwOwner       string           `json:"customer_transit_gateway_owner"`
	TgwRegion      TgwRegion        `json:"customer_transit_gateway_region"`
	PeeringRegions []PeeringRegions `json:"peering_regions"`
}

type AwsCustomerTransitGatewayAssociationsTrait struct {
	CustomerTransitGatewayAssociations []CustomerTransitGatewayAssociation `json:"customer_transit_gateway_associations,omitempty"`
}
type Traits struct {
	TransitGateway *AwsNetworkConnectivityTrait                `json:"AwsNetworkConnectivityTrait,omitempty"`
	AwsInfo        *AwsVpcAttachmentsTrait                     `json:"AwsVpcAttachmentsTrait,omitempty"`
	DxGateway      *AwsDirectConnectGatewayAssociationsTrait   `json:"AwsDirectConnectGatewayAssociationsTrait,omitempty"`
	ExternalTgw    *AwsCustomerTransitGatewayAssociationsTrait `json:"AwsCustomerTransitGatewayAssociationsTrait,omitempty"`
}

type NetworkConnectivityConfigState struct {
	Name string `json:"name"`
}

type NetworkConnectivityConfig struct {
	ID                             string                         `json:"id"`
	GroupID                        string                         `json:"group_id"`
	Name                           string                         `json:"name"`
	NetworkConnectivityConfigState NetworkConnectivityConfigState `json:"state"`
	Traits                         *Traits                        `json:"traits,omitempty"`
}

type ValidationPayload struct {
	DeploymentGroupID string                  `json:"deployment_group_id,omitempty"`
	Members           []DeploymentGroupMember `json:"members"`
}

type Details struct {
	ValidationErrorMessage string   `json:"validation_error_message"`
	Members                []string `json:"members"`
}

type ValidationErrorResponse struct {
	Status  int       `json:"status"`
	Message string    `json:"message"`
	Details []Details `json:"details"`
}

type CreateGroupNetworkConnectivityRequest struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Members     []DeploymentGroupMember `json:"members"`
}

type CreateGroupNetworkConnectivityResponse struct {
	ConfigID string `json:"config_id"`
	GroupID  string `json:"group_id"`
	TaskID   string `json:"operation_id"`
}

type GroupMember struct {
	ID string `json:"deployment_id"`
}

type Creator struct {
	UserName  string `json:"user_name"`
	Timestamp string `json:"timestamp"`
}

type Membership struct {
	Excluded []GroupMember `json:"excluded"`
	Included []GroupMember `json:"included"`
}

type DeploymentGroup struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrgID       string     `json:"org_id"`
	Deleted     bool       `json:"deleted"`
	Membership  Membership `json:"membership"`
	Creator     Creator    `json:"creator"`
}

type NetworkOperation struct {
	ID           string `json:"id,omitempty"`
	OrgID        string `json:"org_id"`
	ResourceID   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	TaskID       string `json:"task_id,omitempty"`
	Type         string `json:"type"`
	Config       Config `json:"config"`
}

const NetworkConnectivityConfigResourceType = "network-connectivity-config"

func NewNetworkOperation(
	orgID string,
	resourceID string,
	networkOperationType string,
	config Config) *NetworkOperation {
	return &NetworkOperation{
		OrgID:      orgID,
		ResourceID: resourceID,
		// Required by the API, without it, it responds with 400
		ResourceType: NetworkConnectivityConfigResourceType,
		Type:         networkOperationType,
		Config:       config,
	}
}

type Config struct {
	Type          string                  `json:"type"`
	OperationID   string                  `json:"operation_id,omitempty"`
	AddMembers    []DeploymentGroupMember `json:"add_members,omitempty"`
	RemoveMembers []DeploymentGroupMember `json:"remove_members,omitempty"`
}

const UpdateMembersNetworkOperationType = "UPDATE_MEMBERS"

func NewAwsUpdateDeploymentGroupMembersConfig(
	addMembers []DeploymentGroupMember,
	removeMembers []DeploymentGroupMember) *Config {
	return &Config{
		Type:          "AwsUpdateDeploymentGroupMembersConfig",
		AddMembers:    addMembers,
		RemoveMembers: removeMembers,
	}
}

const DeleteSddcGroupNetworkOperationType = "DELETE_DEPLOYMENT_GROUP"

func NewAwsDeleteDeploymentGroupConfig() *Config {
	return &Config{
		Type: "AwsDeleteDeploymentGroupConfig",
	}
}
