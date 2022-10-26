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

type TaskV2State struct {
	Name string `json:"name"`
}

type TaskV2 struct {
	Id           string      `json:"id"`
	TaskState    TaskV2State `json:"state"`
	TaskType     string      `json:"type"`
	ErrorMessage string      `json:"error_message"`
}

type TaskV2Client interface {
	connector.Authenticator
	GetTask(taskId string) (TaskV2, error)
}

// HttpClient an interface, that is implemented by the http.DefaultClient,
// intended to enable stubbing for testing purposes
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type TaskV2ClientImpl struct {
	vmcUrl       string
	cspUrl       string
	refreshToken string
	orgId        string
	accessToken  string
	httpClient   HttpClient
}

func NewTaskV2ClientImpl(vmcUrl string, cspUrl string, refreshToken string, orgId string) *TaskV2ClientImpl {
	return &TaskV2ClientImpl{
		vmcUrl:       vmcUrl,
		cspUrl:       cspUrl,
		refreshToken: refreshToken,
		orgId:        orgId,
		httpClient:   http.DefaultClient,
	}
}

// newTestTaskV2ClientImpl intended for injecting dummy accessToken and stubbed httpClient for
// testing purposes.
func newTestTaskV2ClientImpl(vmcUrl string, orgId string, accessToken string, httpClient HttpClient) *TaskV2ClientImpl {
	return &TaskV2ClientImpl{
		vmcUrl:      vmcUrl,
		orgId:       orgId,
		accessToken: accessToken,
		httpClient:  httpClient,
	}
}

// Authenticate grab an access token and set it into the SddcGroupClient instance for later use
func (client *TaskV2ClientImpl) Authenticate() error {
	authUrl := client.cspUrl + constants.CspRefreshUrlSuffix
	securityContext, err := connector.SecurityContextByRefreshToken(client.refreshToken, authUrl)
	if err != nil {
		return err
	}
	client.accessToken = securityContext.Property(security.ACCESS_TOKEN).(string)
	return nil
}

func (client *TaskV2ClientImpl) GetTask(taskId string) (TaskV2, error) {
	getTaskV2Url := client.getBaseUrl() + fmt.Sprintf("/operation/%s/core/operations/%s", client.orgId, taskId)
	req, _ := http.NewRequest(http.MethodGet, getTaskV2Url, nil)
	req.Header.Add(authnHeader, client.accessToken)
	var result TaskV2
	rawResponse, statusCode, err := client.executeRequest(req)
	if err != nil {
		return result, err
	}
	if statusCode == http.StatusNotFound {
		return result, fmt.Errorf("Task with ID: %s not found ", taskId)
	}
	if statusCode == http.StatusOK {
		err := json.NewDecoder(bytes.NewReader(*rawResponse)).Decode(&result)
		return result, err
	}
	return result, fmt.Errorf("GetTask response code: %d", statusCode)
}

func (client *TaskV2ClientImpl) getBaseUrl() string {
	return client.vmcUrl + "/api"
}

// executeRequest Returns the body of the response as byte array pointer, the status code
// or any error that may have occurred during the Http communication.
func (client *TaskV2ClientImpl) executeRequest(
	request *http.Request) (responseBody *[]byte, statusCode int, error error) {
	response, err := client.httpClient.Do(request)
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
