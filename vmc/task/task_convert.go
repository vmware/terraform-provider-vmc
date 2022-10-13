/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package task

import (
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	autoscalerapi "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api"
	draas "github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
)

// The following functions create API clients and poll for tasks with the provided ID.
// They are very difficult to unit-test as the API clients don't implement any interface, so
// passing in (stubbed) API client as a parameter of the functions is not possible.
// This concern has been raised with the owners of the VMC SDKs.

// GetTask returns a model.Task with specified ID
func GetTask(connectorWrapper *connector.ConnectorWrapper, taskId string) (model.Task, error) {
	tasksClient := orgs.NewTasksClient(connectorWrapper)
	return tasksClient.Get(connectorWrapper.OrgID, taskId)
}

// GetAutoscalerTask polls autoscalerapi for task with specified ID and converts it to model.Task
func GetAutoscalerTask(connectorWrapper *connector.ConnectorWrapper, taskId string) (model.Task, error) {
	tasksClient := autoscalerapi.NewAutoscalerClient(connectorWrapper)
	autoscalerTask, err := tasksClient.Get(connectorWrapper.OrgID, taskId)
	// Commented out fields do not exist in the autoscalerapi task
	return model.Task{
		Updated:               autoscalerTask.Updated,
		UserId:                autoscalerTask.UserId,
		UpdatedByUserId:       autoscalerTask.UpdatedByUserId,
		Created:               autoscalerTask.Created,
		UserName:              autoscalerTask.UserName,
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
func GetDraasTask(connectorWrapper *connector.ConnectorWrapper, taskId string) (model.Task, error) {
	tasksClient := draas.NewTaskClient(connectorWrapper)
	draasTask, err := tasksClient.Get(connectorWrapper.OrgID, taskId)
	// Commented out fields do not exist in draas API Task
	return model.Task{
		Updated:               draasTask.Updated,
		UserId:                draasTask.UserId,
		UpdatedByUserId:       draasTask.UpdatedByUserId,
		Created:               draasTask.Created,
		UserName:              draasTask.UserName,
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
		//OrgId:
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
