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
const testVmcUrl = "https://test.vmc.vmware.com"

type HttpClientStub struct {
	expectedUrl   string
	responseJson  string
	responseCode  int
	responseError error
	t             *testing.T
}

func (stub HttpClientStub) Do(req *http.Request) (*http.Response, error) {
	assert.Equal(stub.t, stub.expectedUrl, req.URL.String())
	assert.Equal(stub.t, http.MethodGet, req.Method)
	assert.Equal(stub.t, req.Header.Get(authnHeader), testAccessToken)
	response := http.Response{
		StatusCode: stub.responseCode,
		Body:       io.NopCloser(strings.NewReader(stub.responseJson)),
	}
	return &response, stub.responseError
}

func TestGetTask(t *testing.T) {
	testOrgId := "testOrgId"
	expectedUrl := "https://test.vmc.vmware.com/api/operation/testOrgId/core/operations/lele"
	type inputStruct struct {
		httpClientStub HttpClient
		taskId         string
	}
	type outputStruct struct {
		task  TaskV2
		error error
	}
	type test struct {
		input  inputStruct
		output outputStruct
	}
	tests := []test{
		{
			inputStruct{
				httpClientStub: HttpClientStub{
					expectedUrl:   expectedUrl,
					responseCode:  http.StatusInternalServerError,
					responseError: fmt.Errorf("VMC service down"),
					responseJson:  "",
					t:             t,
				},
				taskId: "lele",
			},
			outputStruct{
				task:  TaskV2{},
				error: fmt.Errorf("VMC service down"),
			},
		},
		{
			inputStruct{
				httpClientStub: HttpClientStub{
					expectedUrl:   expectedUrl,
					responseCode:  http.StatusOK,
					responseError: nil,
					responseJson:  "{\"error_message\":\"SddcGroup creation failed\", \"state\":{\"name\":\"FAILED\"}}",
					t:             t,
				},
				taskId: "lele",
			},
			outputStruct{
				task: TaskV2{
					TaskState: TaskV2State{
						Name: "FAILED",
					},
					ErrorMessage: "SddcGroup creation failed",
				},
				error: nil,
			},
		},
		{
			inputStruct{
				httpClientStub: HttpClientStub{
					expectedUrl:   expectedUrl,
					responseCode:  http.StatusNotFound,
					responseError: nil,
					responseJson:  "",
					t:             t,
				},
				taskId: "lele",
			},
			outputStruct{
				task:  TaskV2{},
				error: fmt.Errorf("Task with ID: lele not found "),
			},
		},
	}
	for _, testCase := range tests {
		sddcGroupClient := newTestTaskV2ClientImpl(testVmcUrl, testOrgId, testAccessToken, testCase.input.httpClientStub)
		task, err := sddcGroupClient.GetTask(testCase.input.taskId)
		assert.Equal(t, testCase.output.task, task)
		assert.Equal(t, testCase.output.error, err)
	}
}
