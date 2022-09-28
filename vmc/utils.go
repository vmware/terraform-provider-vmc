/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
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
		return "SingleAZ"
	} else if s == "MULTI_AZ" {
		return "MultiAZ"
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

// getHostCountOnPrimaryCluster tries to find the amount of hosts on the primary Cluster in
// the ResourceConfig of the provided SDDC. If there is no ResourceConfig/Cluster 0 is returned.
// The primary Cluster is distinguished by its id
func getHostCountOnPrimaryCluster(sddc *model.Sddc, primaryClusterId string) int {
	primaryClusterHostCount := 0
	if sddc != nil && sddc.ResourceConfig != nil && sddc.ResourceConfig.Clusters != nil {
		for _, cluster := range sddc.ResourceConfig.Clusters {
			if cluster.ClusterId == primaryClusterId {
				primaryClusterHostCount += len(cluster.EsxHostList)
			}
		}
	}
	return primaryClusterHostCount
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
