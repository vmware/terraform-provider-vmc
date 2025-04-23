// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
)

func TestKeyedMutexLock(t *testing.T) {
	var keyedMutex = KeyedMutex{}
	var key1 = "key1"
	var key2 = "key2"

	lock1Obtained := false
	lock2Obtained := false

	var unlockFunction = keyedMutex.Lock(key1)
	// Try to flip the flag 1 in a separate thread
	go func() {
		keyedMutex.Lock(key1)
		lock1Obtained = true
	}()
	// Try to flip the flag 2 in a separate thread
	go func() {
		keyedMutex.Lock(key2)
		lock2Obtained = true
	}()

	// Give enough time for the separate threads to try to flip the flags
	time.Sleep(500 * time.Millisecond)
	assert.False(t, lock1Obtained)
	assert.True(t, lock2Obtained)
	// Test the unlock functionality
	unlockFunction()
	time.Sleep(500 * time.Millisecond)
	assert.True(t, lock1Obtained)
}

type AuthenticatorStub struct {
}

func (stub AuthenticatorStub) Authenticate() error {
	return nil
}

type BrokenAuthenticatorStub struct {
}

func (stub BrokenAuthenticatorStub) Authenticate() error {
	return fmt.Errorf("authentication broken")
}

func TestRetryTaskUntilFinished(t *testing.T) {
	type inputStruct struct {
		connectorWrapper connector.Authenticator
		taskSupplier     func() (model.Task, error)
		errorMessage     string
		finishCallback   func(task model.Task)
	}
	type test struct {
		input inputStruct
		want  *retry.RetryError
	}
	var finishCallbackHasBeenCalled = false
	tests := []test{
		// Unauthenticated handling - retry authentication
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					return model.Task{}, errors.Unauthenticated{}
				},
				errorMessage: "",
				finishCallback: func(_ model.Task) {
					assert.Fail(t, "finishCallback should not be called on retrievable errors")
				},
			},
			want: retry.RetryableError(fmt.Errorf("task still in progress")),
		},
		// Unauthenticated handling - fail
		{
			input: inputStruct{
				connectorWrapper: BrokenAuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					return model.Task{Id: "Unauthenticated handling - fail"}, errors.Unauthenticated{}
				},
				errorMessage: "",
				finishCallback: func(task model.Task) {
					assert.Equal(t, "Unauthenticated handling - fail", task.Id)
				},
			},
			want: retry.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : authentication broken")),
		},
		// Service unavailable retry
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					// Set the global counter to the last acceptable value
					serviceUnavailableRetries = 19
					return model.Task{}, errors.ServiceUnavailable{}
				},
				errorMessage: "",
				finishCallback: func(_ model.Task) {
					assert.Fail(t, "finishCallback should not be called on retrievable errors")
				},
			},
			want: retry.RetryableError(fmt.Errorf(
				"VMC backend is experiencing difficulties, retry 20 from 20 to polling the SDDC Create Task")),
		},
		// Service unavailable fail
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					return model.Task{Id: "Service unavailable fail"}, errors.ServiceUnavailable{}
				},
				errorMessage: "",
				finishCallback: func(task model.Task) {
					assert.Equal(t, "Service unavailable fail", task.Id)
				},
			},
			want: retry.NonRetryableError(fmt.Errorf("max ServiceUnavailable retries (20) reached to create SDDC")),
		},
		// Task status failed
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					status := model.Task_STATUS_FAILED
					taskErrorMessage := "mnogoGrumna"
					return model.Task{Status: &status, ErrorMessage: &taskErrorMessage}, nil
				},
				errorMessage: "Cluster creation failed",
				finishCallback: func(task model.Task) {
					// test that service unavailable retries have been reset
					assert.Equal(t, 0, serviceUnavailableRetries)
					assert.Equal(t, model.Task_STATUS_FAILED, *task.Status)
				},
			},
			want: retry.NonRetryableError(fmt.Errorf("task failed: Cluster creation failed: mnogoGrumna")),
		},
		// Task status not finished
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					status := model.Task_STATUS_STARTED
					taskType := "notMyType"
					return model.Task{Status: &status, TaskType: &taskType}, nil
				},
				errorMessage: "Cluster creation failed",
				finishCallback: func(task model.Task) {
					assert.Equal(t, model.Task_STATUS_STARTED, *task.Status)
				},
			},
			want: retry.RetryableError(fmt.Errorf("expected task type: notMyType to be finished STARTED")),
		},
		// Task status invalid
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					status := ""
					return model.Task{Status: &status}, nil
				},
				errorMessage: "Cluster creation failed",
				finishCallback: func(task model.Task) {
					assert.Equal(t, "", *task.Status)
					finishCallbackHasBeenCalled = true
				},
			},
			want: retry.NonRetryableError(fmt.Errorf("task status was empty. Some API error occurred")),
		},
		// Task status finished
		{
			input: inputStruct{
				connectorWrapper: AuthenticatorStub{},
				taskSupplier: func() (model.Task, error) {
					status := model.Task_STATUS_FINISHED
					return model.Task{Status: &status}, nil
				},
				errorMessage: "Cluster creation failed",
				finishCallback: func(task model.Task) {
					assert.Equal(t, model.Task_STATUS_FINISHED, *task.Status)
					finishCallbackHasBeenCalled = true
				},
			},
			want: nil,
		},
	}
	for _, testCase := range tests {
		got := RetryTaskUntilFinished(testCase.input.connectorWrapper,
			testCase.input.taskSupplier,
			testCase.input.errorMessage,
			testCase.input.finishCallback)
		assert.Equal(t, got, testCase.want)
	}
	assert.Equal(t, finishCallbackHasBeenCalled, true)
}
