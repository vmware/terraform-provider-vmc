/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package task

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
	vmcURL       string
	cspURL       string
	refreshToken string
	orgID        string
	accessToken  string
	HTTPClient   HTTPClient
}

func NewV2ClientImpl(vmcURL string, cspURL string, refreshToken string, orgID string) *V2ClientImpl {
	return &V2ClientImpl{
		vmcURL:       vmcURL,
		cspURL:       cspURL,
		refreshToken: refreshToken,
		orgID:        orgID,
		HTTPClient:   http.DefaultClient,
	}
}

// newTestV2ClientImpl intended for injecting dummy accessToken and stubbed HTTPClient for
// testing purposes.
func newTestV2ClientImpl(vmcURL string, orgID string, accessToken string, httpClient HTTPClient) *V2ClientImpl {
	return &V2ClientImpl{
		vmcURL:      vmcURL,
		orgID:       orgID,
		accessToken: accessToken,
		HTTPClient:  httpClient,
	}
}

// Authenticate grab an access token and set it into the SddcGroupClient instance for later use
func (client *V2ClientImpl) Authenticate() error {
	authURL := client.cspURL + constants.CspRefreshURLSuffix
	securityContext, err := connector.SecurityContextByRefreshToken(client.refreshToken, authURL)
	if err != nil {
		return err
	}
	client.accessToken = securityContext.Property(security.ACCESS_TOKEN).(string)
	return nil
}

func (client *V2ClientImpl) GetTask(taskID string) (V2Task, error) {
	getTaskV2URL := client.getBaseURL() + fmt.Sprintf("/operation/%s/core/operations/%s", client.orgID, taskID)
	req, _ := http.NewRequest(http.MethodGet, getTaskV2URL, nil)
	req.Header.Add(authnHeader, client.accessToken)
	var result V2Task
	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return result, err
	}
	if statusCode == http.StatusNotFound {
		return result, fmt.Errorf("Task with ID: %s not found ", taskID)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result, err
	}
	return result, fmt.Errorf("GetTask response code: %d", statusCode)
}

func (client *V2ClientImpl) getBaseURL() string {
	return client.vmcURL + "/api"
}

// executeRequest Returns the body of the response as byte array pointer, the status code
// or any error that may have occurred during the Http communication.
func (client *V2ClientImpl) executeRequest(
	request *http.Request) (responseBody *[]byte, statusCode int, error error) {
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
