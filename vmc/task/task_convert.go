// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package task

import (
	autoscalerapi "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
)

// The following functions create API clients and poll for tasks with the provided ID.
// They are very difficult to unit-test as the API clients don't implement any interface, so
// passing in (stubbed) API client as a parameter of the functions is not possible.
// This concern has been raised with the owners of the VMC SDKs.

// GetTask returns a model.Task with specified ID
func GetTask(connectorWrapper *connector.Wrapper, taskID string) (model.Task, error) {
	tasksClient := orgs.NewTasksClient(connectorWrapper)
	return tasksClient.Get(connectorWrapper.OrgID, taskID)
}

// GetV2Task returns an adapted model.Task with specified ID
func GetV2Task(connectorWrapper *connector.Wrapper, taskID string) (model.Task, error) {
	tasksV2Client := NewV2ClientImpl(*connectorWrapper)
	err := tasksV2Client.Authenticate()
	if err != nil {
		return model.Task{}, err
	}
	taskV2, err := tasksV2Client.GetTask(taskID)
	if err != nil {
		return model.Task{}, err
	}
	taskStatus := taskV2.TaskState.Name
	// convert v2 "finished" status to v1 "finished" status
	if taskStatus == "COMPLETED" {
		taskStatus = model.Task_STATUS_FINISHED
	}
	return model.Task{
		Id:           taskV2.ID,
		TaskType:     &taskV2.TaskType,
		Status:       &taskStatus,
		ErrorMessage: &taskV2.ErrorMessage,
	}, nil
}

// GetAutoscalerTask polls autoscalerapi for task with specified ID and converts it to model.Task
func GetAutoscalerTask(connectorWrapper *connector.Wrapper, taskID string) (model.Task, error) {
	tasksClient := autoscalerapi.NewAutoscalerClient(connectorWrapper)
	autoscalerTask, err := tasksClient.Get(connectorWrapper.OrgID, taskID)
	// Commented out fields do not exist in the autoscalerapi task
	return model.Task{
		Updated:               autoscalerTask.Updated,
		UserId:                autoscalerTask.UserId,
		UpdatedByUserId:       &autoscalerTask.UpdatedByUserId,
		Created:               autoscalerTask.Created,
		UserName:              &autoscalerTask.UserName,
		Id:                    autoscalerTask.Id,
		Status:                autoscalerTask.Status,
		LocalizedErrorMessage: autoscalerTask.ErrorMessage,
		ResourceId:            autoscalerTask.ResourceId,
		TaskVersion:           autoscalerTask.TaskVersion,
		//CorrelationId,
		//StartResourceEntityVersion
		//CustomerErrorMessage
		SubStatus: autoscalerTask.SubStatus,
		TaskType:  autoscalerTask.TaskType,
		StartTime: autoscalerTask.StartTime,
		//TaskProgressPhases
		ErrorMessage: autoscalerTask.ErrorMessage,
		OrgId:        autoscalerTask.OrgId,
		//EndResourceEntityVersion
		//ServiceErrors
		//OrgType
		EstimatedRemainingMinutes: autoscalerTask.EstimatedRemainingMinutes,
		Params:                    autoscalerTask.Params,
		ProgressPercent:           autoscalerTask.ProgressPercent,
		PhaseInProgress:           autoscalerTask.PhaseInProgress,
		ResourceType:              autoscalerTask.ResourceType,
		EndTime:                   autoscalerTask.EndTime,
	}, err
}

// GetDraasTask polls draas API for task with specified ID and converts it to model.Task
func GetDraasTask(connectorWrapper *connector.Wrapper, taskID string) (model.Task, error) {
	tasksClient := draas.NewTaskClient(connectorWrapper)
	draasTask, err := tasksClient.Get(connectorWrapper.OrgID, taskID)
	// Commented out fields do not exist in draas API Task
	return model.Task{
		Updated:               draasTask.Updated,
		UserId:                draasTask.UserId,
		UpdatedByUserId:       &draasTask.UpdatedByUserId,
		Created:               draasTask.Created,
		UserName:              &draasTask.UserName,
		Id:                    draasTask.Id,
		Status:                draasTask.Status,
		LocalizedErrorMessage: draasTask.ErrorMessage,
		ResourceId:            draasTask.ResourceId,
		TaskVersion:           draasTask.TaskVersion,
		//CorrelationId,
		//StartResourceEntityVersion
		//CustomerErrorMessage
		SubStatus: draasTask.SubStatus,
		TaskType:  draasTask.TaskType,
		StartTime: draasTask.StartTime,
		//TaskProgressPhases
		ErrorMessage: draasTask.ErrorMessage,
		//OrgID:
		//EndResourceEntityVersion
		//ServiceErrors
		//OrgType
		EstimatedRemainingMinutes: draasTask.EstimatedRemainingMinutes,
		Params:                    draasTask.Params,
		ProgressPercent:           draasTask.ProgressPercent,
		//PhaseInProgress
		ResourceType: draasTask.ResourceType,
		EndTime:      draasTask.EndTime,
	}, err
}
