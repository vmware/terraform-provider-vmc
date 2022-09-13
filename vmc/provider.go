/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

type ConnectorWrapper struct {
	client.Connector
	RefreshToken string
	OrgID        string
	VmcURL       string
	CspURL       string
}

func (c *ConnectorWrapper) authenticate() error {
	var err error
	httpClient := http.Client{}
	c.Connector, err = NewClientConnectorByRefreshToken(c.RefreshToken, c.VmcURL, c.CspURL, httpClient)
	if err != nil {
		return err
	}
	return nil
}

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(APIToken, nil),
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(OrgID, nil),
			},
			"vmc_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(VMCUrl, DefaultVMCUrl),
			},
			"csp_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(CSPUrl, DefaultCSPUrl),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc":          resourceSddc(),
			"vmc_public_ip":     resourcePublicIp(),
			"vmc_site_recovery": resourceSiteRecovery(),
			"vmc_srm_node":      resourceSRMNode(),
			"vmc_cluster":       resourceCluster(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vmc_org":                dataSourceVmcOrg(),
			"vmc_connected_accounts": dataSourceVmcConnectedAccounts(),
			"vmc_customer_subnets":   dataSourceVmcCustomerSubnets(),
			"vmc_sddc":               dataSourceVmcSddc(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	refreshToken := d.Get("refresh_token").(string)
	if len(refreshToken) <= 0 {
		return nil, fmt.Errorf("refresh token cannot be empty")
	}
	// set refresh token to env variable so that it can be used by other connectors
	os.Setenv(APIToken, refreshToken)
	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)
	os.Setenv(CSPUrl, cspURL)
	orgID := d.Get("org_id").(string)
	httpClient := http.Client{}
	connector, err := NewClientConnectorByRefreshToken(refreshToken, vmcURL, cspURL, httpClient)
	if err != nil {
		return nil, HandleCreateError("Client connector", err)
	}

	return &ConnectorWrapper{connector, refreshToken, orgID, vmcURL, cspURL}, nil
}
