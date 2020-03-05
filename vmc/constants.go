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

	// Env variables used in acceptance tests
	APIToken            string = "API_TOKEN"
	OrgID               string = "ORG_ID"
	OrgDisplayName      string = "ORG_DISPLAY_NAME"
	TestSDDCId          string = "TEST_SDDC_ID"
	AWSAccountNumber    string = "AWS_ACCOUNT_NUMBER"
	NsxtReverseProxyUrl string = "NSXT_REVERSE_PROXY_URL"
)
