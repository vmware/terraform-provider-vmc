/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package sddcgroup

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
const testOrgID = "testOrgID"
const testVmcURL = "https://test.vmc.vmware.com"

type HTTPClientStub struct {
	expectedJSON                         string
	additionalResourceIDRequestJSON      string
	expectedMethod                       string
	expectedURL                          string
	responseJSON                         string
	additionalResourceTraitsResponseJSON string
	responseCode                         int
	responseError                        error
	t                                    *testing.T
}

func (stub *HTTPClientStub) Do(req *http.Request) (*http.Response, error) {
	// Handle getResourceIDFromGroupID preliminary request where resource ID is derived from provided group ID
	if strings.HasPrefix(req.URL.String(),
		"https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/?group_id=") {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(stub.additionalResourceIDRequestJSON)),
		}, nil
	}
	// Handle second request where traits are derived from provided resource ID
	if strings.HasPrefix(req.URL.String(),
		"https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/resourceIdDifferentFromGroupId/?trait=") {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader(stub.additionalResourceTraitsResponseJSON)),
		}, nil
	}
	if req.Body == nil {
		assert.Equal(stub.t, stub.expectedJSON, "")
	} else {
		assert.Equal(stub.t, stub.expectedJSON, readAsString(req.Body))
	}
	assert.Equal(stub.t, stub.expectedURL, req.URL.String())
	assert.Equal(stub.t, stub.expectedMethod, req.Method)
	assert.Equal(stub.t, req.Header.Get(authnHeader), testAccessToken)
	response := http.Response{
		StatusCode: stub.responseCode,
		Body:       io.NopCloser(strings.NewReader(stub.responseJSON)),
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
	type inputStruct struct {
		httpClientStub HTTPClient
		sddcIds        *[]string
	}
	type test struct {
		input inputStruct
		want  error
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJSON:   "",
					t:              t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("VMC service down"),
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJSON:   "",
					t:              t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			nil,
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJSON: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"," +
						"\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"male\"\n]\n}\n]\n}",
					t: t,
				},
				sddcIds: &[]string{"lele", "male"},
			},
			fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: male "),
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		assert.Equal(t, testCase.want, sddcGroupClient.ValidateCreateSddcGroup(testCase.input.sddcIds))
	}
}

func TestValidateCreateSddcGroupMembers(t *testing.T) {
	t.Setenv(constants.VmcURL, testVmcURL)
	type inputStruct struct {
		httpClientStub HTTPClient
		groupID        string
		sddcIDs        *[]string
	}
	type test struct {
		input inputStruct
		want  error
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJSON:   "",
					t:              t,
				},
				groupID: "testGroupId",
				sddcIDs: &[]string{"lele", "male"},
			},
			fmt.Errorf("VMC service down"),
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJSON:   "",
					t:              t,
				},
				groupID: "testGroupId",
				sddcIDs: &[]string{"lele", "male"},
			},
			nil,
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"deployment_group_id\":\"testGroupId\",\"members\":[{\"id\":\"lele\"},{\"id\":\"male\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/validate-members",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJSON: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"," +
						"\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"male\"\n]\n}\n]\n}",
					t: t,
				},
				groupID: "testGroupId",
				sddcIDs: &[]string{"lele", "male"},
			},
			fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: male "),
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		assert.Equal(t, testCase.want, sddcGroupClient.ValidateUpdateSddcGroupMembers(testCase.input.groupID, testCase.input.sddcIDs))
	}
}

func TestGetSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcURL, testVmcURL)
	type inputStruct struct {
		httpClientStub HTTPClient
		groupID        string
	}
	type outputStruct struct {
		sddcGroup                 *DeploymentGroup
		networkConnectivityConfig *NetworkConnectivityConfig
		error                     error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodGet,
					expectedJSON:   "",
					expectedURL:    "https://test.vmc.vmware.com/api/inventory/testOrgID/core/deployment-groups/testGroupId",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJSON:   "",
					t:              t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				sddcGroup:                 nil,
				networkConnectivityConfig: nil,
				error:                     fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodGet,
					expectedJSON:   "",
					expectedURL:    "https://test.vmc.vmware.com/api/inventory/testOrgID/core/deployment-groups/testGroupId",
					responseCode:   http.StatusUnauthorized,
					responseError:  nil,
					responseJSON:   "",
					t:              t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				sddcGroup:                 nil,
				networkConnectivityConfig: nil,
				error:                     fmt.Errorf("Unauthenticated request "),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodGet,
					expectedJSON:   "",
					expectedURL:    "https://test.vmc.vmware.com/api/inventory/testOrgID/core/deployment-groups/testGroupId",
					responseCode:   http.StatusNotFound,
					responseError:  nil,
					responseJSON:   "",
					t:              t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				sddcGroup:                 nil,
				networkConnectivityConfig: nil,
				error:                     fmt.Errorf("SDDC Group with ID: testGroupId not found"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodGet,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON:                    "",
					expectedURL:                     "https://test.vmc.vmware.com/api/inventory/testOrgID/core/deployment-groups/testGroupId",
					responseCode:                    http.StatusOK,
					responseError:                   nil,
					responseJSON: "{\n \"id\": \"groupID\",\n  \"name\": \"test_sddc_group\",\n  \"org_id\": \"testOrgID\",\n  \"description\": \"sdfa\",\n " +
						"\"membership\": {\n \"included\": [\n {\n\"deployment_id\": \"sddcId1\"\n },\n {\n \"deployment_id\": \"sddcId2\"\n }\n ]\n },\n \"deleted\": false\n}",
					additionalResourceTraitsResponseJSON: "{\n\"id\":\"1ed56a5b-e89d-6581-91ec-f5243e2d14a5\",\n\"" +
						"group_id\":\"1ed56a5b-e825-6b15-a726-a9eef692a3e2\",\n\"name\":\"test_sddc_group\",\n\"" +
						"state\":{\n\"name\":\"CONNECTED\",\n\"error_msgs\":{}\n},\n\"traits\":{\n\"AwsNetworkConnectivityTrait\":{" +
						"\n\"l3connectors\":[\n{\n\"id\":\"tgw-06794c465682ef15d\",\n\"region\":\"us-west-2\"\n}\n]\n},\n\"" +
						"AwsVpcAttachmentsTrait\":{\n\"accounts\":[\n{\n\"account_number\":\"123420333\",\n\"resource_share_name\":" +
						"\"ram_share_id_xxx\",\n\"state\":\"CONNECTED\",\n\"attachments\":[\n{\n\"vpc_id\":\"pl-0da219e504e5dc6e8\"," +
						"\n\"state\":\"CONNECTED\",\n\"attach_id\":\"attach_id_123\",\n\"configured_prefixes\":[\n\"10.20.30.40/24\"," +
						"\n\"40.30.20.10/24\"\n]\n}\n]\n}\n]\n},\n\"AwsDirectConnectGatewayAssociationsTrait\":{\n\"" +
						"direct_connect_gateway_associations\":[\n{\n\"direct_connect_gateway_id\":\"direct_connect_gateway_id_123\",\n" +
						"\"direct_connect_gateway_owner\":\"Boba\",\n\"state\":\"CONNECTED\",\n\"peering_regions\":[\n{\n\"" +
						"allowed_prefixes\":[\n\"10.20.30.40/24\",\n\"40.30.20.10/24\"\n]\n}\n]\n}\n]\n},\n\"" +
						"AwsCustomerTransitGatewayAssociationsTrait\":{\n\"customer_transit_gateway_associations\":[\n{\n\"" +
						"customer_transit_gateway_id\":\"customer_transit_gateway_id_123\",\n\"customer_transit_gateway_owner\":\"Fett\"," +
						"\n\"customer_transit_gateway_region\":{\n\"code\":\"us-east-1\"\n},\n\"peering_regions\":[\n{\n\"" +
						"configured_prefixes\":[\n\"10.20.30.40/24\",\n\"40.30.20.10/24\"\n]\n}\n]\n}\n]\n}\n}\n}",
					t: t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				sddcGroup: &DeploymentGroup{
					ID:          "groupID",
					Name:        "test_sddc_group",
					OrgID:       testOrgID,
					Description: "sdfa",
					Deleted:     false,
					Membership: Membership{
						Included: []GroupMember{
							{ID: "sddcId1"},
							{ID: "sddcId2"},
						},
					},
				},
				networkConnectivityConfig: &NetworkConnectivityConfig{
					ID:      "1ed56a5b-e89d-6581-91ec-f5243e2d14a5",
					GroupID: "1ed56a5b-e825-6b15-a726-a9eef692a3e2",
					Name:    "test_sddc_group",
					NetworkConnectivityConfigState: NetworkConnectivityConfigState{
						Name: "CONNECTED",
					},
					Traits: &Traits{
						TransitGateway: &AwsNetworkConnectivityTrait{
							L3Connectors: []L3Connector{
								{
									ID:     "tgw-06794c465682ef15d",
									Region: "us-west-2"},
							},
						},
						AwsInfo: &AwsVpcAttachmentsTrait{
							Accounts: []AwsAccount{
								{AccountNumber: "123420333",
									RAMShareID: "ram_share_id_xxx",
									Status:     "CONNECTED",
									AccountAttachments: []AccountAttachment{
										{
											VpcID:        "pl-0da219e504e5dc6e8",
											State:        "CONNECTED",
											AttachmentID: "attach_id_123",
											StaticRoutes: []string{
												"10.20.30.40/24",
												"40.30.20.10/24",
											}},
									}},
							},
						},
						DxGateway: &AwsDirectConnectGatewayAssociationsTrait{
							DirectConnectGatewayAssociations: []DirectConnectGatewayAssociation{
								{
									DxgwID:    "direct_connect_gateway_id_123",
									DxgwOwner: "Boba",
									Status:    "CONNECTED",
									PeeringRegions: []PeeringRegions{
										{
											AllowedPrefixes: []string{
												"10.20.30.40/24",
												"40.30.20.10/24",
											},
										},
									},
								},
							},
						},
						ExternalTgw: &AwsCustomerTransitGatewayAssociationsTrait{
							CustomerTransitGatewayAssociations: []CustomerTransitGatewayAssociation{
								{
									TgwID:    "customer_transit_gateway_id_123",
									TgwOwner: "Fett",
									TgwRegion: TgwRegion{
										Region: "us-east-1",
									},
									PeeringRegions: []PeeringRegions{
										{
											ConfiguredPrefixes: []string{
												"10.20.30.40/24",
												"40.30.20.10/24",
											},
										},
									},
								},
							},
						},
					},
				},
				error: nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		sddcGroup, networkConnectivityConfig, err := sddcGroupClient.GetSddcGroup(testCase.input.groupID)
		assert.Equal(t, testCase.output.sddcGroup, sddcGroup)
		assert.Equal(t, testCase.output.networkConnectivityConfig, networkConnectivityConfig)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestCreateSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcURL, testVmcURL)
	type inputStruct struct {
		httpClientStub HTTPClient
		name           string
		description    string
		sddcIDs        *[]string
	}
	type outputStruct struct {
		groupID string
		taskID  string
		error   error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusInternalServerError,
					responseError:  fmt.Errorf("VMC service down"),
					responseJSON:   "",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIDs:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupID: "",
				taskID:  "",
				error:   fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusConflict,
					responseError:  nil,
					responseJSON: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"sddcId1\"\n]\n}\n]\n}",
					t: t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIDs:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupID: "",
				taskID:  "",
				error:   fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: sddcId1 "),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusBadRequest,
					responseError:  nil,
					responseJSON:   "",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIDs:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupID: "",
				taskID:  "",
				error:   fmt.Errorf("CreateSddcGroup response code: 400"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod: http.MethodPost,
					expectedJSON:   "{\"name\":\"testGroupName\",\"description\":\"testGroupDescription\",\"members\":[{\"id\":\"sddcId1\"},{\"id\":\"sddcId2\"}]}",
					expectedURL:    "https://test.vmc.vmware.com/api/network/testOrgID/core/network-connectivity-configs/create-group-network-connectivity",
					responseCode:   http.StatusOK,
					responseError:  nil,
					responseJSON:   "{\"config_id\":\"configId\",\"group_id\":\"newGroupId\",\"operation_id\":\"createSddcGroupTaskId\"}",
					t:              t,
				},
				name:        "testGroupName",
				description: "testGroupDescription",
				sddcIDs:     &[]string{"sddcId1", "sddcId2"},
			},
			outputStruct{
				groupID: "newGroupId",
				taskID:  "createSddcGroupTaskId",
				error:   nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		groupID, taskID, err := sddcGroupClient.CreateSddcGroup(
			testCase.input.name, testCase.input.description, testCase.input.sddcIDs)
		assert.Equal(t, testCase.output.groupID, groupID)
		assert.Equal(t, testCase.output.taskID, taskID)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestUpdateSddcGroupMembers(t *testing.T) {
	t.Setenv(constants.VmcURL, testVmcURL)
	type inputStruct struct {
		httpClientStub  HTTPClient
		groupID         string
		sddcIDsToAdd    *[]string
		sddcIDsToRemove *[]string
	}
	type outputStruct struct {
		taskID string
		error  error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\"," +
						"\"resource_type\":\"network-connectivity-config\",\"type\":\"UPDATE_MEMBERS\",\"config\":{" +
						"\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"add_members\":[{\"id\":\"sddcId1\"}]," +
						"\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJSON:  "",
					t:             t,
				},
				groupID:         "testGroupId",
				sddcIDsToAdd:    &[]string{"sddcId1"},
				sddcIDsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskID: "",
				error:  fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\"," +
						"\"resource_type\":\"network-connectivity-config\",\"type\":\"UPDATE_MEMBERS\"," +
						"\"config\":{\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"add_members\":[{\"id\":\"sddcId1\"}]," +
						"\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusConflict,
					responseError: nil,
					responseJSON: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deployment(s) for group membership failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"Found invalid or overlapping CIDR blocks.\",\n\"members\": [\n\"sddcId1\"\n]\n}\n]\n}",
					t: t,
				},
				groupID:         "testGroupId",
				sddcIDsToAdd:    &[]string{"sddcId1"},
				sddcIDsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskID: "",
				error:  fmt.Errorf("Found invalid or overlapping CIDR blocks. For members: sddcId1 "),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"UPDATE_MEMBERS\",\"config\":{\"type\":\"AwsUpdateDeploymentGroupMembersConfig\"," +
						"\"add_members\":[{\"id\":\"sddcId1\"}],\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusOK,
					responseError: nil,
					responseJSON: "{\"task_id\":\"notImportant\",\"org_id\":\"testOrgID\",\"resource_id\":\"testGroupId\",\"type\":\"UPDATE_MEMBERS\",\"config\":{" +
						"\"type\":\"AwsUpdateDeploymentGroupMembersConfig\",\"operation_id\":\"updateSddcMembersTaskId\",\"add_members\":[{\"id\":\"sddcId1\"}],\"remove_members\":[{\"id\":\"sddcId2\"}]}}",
					t: t,
				},
				groupID:         "testGroupId",
				sddcIDsToAdd:    &[]string{"sddcId1"},
				sddcIDsToRemove: &[]string{"sddcId2"},
			},
			outputStruct{
				taskID: "updateSddcMembersTaskId",
				error:  nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		taskID, err := sddcGroupClient.UpdateSddcGroupMembers(
			testCase.input.groupID, testCase.input.sddcIDsToAdd, testCase.input.sddcIDsToRemove)
		assert.Equal(t, testCase.output.taskID, taskID)
		assert.Equal(t, testCase.output.error, err)
	}
}

func TestDeleteSddcGroup(t *testing.T) {
	t.Setenv(constants.VmcURL, testVmcURL)
	type inputStruct struct {
		httpClientStub HTTPClient
		groupID        string
	}
	type outputStruct struct {
		taskID string
		error  error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJSON:  "",
					t:             t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				taskID: "",
				error:  fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusConflict,
					responseError: nil,
					responseJSON: "{\n\"status\": 409,\n\"error\": \"Conflict\",\n\"message\": \"Validation of deletion of SDDC group failed.\"" +
						",\n\"details\": [\n{\n\"validation_error_message\": \"cannot delete SDDC Group. There are still sddcs attached\"\n}]}",
					t: t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				taskID: "",
				error:  fmt.Errorf("cannot delete SDDC Group. There are still sddcs attached"),
			},
		},
		{
			inputStruct{
				httpClientStub: &HTTPClientStub{
					expectedMethod:                  http.MethodPost,
					additionalResourceIDRequestJSON: "[{\"id\":\"resourceIdDifferentFromGroupId\"}]",
					expectedJSON: "{\"org_id\":\"testOrgID\",\"resource_id\":\"resourceIdDifferentFromGroupId\",\"resource_type\":\"network-connectivity-config\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					expectedURL:   "https://test.vmc.vmware.com/api/network/testOrgID/aws/operations",
					responseCode:  http.StatusCreated,
					responseError: nil,
					responseJSON: "{\"id\":\"deleteSddcTaskId\",\"task_id\":\"notImportant\",\"org_id\":\"testOrgID\",\"resource_id\":\"testGroupId\"," +
						"\"type\":\"DELETE_DEPLOYMENT_GROUP\",\"config\":{\"type\":\"AwsDeleteDeploymentGroupConfig\"}}",
					t: t,
				},
				groupID: "testGroupId",
			},
			outputStruct{
				taskID: "deleteSddcTaskId",
				error:  nil,
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestSddcGroupClient(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		taskID, err := sddcGroupClient.DeleteSddcGroup(testCase.input.groupID)
		assert.Equal(t, testCase.output.taskID, taskID)
		assert.Equal(t, testCase.output.error, err)
	}
}
