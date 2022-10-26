/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package sddc_group

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

const testAccessToken = "testAccessToken"
const testOrgId = "testOrgId"
const testVmcUrl = "https://test.vmc.vmware.com"

type HttpClientStub struct {
	expectedJson                    string
	additionalResourceIdRequestJson string
	additionalRequestPassed         bool
	expectedMethod                  string
	expectedUrl                     string
	responseJson                    string
	responseCode                    int
	responseError                   error
	t                               *testing.T
}

func (stub *HttpClientStub) Do(req *http.Request) (*http.Response, error) {
	if len(stub.additionalResourceIdRequestJson) > 0 && !stub.additionalRequestPassed {
		stub.additionalRequestPassed = true
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(stub.additionalResourceIdRequestJson)),
		}, nil
	}
	if req.Body == nil {
		assert.Equal(stub.t, stub.expectedJson, "")
	} else {
		assert.Equal(stub.t, stub.expectedJson, readAsString(req.Body))
	}
	assert.Equal(stub.t, stub.expectedUrl, req.URL.String())
	assert.Equal(stub.t, stub.expectedMethod, req.Method)
	assert.Equal(stub.t, req.Header.Get(authnHeader), testAccessToken)
	response := http.Response{
		StatusCode: stub.responseCode,
		Body:       io.NopCloser(strings.NewReader(stub.responseJson)),
	}
	return &response, stub.responseError
}

func readAsString(reader io.ReadCloser) string {
	bodyBytes, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyBytes)
}

func TestValidateCreateSddcGroup(t *testing.T) {
	testOrgId := "testOrgId"
	type inputStruct struct {
		httpClientStub HttpClient
		sddcIds        *[]string
	}
	type test struct {
		input inputStruct
		want  error
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJson:   "",
					t:              t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("VMC service down"),
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJson:   "",
					t:              t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			nil,
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJson: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"," +
						"\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"male\"\n]\n}\n]\n}",
					t: t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: male "),
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		assert.Equal(t, testCase.want, sddcGroupClient.ValidateCreateSddcGroup(testCase.input.sddcIds))
	}
}

func TestValidateCreateSddcGroupMembers(t *testing.T) {
	t.Setenv(constants.VmcUrl, testVmcUrl)
	type inputStruct struct {
		httpClientStub HttpClient
		groupId        string
		sddcIds        *[]string
	}
	type test struct {
		input inputStruct
		want  error
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJson:   "",
					t:              t,
				},
				groupId: "testGroupId",
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("VMC service down"),
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJson:   "",
					t:              t,
				},
				groupId: "testGroupId",
				sddcIds: &[]string{"lele", "male"},
			},
			nil,
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJson: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"," +
						"\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"male\"\n]\n}\n]\n}",
					t: t,
				},
				groupId: "testGroupId",
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: male "),
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		assert.Equal(t, testCase.want, sddcGroupClient.ValidateUpdateSddcGroupMembers(testCase.input.groupId, testCase.input.sddcIds))
	}
}

func TestGetSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcUrl, testVmcUrl)
	type inputStruct struct {
		httpClientStub HttpClient
		groupId        string
	}
	type outputStruct struct {
		sddcGroup DeploymentGroup
		error     error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodGet,
					expectedJson:   "",
					expectedUrl:    "https://test.vmc.vmware.com/api/inventory/testOrgId/core/deployment-groups/testGroupId",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJson:   "",
					t:              t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				sddcGroup: DeploymentGroup{},
				error:     fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodGet,
					expectedJson:   "",
					expectedUrl:    "https://test.vmc.vmware.com/api/inventory/testOrgId/core/deployment-groups/testGroupId",
					responseCode:   http.StatusUnauthorized,
					responseError:  nil,
					responseJson:   "",
					t:              t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				sddcGroup: DeploymentGroup{},
				error:     fmt.Errorf("GetSddcGroup response code: 401"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodGet,
					expectedJson:   "",
					expectedUrl:    "https://test.vmc.vmware.com/api/inventory/testOrgId/core/deployment-groups/testGroupId",
					responseCode:   http.StatusNotFound,
					responseError:  nil,
					responseJson:   "",
					t:              t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				sddcGroup: DeploymentGroup{},
				error:     fmt.Errorf("SDDC Group with ID: testGroupId not found"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodGet,
					expectedJson:   "",
					expectedUrl:    "https://test.vmc.vmware.com/api/inventory/testOrgId/core/deployment-groups/testGroupId",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJson: "{\n \"id\": \"groupId\",\n  \"name\": \"test_sddc_group\",\n  \"org_id\": \"testOrgId\",\n  \"description\": \"sdfa\",\n " +
						"\"membership\": {\n \"included\": [\n {\n\"deployment_id\": \"sddcId1\"\n },\n {\n \"deployment_id\": \"sddcId2\"\n }\n ]\n },\n \"deleted\": false\n}",
					t: t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				sddcGroup: DeploymentGroup{
					Id:          "groupId",
					Name:        "test_sddc_group",
					OrgId:       testOrgId,
					Description: "sdfa",
					Deleted:     false,
					Membership: Membership{
						Included: []GroupMember{
							{Id: "sddcId1"},
							{Id: "sddcId2"},
						},
					},
				},
				error: nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		sddcGroup, err := sddcGroupClient.GetSddcGroup(testCase.input.groupId)
		assert.Equal(t, testCase.output.sddcGroup, sddcGroup)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestCreateSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcUrl, testVmcUrl)
	type inputStruct struct {
		httpClientStub HttpClient
		name           string
		description    string
		sddcIds        *[]string
	}
	type outputStruct struct {
		groupId string
		taskId  string
		error   error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJson:   "",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIds:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupId: "",
				taskId:  "",
				error:   fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJson: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"sddcId1\"\n]\n}\n]\n}",
					t: t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIds:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupId: "",
				taskId:  "",
				error:   fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: sddcId1 "),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusUnauthorized,
					responseError:  nil,
					responseJson:   "",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIds:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupId: "",
				taskId:  "",
				error:   fmt.Errorf("CreateSddcGroup response code: 401"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod: http.MethodPost,
					expectedJson:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedUrl:    "https://test.vmc.vmware.com/api/network/testOrgId/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJson:   "{\"config_id\":\"configId\",\"group_id\":\"newGroupId\",\"operation_id\":\"createSddcGroupTaskId\"}",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIds:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupId: "newGroupId",
				taskId:  "createSddcGroupTaskId",
				error:   nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		groupId, taskId, err := sddcGroupClient.CreateSddcGroup(
			testCase.input.name, testCase.input.description, testCase.input.sddcIds)
		assert.Equal(t, testCase.output.groupId, groupId)
		assert.Equal(t, testCase.output.taskId, taskId)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestUpdateSddcGroupMembers(t *testing.T) {
	t.Setenv(constants.VmcUrl, testVmcUrl)
	type inputStruct struct {
		httpClientStub  HttpClient
		groupId         string
		sddcIdsToAdd    *[]string
		sddcIdsToRemove *[]string
	}
	type outputStruct struct {
		taskId string
		error  error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\"," +
						"\"resource_type\":\"network-connectivity-config\",\"type\":\"UPDATE_MEMBERS\",\"config\":{" +
						"\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"add_members\":[{\"id\":\"sddcId1\"}]," +
						"\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJson:  "",
					t:             t,
				},
				groupId:         "testGroupId",
				sddcIdsToAdd:    &[]string{"sddcId1"},
				sddcIdsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskId: "",
				error:  fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\"," +
						"\"resource_type\":\"network-connectivity-config\",\"type\":\"UPDATE_MEMBERS\"," +
						"\"config\":{\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"add_members\":[{\"id\":\"sddcId1\"}]," +
						"\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusConflict,
					responseError: nil,
					responseJson: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"sddcId1\"\n]\n}\n]\n}",
					t: t,
				},
				groupId:         "testGroupId",
				sddcIdsToAdd:    &[]string{"sddcId1"},
				sddcIdsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskId: "",
				error:  fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: sddcId1 "),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"UPDATE_MEMBERS\",\"config\":{\"type\":\"AwsUpdateDeploymentGroupMembersConfig\"," +
						"\"add_members\":[{\"id\":\"sddcId1\"}],\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusOK,
					responseError: nil,
					responseJson: "{\"task_id\":\"notImportant\",\"org_id\":\"testOrgId\",\"resource_id\":\"testGroupId\",\"type\":\"UPDATE_MEMBERS\",\"config\":{" +
						"\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"operation_id\":\"updateSddcMembersTaskId\",\"add_members\":[{\"id\":\"sddcId1\"}],\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					t: t,
				},
				groupId:         "testGroupId",
				sddcIdsToAdd:    &[]string{"sddcId1"},
				sddcIdsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskId: "updateSddcMembersTaskId",
				error:  nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		taskId, err := sddcGroupClient.UpdateSddcGroupMembers(
			testCase.input.groupId, testCase.input.sddcIdsToAdd, testCase.input.sddcIdsToRemove)
		assert.Equal(t, testCase.output.taskId, taskId)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestDeleteSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcUrl, testVmcUrl)
	type inputStruct struct {
		httpClientStub HttpClient
		groupId        string
	}
	type outputStruct struct {
		taskId string
		error  error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJson:  "",
					t:             t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				taskId: "",
				error:  fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusConflict,
					responseError: nil,
					responseJson: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deletion of SDDC group failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"cannot delete SDDC Group. There are still sddcs attached\"\n}]}",
					t: t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				taskId: "",
				error:  fmt.Errorf("cannot delete SDDC Group. There are still sddcs attached"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HttpClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIdRequestJson: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJson: "{\"org_id\":\"testOrgId\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedUrl:   "https://test.vmc.vmware.com/api/network/testOrgId/aws/operations",
					responseCode:  http.StatusCreated,
					responseError: nil,
					responseJson: "{\"id\":\"deleteSddcTaskId\",\"task_id\":\"notImportant\",\"org_id\":\"testOrgId\",\"resource_id\":\"testGroupId\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					t: t,
				},
				groupId: "testGroupId",
			},
			outputStruct{
				taskId: "deleteSddcTaskId",
				error:  nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		taskId, err := sddcGroupClient.DeleteSddcGroup(testCase.input.groupId)
		assert.Equal(t, testCase.output.taskId, taskId)
		assert.Equal(t, testCase.output.error, err)
	}
}
