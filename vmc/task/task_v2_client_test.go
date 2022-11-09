/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package task

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

const testAccessToken = "testAccessToken"
const testVmcURL = "https://test.vmc.vmware.com"

type HTTPClientStub struct {
	expectedURL   string
	responseJSON  string
	responseCode  int
	responseError error
	t             *testing.T
}

func (stub HTTPClientStub) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(stub.t, stub.expectedURL, req.URL.String())
	assert.Equal(stub.t, http.MethodGet, req.Method)
	assert.Equal(stub.t, req.Header.Get(authnHeader), testAccessToken)
	response := http.Response{
		StatusCode: stub.responseCode,
		Body:       io.NopCloser(strings.NewReader(stub.responseJSON)),
	}
	return &response, stub.responseError
}

func TestGetTask(t *testing.T) {
	testOrgID := "testOrgId"
	expectedURL := "https://test.vmc.vmware.com/api/operation/testOrgId/core/operations/lele"
	type inputStruct struct {
		httpClientStub HTTPClient
		taskID         string
	}
	type outputStruct struct {
		task  V2Task
		error error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: HTTPClientStub{
					expectedURL:   expectedURL,
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJSON:  "",
					t:             t,
				},
				taskID: "lele",
			},
			outputStruct{
				task:  V2Task{},
				error: fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: HTTPClientStub{
					expectedURL:   expectedURL,
					responseCode:  http.StatusOK,
					responseError: nil,
					responseJSON:  "{\"error_message\":\"SddcGroup creation failed\", \"state\":{\"name\":\"FAILED\"}}",
					t:             t,
				},
				taskID: "lele",
			},
			outputStruct{
				task: V2Task{
					TaskState: V2State{
						Name: "FAILED",
					},
					ErrorMessage: "SddcGroup creation failed",
				},
				error: nil,
			},
		},
		{
			inputStruct{
				httpClientStub: HTTPClientStub{
					expectedURL:   expectedURL,
					responseCode:  http.StatusNotFound,
					responseError: nil,
					responseJSON:  "",
					t:             t,
				},
				taskID: "lele",
			},
			outputStruct{
				task:  V2Task{},
				error: fmt.Errorf("Task with ID: lele not found "),
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestV2ClientImpl(testVmcURL, testOrgID, testAccessToken, testCase.input.httpClientStub)
		task, err := sddcGroupClient.GetTask(testCase.input.taskID)
		assert.Equal(t, testCase.output.task, task)
		assert.Equal(t, testCase.output.error, err)
	}
}
