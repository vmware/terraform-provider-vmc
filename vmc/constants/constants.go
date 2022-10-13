/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package constants

const (
	// DefaultVmcUrl defines the default VMC server url.
	DefaultVmcUrl string = "https://vmc.vmware.com"

	// DefaultCspUrl defines the default URL for CSP.
	DefaultCspUrl string = "https://console.cloud.vmware.com"

	// CspRefreshUrlSuffix defines the CSP Refresh API endpoint.
	CspRefreshUrlSuffix string = "/csp/gateway/am/api/auth/api-tokens/authorize"

	// sksNSXTManager to be stripped from nsxt reverse proxy url for public IP resource
	SksNsxtManager string = "/sks-nsxt-manager"

	// ESX Host instance types supported for SDDC creation.
	HostInstancetypeI3   string = "I3_METAL"
	HostInstancetypeR5   string = "R5_METAL"
	HostInstancetypeI3EN string = "I3EN_METAL"
	HostInstancetypeI4I  string = "I4I_METAL"

	// Availability Zones
	SingleAvailabilityZone string = "SingleAZ"
	MultiAvailabilityZone  string = "MultiAZ"
	MinMultiAZHosts        int    = 6

	// SDDC Size
	MediumSddcSize        = "medium"
	CapitalMediumSddcSize = "MEDIUM"
	LargeSddcSize         = "large"
	CapitalLargeSddcSize  = "LARGE"

	ClusterIdFieldName = "clusterId"
	SrmPrefix          = "srm-"
	SddcSuffix         = ".sddc-"

	// EDRS Policy types
	CostPolicyType           = "cost"
	PerformancePolicyType    = "performance"
	StorageScaleUpPolicyType = "storage-scaleup"
	RapidScaleUpPolicyType   = "rapid-scaleup"

	// Microsoft licensing config actions
	LicenseConfigEnabled         = "enabled"
	LicenseConfigDisabled        = "disabled"
	CapitalLicenseConfigEnabled  = "ENABLED"
	CapitalLicenseConfigDisabled = "DISABLED"

	// SDDC Type
	OneNodeSddcType = "1NODE"

	// Provider Types
	AwsProviderType       = "AWS"
	ZeroCloudProviderType = "ZEROCLOUD"

	// Intranet Uplink MTU Range
	MinIntranetMtuLink = 1500
	MaxIntranetMtuLink = 8900

	// Range for number of hosts
	MinHosts = 2
	MaxHosts = 16

	// Env variables used in acceptance tests
	VmcUrl         string = "VMC_URL"
	CspUrl         string = "CSP_URL"
	ApiToken       string = "API_TOKEN"
	OrgID          string = "ORG_ID"
	OrgDisplayName string = "ORG_DISPLAY_NAME"
	// TestSddcId ID of an existing SDDC used for sddc data source, site recovery and srm node tests
	TestSddcId string = "TEST_SDDC_ID"
	// TestSddcName Name of an existing SDDC used for sddc data source tests
	TestSddcName        string = "TEST_SDDC_NAME"
	AwsAccountNumber    string = "AWS_ACCOUNT_NUMBER"
	NsxtReverseProxyUrl string = "NSXT_REVERSE_PROXY_URL"
)
