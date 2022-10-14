/* Copyright 2020-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
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

func GetSddc(connector client.Connector, orgID string, sddcID string) (model.Sddc, error) {
	sddcClient := orgs.NewSddcsClient(connector)
	sddc, err := sddcClient.Get(orgID, sddcID)
	return sddc, err
}

func ConvertStorageCapacityToInt(s string) int64 {
	storageCapacity := storageCapacityMap[s]
	return storageCapacity
}

// ConvertDeployType Mapping for deployment_type field
// During refresh/import state, return value of VMC API should be converted to uppercamel case in terraform
// to maintain consistency
func ConvertDeployType(s string) string {
	if s == "SINGLE_AZ" {
		return constants.SingleAvailabilityZone
	} else if s == "MULTI_AZ" {
		return constants.MultiAvailabilityZone
	} else {
		return ""
	}
}

func IsValidUuid(u string) error {
	_, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	return nil
}

func IsValidUrl(s string) error {
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

func getNsxtReverseProxyURLConnector(nsxtReverseProxyUrl string) (client.Connector, error) {
	apiToken := os.Getenv(constants.ApiToken)
	if len(nsxtReverseProxyUrl) == 0 {
		return nil, fmt.Errorf("NSX reverse proxy url is required for public IP resource creation")
	}
	nsxtReverseProxyUrl = strings.Replace(nsxtReverseProxyUrl, constants.SksNsxtManager, "", -1)
	httpClient := http.Client{}
	cspUrl := os.Getenv(constants.CspUrl)
	apiConnector, err := connector.NewClientConnectorByRefreshToken(apiToken, nsxtReverseProxyUrl, cspUrl, httpClient)
	if err != nil {
		return nil, HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	return apiConnector, nil
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
	case constants.HostInstancetypeI3:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL, nil
	case constants.HostInstancetypeI3EN:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL, nil
	case constants.HostInstancetypeI4I:
		return model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL, nil
	case constants.HostInstancetypeR5:
		return model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL, nil
	default:
		return "", fmt.Errorf("unknown host instance type: %s", userPassedHostInstanceType)
	}
}
