/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package sddc_group

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"io"
	"net/http"
)

const authnHeader = "csp-auth-token"

type SddcGroupClient interface {
	connector.Authenticator
	ValidateCreateSddcGroup(sddcIds *[]string) error
	ValidateUpdateSddcGroupMembers(groupId string, sddcIds *[]string) error
	GetSddcGroup(groupId string) (sddcGroup DeploymentGroup, error error)
	CreateSddcGroup(name string, description string, sddcIds *[]string) (groupId string, taskId string, error error)
	UpdateSddcGroupMembers(groupId string, sddcIdsToAdd *[]string, sddcIdsToRemove *[]string) (taskId string, error error)
	DeleteSddcGroup(groupId string) (taskId string, error error)
}

// HttpClient an interface, that is implemented by the http.DefaultClient,
// intended to enable stubbing for testing purposes
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type SddcGroupClientImpl struct {
	vmcUrl       string
	cspUrl       string
	refreshToken string
	orgId        string
	accessToken  string
	httpClient   HttpClient
}

func NewSddcGroupClient(vmcUrl string, cspUrl string, refreshToken string, orgId string) *SddcGroupClientImpl {
	return &SddcGroupClientImpl{
		vmcUrl:       vmcUrl,
		cspUrl:       cspUrl,
		refreshToken: refreshToken,
		orgId:        orgId,
		httpClient:   http.DefaultClient,
	}
}

// newTestSddcGroupClient intended for injecting dummy accessToken and stubbed httpClient for
// testing purposes.
func newTestSddcGroupClient(vmcUrl string, orgId string, accessToken string, httpClient HttpClient) *SddcGroupClientImpl {
	return &SddcGroupClientImpl{
		vmcUrl:      vmcUrl,
		orgId:       orgId,
		accessToken: accessToken,
		httpClient:  httpClient,
	}
}

// Authenticate grab an access token and set it into the SddcGroupClient instance for later use
func (client *SddcGroupClientImpl) Authenticate() error {
	authUrl := client.cspUrl + constants.CspRefreshUrlSuffix
	securityContext, err := connector.SecurityContextByRefreshToken(client.refreshToken, authUrl)
	if err != nil {
		return err
	}
	client.accessToken = securityContext.Property(security.ACCESS_TOKEN).(string)
	return nil
}

func (client *SddcGroupClientImpl) ValidateCreateSddcGroup(sddcIds *[]string) error {
	return client.validateCreateUpdateSddcGroupInternal("", sddcIds)
}

func (client *SddcGroupClientImpl) ValidateUpdateSddcGroupMembers(groupId string, sddcIds *[]string) error {
	return client.validateCreateUpdateSddcGroupInternal(groupId, sddcIds)
}

func (client *SddcGroupClientImpl) validateCreateUpdateSddcGroupInternal(groupId string, sddcIds *[]string) error {
	validationPayload := ValidationPayload{}

	if groupId != "" {
		validationPayload.DeploymentGroupId = groupId
	}

	for _, sddcId := range *sddcIds {
		validationPayload.Members = append(validationPayload.Members, DeploymentGroupMember{Id: sddcId})
	}

	requestPayload, err := json.Marshal(validationPayload)
	if err != nil {
		return err
	}
	validateCreateUrl := client.getBaseUrl() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/validate-members", client.orgId)

	req, _ := http.NewRequest(http.MethodPost, validateCreateUrl, bytes.NewBuffer(requestPayload))
	req.Header.Add(authnHeader, client.accessToken)
	req.Header.Add("content-type", "application/json")

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

func (client *SddcGroupClientImpl) GetSddcGroup(groupId string) (sddcGroup DeploymentGroup, error error) {
	getSddcGroupUrl := client.getBaseUrl() + fmt.Sprintf("/inventory/%s/core/deployment-groups/%s", client.orgId, groupId)
	req, _ := http.NewRequest(http.MethodGet, getSddcGroupUrl, nil)
	req.Header.Add(authnHeader, client.accessToken)
	var result DeploymentGroup
	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return result, err
	}
	if statusCode == http.StatusNotFound {
		return result, fmt.Errorf("SDDC Group with ID: %s not found", groupId)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result, err
	}
	return result, fmt.Errorf("GetSddcGroup response code: %d", statusCode)
}

func (client *SddcGroupClientImpl) CreateSddcGroup(
	name string,
	description string,
	sddcIds *[]string) (groupId string, taskId string, error error) {
	createGroupNetworkConnectivityRequest := CreateGroupNetworkConnectivityRequest{}

	createGroupNetworkConnectivityRequest.Name = name
	createGroupNetworkConnectivityRequest.Description = description
	for _, sddcId := range *sddcIds {
		createGroupNetworkConnectivityRequest.Members =
			append(createGroupNetworkConnectivityRequest.Members, DeploymentGroupMember{Id: sddcId})
	}

	requestPayload, err := json.Marshal(createGroupNetworkConnectivityRequest)
	if err != nil {
		return "", "", err
	}
	createSddcGroupUrl := client.getBaseUrl() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/create-group-network-connectivity", client.orgId)

	req, _ := http.NewRequest(http.MethodPost, createSddcGroupUrl, bytes.NewBuffer(requestPayload))
	req.Header.Add(authnHeader, client.accessToken)
	req.Header.Add("content-type", "application/json")

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
		return createGroupNetworkConnectivityResponse.GroupId, createGroupNetworkConnectivityResponse.TaskId, nil
	}
	return "", "", fmt.Errorf("CreateSddcGroup response code: %d", statusCode)
}

func (client *SddcGroupClientImpl) UpdateSddcGroupMembers(
	groupId string, sddcIdsToAdd *[]string, sddcIdsToRemove *[]string) (taskId string, error error) {
	var addMembers []DeploymentGroupMember
	var removeMembers []DeploymentGroupMember

	for _, sddcIdToAdd := range *sddcIdsToAdd {
		addMembers = append(addMembers, DeploymentGroupMember{
			Id: sddcIdToAdd,
		})
	}
	for _, sddcIdToRemove := range *sddcIdsToRemove {
		removeMembers = append(removeMembers, DeploymentGroupMember{
			Id: sddcIdToRemove,
		})
	}
	resourceId, err := client.getResourceIdFromGroupId(groupId)
	if err != nil {
		return "", err
	}
	config := NewAwsUpdateDeploymentGroupMembersConfig(addMembers, removeMembers)
	networkOperation := NewNetworkOperation(client.orgId, resourceId, UpdateMembersNetworkOperationType, *config)
	networkOperationResponse, err := client.executeNetworkOperation(networkOperation)
	if err != nil {
		return "", err
	}
	return networkOperationResponse.Config.OperationId, nil
}

func (client *SddcGroupClientImpl) DeleteSddcGroup(groupId string) (taskId string, error error) {
	resourceId, err := client.getResourceIdFromGroupId(groupId)
	if err != nil {
		return "", err
	}
	config := NewAwsDeleteDeploymentGroupConfig()
	networkOperation := NewNetworkOperation(client.orgId, resourceId, DeleteSddcGroupNetworkOperationType, *config)
	networkOperationResponse, err := client.executeNetworkOperation(networkOperation)
	if err != nil {
		return "", err
	}
	return networkOperationResponse.Id, nil
}

func (client *SddcGroupClientImpl) getResourceIdFromGroupId(groupId string) (resourceId string, error error) {
	getResourceIdUrl := client.getBaseUrl() + fmt.Sprintf(
		"/network/%s/core/network-connectivity-configs/?group_id=%s", client.orgId, groupId)

	req, _ := http.NewRequest(http.MethodGet, getResourceIdUrl, nil)
	req.Header.Add(authnHeader, client.accessToken)

	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return "", err
	}
	var result []NetworkConnectivityConfig
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result[0].Id, err
	}
	return "", fmt.Errorf("getResourceIdFromGroupId failed with status %d body: %s",
		statusCode, string(*rawResponse))
}

func (client *SddcGroupClientImpl) executeNetworkOperation(networkOperation *NetworkOperation) (networkOperationResponse *NetworkOperation, error error) {
	networkOperationResponse = nil
	requestPayload, err := json.Marshal(networkOperation)
	if err != nil {
		return networkOperationResponse, err
	}
	networkOperationsUrl := client.getNetworkOperationsUrl()

	req, _ := http.NewRequest(http.MethodPost, networkOperationsUrl, bytes.NewBuffer(requestPayload))
	req.Header.Add(authnHeader, client.accessToken)
	req.Header.Add("content-type", "application/json")

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

func (client *SddcGroupClientImpl) getNetworkOperationsUrl() string {
	return client.getBaseUrl() + fmt.Sprintf(
		"/network/%s/aws/operations", client.orgId)
}

// executeRequest Returns the body of the response as byte array pointer, the status code
// or any error that may have occurred during the Http communication.
func (client *SddcGroupClientImpl) executeRequest(
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
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}

	return &result, response.StatusCode, nil
}

func (client *SddcGroupClientImpl) getBaseUrl() string {
	return client.vmcUrl + "/api"
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
