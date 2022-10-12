/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	autoscalerapi "github.com/vmware/vsphere-automation-sdk-go/services/vmc/autoscaler/api"
	draas "github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas"
	"net/http"
	"net/url"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
)

var storageCapacityMap = map[string]int64{
	"15TB": 15003,
	"20TB": 20004,
	"25TB": 25005,
	"30TB": 30006,
	"35TB": 35007,
}

func GetSDDC(connector client.Connector, orgID string, sddcID string) (model.Sddc, error) {
	sddcClient := orgs.NewSddcsClient(connector)
	sddc, err := sddcClient.Get(orgID, sddcID)
	return sddc, err
}

func ConvertStorageCapacitytoInt(s string) int64 {
	storageCapacity := storageCapacityMap[s]
	return storageCapacity
}

// ConvertDeployType Mapping for deployment_type field
// During refresh/import state, return value of VMC API should be converted to uppercamel case in terraform
// to maintain consistency
func ConvertDeployType(s string) string {
	if s == "SINGLE_AZ" {
		return SingleAvailabilityZone
	} else if s == "MULTI_AZ" {
		return MultiAvailabilityZone
	} else {
		return ""
	}
}

func IsValidUUID(u string) error {
	_, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	return nil
}

func IsValidURL(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	return nil
}
func expandMsftLicenseConfig(config []interface{}) *model.MsftLicensingConfig {
	if len(config) == 0 {
		return nil
	}
	var licenseConfig model.MsftLicensingConfig
	licenseConfigMap := config[0].(map[string]interface{})
	mssqlLicensing := strings.ToUpper(licenseConfigMap["mssql_licensing"].(string))
	windowsLicensing := strings.ToUpper(licenseConfigMap["windows_licensing"].(string))
	licenseConfig = model.MsftLicensingConfig{MssqlLicensing: &mssqlLicensing, WindowsLicensing: &windowsLicensing}
	return &licenseConfig
}

func getNSXTReverseProxyURLConnector(nsxtReverseProxyUrl string) (client.Connector, error) {
	apiToken := os.Getenv(APIToken)
	if len(nsxtReverseProxyUrl) == 0 {
		return nil, fmt.Errorf("NSX reverse proxy url is required for public IP resource creation")
	}
	nsxtReverseProxyUrl = strings.Replace(nsxtReverseProxyUrl, SksNSXTManager, "", -1)
	httpClient := http.Client{}
	cspUrl := os.Getenv(CSPUrl)
	connector, err := NewClientConnectorByRefreshToken(apiToken, nsxtReverseProxyUrl, cspUrl, httpClient)
	if err != nil {
		return nil, HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	return connector, nil
}

// getHostCountCluster tries to find the amount of hosts on a Cluster in
// the ResourceConfig of the provided SDDC. If there is no ResourceConfig/Cluster 0 is returned.
// A Cluster is distinguished by its id
func getHostCountCluster(sddc *model.Sddc, clusterId string) int {
	if sddc != nil && sddc.ResourceConfig != nil && sddc.ResourceConfig.Clusters != nil {
		for _, cluster := range sddc.ResourceConfig.Clusters {
			if cluster.ClusterId == clusterId {
				return len(cluster.EsxHostList)
			}
		}
	}
	return 0
}

// toHostInstanceType converts from the Schema format of the host_instance_type to
// the possible string values defined in the VMC SDK
func toHostInstanceType(userPassedHostInstanceType string) (string, error) {
	switch userPassedHostInstanceType {
	case HostInstancetypeI3:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL, nil
	case HostInstancetypeI3EN:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL, nil
	case HostInstancetypeI4I:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL, nil
	case HostInstancetypeR5:
		return model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL, nil
	default:
		return "", fmt.Errorf("unknown host instance type: %s", userPassedHostInstanceType)
	}
}

// getTask returns a model.Task with specified ID
func getTask(connectorWrapper *ConnectorWrapper, taskId string) (model.Task, error) {
	tasksClient := orgs.NewTasksClient(connectorWrapper)
	return tasksClient.Get(connectorWrapper.OrgID, taskId)
}

// getAutoscalerTask polls autoscalerapi for task with specified ID and converts it to model.Task
func getAutoscalerTask(connectorWrapper *ConnectorWrapper, taskId string) (model.Task, error) {
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

// getDraasTask polls draas API for task with specified ID and converts it to model.Task
func getDraasTask(connectorWrapper *ConnectorWrapper, taskId string) (model.Task, error) {
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
