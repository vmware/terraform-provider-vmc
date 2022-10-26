/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package sddc_group

type DeploymentGroupMember struct {
	Id string `json:"id"`
}

type NetworkConnectivityConfig struct {
	Id      string `json:"id"`
	GroupId string `json:"group_id"`
	Name    string `json:"name"`
}

type ValidationPayload struct {
	DeploymentGroupId string                  `json:"deployment_group_id,omitempty"`
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
	ConfigId string `json:"config_id"`
	GroupId  string `json:"group_id"`
	TaskId   string `json:"operation_id"`
}

type GroupMember struct {
	Id string `json:"deployment_id"`
}

type Membership struct {
	Excluded []GroupMember `json:"excluded"`
	Included []GroupMember `json:"included"`
}

type DeploymentGroup struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	OrgId       string     `json:"org_id"`
	Deleted     bool       `json:"deleted"`
	Membership  Membership `json:"membership"`
}

type NetworkOperation struct {
	Id           string `json:"id,omitempty"`
	OrgId        string `json:"org_id"`
	ResourceId   string `json:"resource_id"`
	ResourceType string `json:"resource_type"`
	TaskId       string `json:"task_id,omitempty"`
	Type         string `json:"type"`
	Config       Config `json:"config"`
}

const NetworkConnectivityConfigResourceType = "network-connectivity-config"

func NewNetworkOperation(
	orgId string,
	resourceId string,
	networkOperationType string,
	config Config) *NetworkOperation {
	return &NetworkOperation{
		OrgId:      orgId,
		ResourceId: resourceId,
		// Required by the API, without it, it responds with 400
		ResourceType: NetworkConnectivityConfigResourceType,
		Type:         networkOperationType,
		Config:       config,
	}
}

type Config struct {
	Type          string                  `json:"type"`
	OperationId   string                  `json:"operation_id,omitempty"`
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
