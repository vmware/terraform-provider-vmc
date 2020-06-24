/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

const (
	// DefaultVMCServer defines the default VMC server.
	DefaultVMCServer string = "https://vmc.vmware.com"

	// DefaultCSPUrl defines the default URL for CSP.
	DefaultCSPUrl string = "https://console.cloud.vmware.com"

	// CSPRefreshUrlSuffix defines the CSP Refresh API endpoint.
	CSPRefreshUrlSuffix string = "/csp/gateway/am/api/auth/api-tokens/authorize"

	// sksNSXTManager to be stripped from nsxt reverse proxy url for public IP resource
	SksNSXTManager string = "/sks-nsxt-manager"

	// ESX Host instance types supported for SDDC creation.
	HostInstancetypeI3   string = "I3_METAL"
	HostInstancetypeR5   string = "R5_METAL"
	HostInstancetypeI3EN string = "I3EN_METAL"

	// Availability Zones
	SingleAvailabilityZone string = "SingleAZ"
	MultiAvailabilityZone  string = "MultiAZ"
	MinMultiAZHosts        int    = 6

	ClusterIdFieldName = "clusterId"

	// Env variables used in acceptance tests
	APIToken            string = "API_TOKEN"
	OrgID               string = "ORG_ID"
	OrgDisplayName      string = "ORG_DISPLAY_NAME"
	TestSDDCId          string = "TEST_SDDC_ID"
	AWSAccountNumber    string = "AWS_ACCOUNT_NUMBER"
	NSXTReverseProxyUrl string = "NSXT_REVERSE_PROXY_URL"
)
