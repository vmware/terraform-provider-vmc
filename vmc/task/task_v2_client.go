// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
)

const authnHeader = "csp-auth-token"

type V2State struct {
	Name string `json:"name"`
}

type V2Task struct {
	ID           string  `json:"id"`
	TaskState    V2State `json:"state"`
	TaskType     string  `json:"type"`
	ErrorMessage string  `json:"error_message"`
}

type V2Client interface {
	connector.Authenticator
	GetTask(taskID string) (V2Task, error)
}

// HTTPClient an interface, that is implemented by the http.DefaultClient,
// intended to enable stubbing for testing purposes
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type V2ClientImpl struct {
	connector  connector.Wrapper
	HTTPClient HTTPClient
}

func NewV2ClientImpl(wrapper connector.Wrapper) *V2ClientImpl {
	copyWrapper := connector.CopyWrapper(wrapper)
	return &V2ClientImpl{
		connector:  *copyWrapper,
		HTTPClient: http.DefaultClient,
	}
}

// newTestV2ClientImpl intended for injecting dummy accessToken and stubbed HTTPClient for
// testing purposes.
func newTestV2ClientImpl(vmcURL string, orgID string, accessToken string, httpClient HTTPClient) *V2ClientImpl {
	testConnector := connector.Wrapper{
		VmcURL: vmcURL,
		OrgID:  orgID,
	}
	// Create a dummy connector to house the access token in a security context
	testConnector.Connector = client.NewConnector("", client.WithHttpClient(&http.Client{}),
		client.WithSecurityContext(security.NewOauthSecurityContext(accessToken)))
	return &V2ClientImpl{
		connector:  testConnector,
		HTTPClient: httpClient,
	}
}

// Authenticate grab an access token and set it into the V2Client instance for later use
func (client *V2ClientImpl) Authenticate() error {
	return client.connector.Authenticate()
}

func (client *V2ClientImpl) GetTask(taskID string) (V2Task, error) {
	getTaskV2URL := client.getBaseURL() + fmt.Sprintf("/operation/%s/core/operations/%s", client.connector.OrgID, taskID)
	req, _ := http.NewRequest(http.MethodGet, getTaskV2URL, nil)
	req.Header.Add(authnHeader, client.connector.Connector.SecurityContext().Property(security.ACCESS_TOKEN).(string))
	var result V2Task
	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return result, err
	}
	if statusCode == http.StatusNotFound {
		return result, fmt.Errorf("task with ID: %s not found ", taskID)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result, err
	}
	return result, fmt.Errorf("GetTask response code: %d", statusCode)
}

func (client *V2ClientImpl) getBaseURL() string {
	return client.connector.VmcURL + "/api"
}

// executeRequest Returns the body of the response as byte array pointer, the status code
// or any error that may have occurred during the Http communication.
func (client *V2ClientImpl) executeRequest(
	request *http.Request) (responseBody *[]byte, statusCode int, responseErr error) {
	response, err := client.HTTPClient.Do(request)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Eror closing body of http response")
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
