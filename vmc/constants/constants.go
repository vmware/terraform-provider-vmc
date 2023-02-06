/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package constants

const (
	// DefaultVmcURL defines the default VMC server url.
	DefaultVmcURL string = "https://vmc.vmware.com"

	// DefaultCspURL defines the default URL for CSP.
	DefaultCspURL string = "https://console.cloud.vmware.com"

	// CspRefreshURLSuffix defines the CSP Refresh Token API endpoint.
	CspRefreshURLSuffix string = "/csp/gateway/am/api/auth/api-tokens/authorize"

	// CspOauthURLSuffix defines the CSP Oauth API endpoint.
	CspOauthURLSuffix string = "/csp/gateway/am/api/auth/token"

	// sksNSXTManager to be stripped from nsxt reverse proxy url for public IP resource
	SksNSXTManager string = "/sks-nsxt-manager"

	// ESX Host instance types supported for SDDC creation.
	HostInstancetypeI3   string = "I3_METAL"
	HostInstancetypeR5   string = "R5_METAL"
	HostInstancetypeI3EN string = "I3EN_METAL"
	HostInstancetypeI4I  string = "I4I_METAL"

	// Availability Zones
	SingleAvailabilityZone string = "SingleAZ"
	MultiAvailabilityZone  string = "MultiAZ"

	// SDDC Size
	MediumSddcSize        = "medium"
	CapitalMediumSddcSize = "MEDIUM"
	LargeSddcSize         = "large"
	CapitalLargeSddcSize  = "LARGE"

	ClusterIDFieldName = "clusterId"
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
	VmcURL         string = "VMC_URL"
	CspURL         string = "CSP_URL"
	APIToken       string = "API_TOKEN"
	ClientID       string = "CLIENT_ID"
	ClientSecret   string = "CLIENT_SECRET"
	OrgID          string = "ORG_ID"
	OrgDisplayName string = "ORG_DISPLAY_NAME"
	// TestSddcID ID of an existing SDDC used for sddc data source, site recovery and srm node tests
	TestSddcID string = "TEST_SDDC_ID"
	// TestSddcName Name of an existing SDDC used for sddc data source tests
	TestSddcName        string = "TEST_SDDC_NAME"
	AwsAccountNumber    string = "AWS_ACCOUNT_NUMBER"
	NsxtReverseProxyURL string = "NSXT_REVERSE_PROXY_URL"
)
