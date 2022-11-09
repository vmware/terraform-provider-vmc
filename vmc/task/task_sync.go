/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package task

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"log"
	"sync"
)

// KeyedMutex Mutex that operates multiple locks, based  on a string key.
type KeyedMutex struct {
	mutexes sync.Map // Zero value thread-safe map is empty and ready for use
}

// Lock Locks on a key, allowing multiple threads to operate on separate keys. Returns
// a function, that clients should use to unlock the locks they've obtained.
func (keyedMutex *KeyedMutex) Lock(key string) func() {
	value, _ := keyedMutex.mutexes.LoadOrStore(key, &sync.Mutex{})
	mutex := value.(*sync.Mutex)
	mutex.Lock()

	// Encapsulate the access to the underlying mutex, but allow clients to unlock the
	// mutex they've just locked on.
	return func() {
		mutex.Unlock()
	}
}

// global retry counter. The added complexity of a counter per task is unlikely to pay off.
var serviceUnavailableRetries = 0

// max amount of retries for "service unavailable" errors before giving up
var maxServiceUnavailableRetries = 20

// RetryTaskUntilFinished function that will poll (using provided task supplier) for a
// task state until a non-recoverable error is encountered, like task failure or
// authentication error or until the task finishes. An option to execute a callback after task
// finish (either successfully or not) is provided.
func RetryTaskUntilFinished(authenticator connector.Authenticator,
	taskSupplier func() (model.Task, error),
	errorMessage string,
	finishCallback func(task model.Task)) *resource.RetryError {
	task, err := taskSupplier()
	if err != nil {
		// Try to reauthenticate (if access token expired)
		if err.Error() == (errors.Unauthenticated{}.Error()) {
			log.Printf("Authentication error : %v", errors.Unauthenticated{}.Error())
			err = authenticator.Authenticate()
			if err != nil {
				if finishCallback != nil {
					finishCallback(task)
				}
				return resource.NonRetryableError(fmt.Errorf("authentication error from Cloud Service Provider : %v", err))
			}
			return resource.RetryableError(fmt.Errorf("task still in progress"))
		}
		// Best-effort resiliency in case of difficulties the VMC service may experience,
		// during long-running tasks
		if err.Error() == (errors.ServiceUnavailable{}.Error()) {
			serviceUnavailableRetries++
			if serviceUnavailableRetries <= maxServiceUnavailableRetries {
				return resource.RetryableError(fmt.Errorf(
					"VMC backend is experiencing difficulties, retry %d from %d to polling the SDDC Create Task",
					serviceUnavailableRetries, maxServiceUnavailableRetries))
			}
			if finishCallback != nil {
				finishCallback(task)
			}
			return resource.NonRetryableError(fmt.Errorf("max ServiceUnavailable retries (20) reached to create SDDC"))
		}
		if finishCallback != nil {
			finishCallback(task)
		}
		return resource.NonRetryableError(fmt.Errorf(errorMessage+": %v", err))

	}
	// If code reached this point it is safe to assume "service unavailable" window passed,
	// so reset the global counter
	if serviceUnavailableRetries > 0 {
		serviceUnavailableRetries = 0
	}
	if *task.Status == "" {
		if finishCallback != nil {
			finishCallback(task)
		}
		return resource.NonRetryableError(fmt.Errorf("task status was empty. Some API error occurred"))
	} else if *task.Status == model.Task_STATUS_FAILED {
		if finishCallback != nil {
			finishCallback(task)
		}
		return resource.NonRetryableError(fmt.Errorf("task failed: "+errorMessage+": %s", *task.ErrorMessage))
	} else if *task.Status != model.Task_STATUS_FINISHED {
		return resource.RetryableError(fmt.Errorf("expected task type: %s to be finished %s", *task.TaskType, *task.Status))
	}
	if finishCallback != nil {
		finishCallback(task)
	}
	return nil
}
