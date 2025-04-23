// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
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
	}
	if s == "MULTI_AZ" {
		return constants.MultiAvailabilityZone
	}
	return ""
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
	academicLicense := licenseConfigMap["academic_license"].(bool)
	licenseConfig = model.MsftLicensingConfig{MssqlLicensing: &mssqlLicensing, WindowsLicensing: &windowsLicensing, AcademicLicense: &academicLicense}
	return &licenseConfig
}

func getNsxtReverseProxyURLConnector(nsxtReverseProxyURL string, wrapper *connector.Wrapper) (client.Connector, error) {
	if len(nsxtReverseProxyURL) == 0 {
		return nil, fmt.Errorf("NSX reverse proxy url is required for public IP resource creation")
	}
	if wrapper == nil {
		return nil, fmt.Errorf("nil connector.Wrapper provided")
	}
	nsxtReverseProxyURL = strings.Replace(nsxtReverseProxyURL, constants.SksNSXTManager, "", -1)
	copyWrapper := connector.CopyWrapper(*wrapper)
	// The wrapper uses the VmcURL as service URL, so setting it to the NSX URL will
	// force authentication against the NSX instance
	copyWrapper.VmcURL = nsxtReverseProxyURL
	err := copyWrapper.Authenticate()
	if err != nil {
		return nil, err
	}
	return copyWrapper.Connector, nil
}

// getHostCountCluster tries to find the amount of hosts on a Cluster in
// the ResourceConfig of the provided SDDC. If there is no ResourceConfig/Cluster 0 is returned.
// A Cluster is distinguished by its id
func getHostCountCluster(sddc *model.Sddc, clusterID string) int {
	if sddc != nil && sddc.ResourceConfig != nil && sddc.ResourceConfig.Clusters != nil {
		for _, cluster := range sddc.ResourceConfig.Clusters {
			if cluster.ClusterId == clusterID {
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
	case constants.HostInstancetypeC6I:
		return model.SddcConfig_HOST_INSTANCE_TYPE_C6I_METAL, nil
	case constants.HostInstancetypeM7i24xl:
		return model.SddcConfig_HOST_INSTANCE_TYPE_M7I_METAL_24XL, nil
	case constants.HostInstancetypeM7i48xl:
		return model.SddcConfig_HOST_INSTANCE_TYPE_M7I_METAL_48XL, nil
	default:
		return "", fmt.Errorf("unknown host instance type: %s", userPassedHostInstanceType)
	}
}
