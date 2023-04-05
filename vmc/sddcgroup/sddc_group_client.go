/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package sddcgroup

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"io"
	"net/http"
)

const authnHeader = "csp-auth-token"

type Client interface {
	connector.Authenticator
	ValidateCreateSddcGroup(sddcIDs *[]string) error
	ValidateUpdateSddcGroupMembers(groupID string, sddcIDs *[]string) error
	GetSddcGroup(groupID string) (sddcGroup DeploymentGroup, error error)
	CreateSddcGroup(name string, description string, sddcIDs *[]string) (groupID string, taskID string, error error)
	UpdateSddcGroupMembers(groupID string, sddcIDsToAdd *[]string, sddcIDsToRemove *[]string) (taskID string, error error)
	DeleteSddcGroup(groupID string) (taskID string, error error)
}

// HTTPClient an interface, that is implemented by the http.DefaultClient,
// intended to enable stubbing for testing purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ClientImpl struct {
	connector  connector.Wrapper
	httpClient HTTPClient
}

func NewSddcGroupClient(wrapper connector.Wrapper) *ClientImpl {
	copyWrapper := connector.CopyWrapper(wrapper)
	return &ClientImpl{
		connector:  *copyWrapper,
		httpClient: http.DefaultClient,
	}
}

// newTestSddcGroupClient intended for injecting dummy accessToken and stubbed httpClient for
// testing purposes.
func newTestSddcGroupClient(vmcURL string, orgID string, accessToken string, httpClient HTTPClient) *ClientImpl {
	testConnector := connector.Wrapper{
		VmcURL: vmcURL,
		OrgID:  orgID,
	}
	// Create a dummy connector to house the access token in a security context
	testConnector.Connector = client.NewConnector("", client.WithHttpClient(&http.Client{}),
		client.WithSecurityContext(security.NewOauthSecurityContext(accessToken)))
	return &ClientImpl{
		connector:  testConnector,
		httpClient: httpClient,
	}
}

// Authenticate grab an access token and set it into the Client instance for later use
func (client *ClientImpl) Authenticate() error {
	return client.connector.Authenticate()
}

func (client *ClientImpl) ValidateCreateSddcGroup(sddcIDs *[]string) error {
	return client.validateCreateUpdateSddcGroupInternal("", sddcIDs)
}

func (client *ClientImpl) ValidateUpdateSddcGroupMembers(groupID string, sddcIDs *[]string) error {
	return client.validateCreateUpdateSddcGroupInternal(groupID, sddcIDs)
}

func (client *ClientImpl) validateCreateUpdateSddcGroupInternal(groupID string, sddcIDs *[]string) error {
	validationPayload := ValidationPayload{}

	if groupID != "" {
		validationPayload.DeploymentGroupID = groupID
	}

	for _, sddcID := range *sddcIDs {
		validationPayload.Members = append(validationPayload.Members, DeploymentGroupMember{ID: sddcID})
	}

	requestPayload, err := json.Marshal(validationPayload)
	if err != nil {
		return err
	}
	validateCreateURL := client.getBaseURL() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/validate-members", client.connector.OrgID)

	req := client.createNewRequest(http.MethodPost, validateCreateURL, bytes.NewBuffer(requestPayload))

	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return err
	}
	// The API returns 409 id there are validation errors
	if statusCode == http.StatusConflict {
		return toConflictError(rawResponse)
	}
	return nil
}

func (client *ClientImpl) GetSddcGroup(groupID string) (*DeploymentGroup,
	*NetworkConnectivityConfig, error) {
	getSddcGroupURL := client.getBaseURL() + fmt.Sprintf("/inventory/%s/core/deployment-groups/%s",
		client.connector.OrgID, groupID)
	req := client.createNewRequest(http.MethodGet, getSddcGroupURL, nil)
	rawResponse, statusCode, err := client.executeRequest(req)
	var group *DeploymentGroup
	var config *NetworkConnectivityConfig
	if err != nil {
		return group, config, err
	}
	if statusCode == http.StatusNotFound {
		return group, config, fmt.Errorf("SDDC Group with ID: %s not found", groupID)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&group)
		if err != nil {
			return group, config, err
		}
	}
	// No need to query for resource ID and NetworkConnectivityConfig for deleted sddc groups
	if group != nil && group.Deleted {
		return group, config, err
	}
	resourceID, err := client.getResourceIDFromGroupID(groupID)
	if err != nil {
		return group, config, err
	}
	getTraitsURL := client.getBaseURL() + fmt.Sprintf("/network/%s/core/network-connectivity-configs/%s/"+
		"?trait=AwsVpcAttachmentsTrait,AwsDirectConnectGatewayAssociationsTrait,"+
		"AwsNetworkConnectivityTrait,AwsCustomerTransitGatewayAssociationsTrait",
		client.connector.OrgID, resourceID)
	req = client.createNewRequest(http.MethodGet, getTraitsURL, nil)
	rawResponse, statusCode, err = client.executeRequest(req)
	if err != nil {
		return group, config, err
	}
	if statusCode == http.StatusNotFound {
		return group, config, fmt.Errorf("No NetworkConnectivityConfig with ID: %s not found ", resourceID)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&config)
		return group, config, err
	}

	return group, config, fmt.Errorf("GetSddcGroup response code: %d", statusCode)
}

func (client *ClientImpl) CreateSddcGroup(
	name string,
	description string,
	sddcIDs *[]string) (groupID string, taskID string, error error) {
	createGroupNetworkConnectivityRequest := CreateGroupNetworkConnectivityRequest{}

	createGroupNetworkConnectivityRequest.Name = name
	createGroupNetworkConnectivityRequest.Description = description
	for _, sddcID := range *sddcIDs {
		createGroupNetworkConnectivityRequest.Members =
			append(createGroupNetworkConnectivityRequest.Members, DeploymentGroupMember{ID: sddcID})
	}

	requestPayload, err := json.Marshal(createGroupNetworkConnectivityRequest)
	if err != nil {
		return "", "", err
	}
	createSddcGroupURL := client.getBaseURL() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/create-group-network-connectivity",
		client.connector.OrgID)

	req := client.createNewRequest(http.MethodPost, createSddcGroupURL, bytes.NewBuffer(requestPayload))

	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return "", "", err
	}
	// The API returns 409 id there are validation errors
	if statusCode == http.StatusConflict {
		return "", "", toConflictError(rawResponse)
	}
	if statusCode == http.StatusOK {
		var createGroupNetworkConnectivityResponse CreateGroupNetworkConnectivityResponse
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&createGroupNetworkConnectivityResponse)
		if err != nil {
			return "", "", err
		}
		return createGroupNetworkConnectivityResponse.GroupID, createGroupNetworkConnectivityResponse.TaskID, nil
	}
	return "", "", fmt.Errorf("CreateSddcGroup response code: %d", statusCode)
}

func (client *ClientImpl) UpdateSddcGroupMembers(
	groupID string, sddcIDsToAdd *[]string, sddcIDsToRemove *[]string) (taskID string, error error) {
	var addMembers []DeploymentGroupMember
	var removeMembers []DeploymentGroupMember

	for _, sddcIDToAdd := range *sddcIDsToAdd {
		addMembers = append(addMembers, DeploymentGroupMember{
			ID: sddcIDToAdd,
		})
	}
	for _, sddcIDToRemove := range *sddcIDsToRemove {
		removeMembers = append(removeMembers, DeploymentGroupMember{
			ID: sddcIDToRemove,
		})
	}
	resourceID, err := client.getResourceIDFromGroupID(groupID)
	if err != nil {
		return "", err
	}
	config := NewAwsUpdateDeploymentGroupMembersConfig(addMembers, removeMembers)
	networkOperation := NewNetworkOperation(client.connector.OrgID, resourceID, UpdateMembersNetworkOperationType, *config)
	networkOperationResponse, err := client.executeNetworkOperation(networkOperation)
	if err != nil {
		return "", err
	}
	return networkOperationResponse.Config.OperationID, nil
}

func (client *ClientImpl) DeleteSddcGroup(groupID string) (taskID string, error error) {
	resourceID, err := client.getResourceIDFromGroupID(groupID)
	if err != nil {
		return "", err
	}
	config := NewAwsDeleteDeploymentGroupConfig()
	networkOperation := NewNetworkOperation(client.connector.OrgID, resourceID, DeleteSddcGroupNetworkOperationType, *config)
	networkOperationResponse, err := client.executeNetworkOperation(networkOperation)
	if err != nil {
		return "", err
	}
	return networkOperationResponse.ID, nil
}

func (client *ClientImpl) getResourceIDFromGroupID(groupID string) (resourceID string, error error) {
	getResourceIDURL := client.getBaseURL() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/?group_id=%s", client.connector.OrgID, groupID)

	req := client.createNewRequest(http.MethodGet, getResourceIDURL, nil)

	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return "", err
	}
	var result []NetworkConnectivityConfig
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result[0].ID, err
	}
	return "", fmt.Errorf("getResourceIDFromGroupID failed with status %d body: %s",
		statusCode, string(*rawResponse))
}

func (client *ClientImpl) executeNetworkOperation(networkOperation *NetworkOperation) (networkOperationResponse *NetworkOperation, error error) {
	networkOperationResponse = nil
	requestPayload, err := json.Marshal(networkOperation)
	if err != nil {
		return networkOperationResponse, err
	}
	networkOperationsURL := client.getNetworkOperationsURL()

	req := client.createNewRequest(http.MethodPost, networkOperationsURL, bytes.NewBuffer(requestPayload))

	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return networkOperationResponse, err
	}

	// The API returns 409 id there are validation errors
	if statusCode == http.StatusConflict {
		return networkOperationResponse, toConflictError(rawResponse)
	}
	if statusCode == http.StatusOK || statusCode == http.StatusCreated {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&networkOperationResponse)
		if err != nil {
			return networkOperationResponse, err
		}
		return networkOperationResponse, nil
	}
	return networkOperationResponse, fmt.Errorf("%s response code: %d \n body: %s",
		networkOperation.Type, statusCode, string(*rawResponse))
}

func (client *ClientImpl) getNetworkOperationsURL() string {
	return client.getBaseURL() + fmt.Sprintf(
		"/network/%s/aws/operations", client.connector.OrgID)
}

// executeRequest Returns the body of the response as byte array pointer, the status code
// or any error that may have occurred during the Http communication.
func (client *ClientImpl) executeRequest(
	request *http.Request) (responseBody *[]byte, statusCode int, error error) {
	response, err := client.httpClient.Do(request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing body of http response")
		}
	}(response.Body)

	if err != nil {
		return nil, response.StatusCode, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, response.StatusCode, fmt.Errorf("Unauthenticated request ")
	}
	if response.StatusCode == http.StatusForbidden {
		return nil, response.StatusCode, fmt.Errorf("Unauthorized request ")
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return &result, response.StatusCode, nil
}

func (client *ClientImpl) getBaseURL() string {
	return client.connector.VmcURL + "/api"
}

func (client *ClientImpl) createNewRequest(method string, URL string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, URL, body)
	req.Header.Add(authnHeader, client.connector.Connector.SecurityContext().Property(security.ACCESS_TOKEN).(string))
	if method == http.MethodPost {
		req.Header.Add("content-type", "application/json")
	}
	return req
}

func validationErrorResponseToString(errorResponse *ValidationErrorResponse) string {
	errorMessage := ""
	for _, detail := range errorResponse.Details {
		errorMessage += detail.ValidationErrorMessage
		if len(detail.Members) > 0 {
			errorMessage += " For members: "
			for _, member := range detail.Members {
				errorMessage += member + " "
			}
		}
	}
	return errorMessage
}

func toConflictError(rawResponse *[]byte) error {
	var validationErrorResponse ValidationErrorResponse
	err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&validationErrorResponse)
	if err != nil {
		return err
	}
	return fmt.Errorf(validationErrorResponseToString(&validationErrorResponse))
}
