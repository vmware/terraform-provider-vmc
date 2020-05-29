/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
)

type ConnectorWrapper struct {
	client.Connector
	RefreshToken string
	ClientID     string
	ClientSecret string
	OrgID        string
	VmcURL       string
	CspURL       string
}

func (c *ConnectorWrapper) authenticate() error {
	var err error
	httpClient := http.Client{}
	if len(c.RefreshToken) > 0 {
		c.Connector, err = NewClientConnectorByRefreshToken(c.RefreshToken, c.VmcURL, c.CspURL, httpClient)
		if err != nil {
			return err
		}
	} else {
		c.Connector, err = NewClientConnectorByClientID(c.ClientID, c.ClientSecret, c.VmcURL, c.CspURL, httpClient)
		if err != nil {
			return err
		}
	}

	return nil
}

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ApiToken, nil),
			},
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ClientID, nil),
			},
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(ClientSecret, nil),
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(OrgID, nil),
			},
			"vmc_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://vmc.vmware.com",
			},
			"csp_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console.cloud.vmware.com/csp/gateway",
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
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	refreshToken := d.Get("refresh_token").(string)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	if len(refreshToken) <= 0 && len(clientID) <= 0 && len(clientSecret) <= 0 {
		return nil, fmt.Errorf("must provide value for refresh_token or client_id and client_secret")
	}

	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)
	orgID := d.Get("org_id").(string)
	httpClient := http.Client{}

	if len(refreshToken) > 0 {
		// set refresh token to env variable so that it can be used by other connectors
		os.Setenv(ApiToken, refreshToken)
		connector, err := NewClientConnectorByRefreshToken(refreshToken, vmcURL, cspURL, httpClient)
		if err != nil {
			return nil, HandleCreateError("Client connector using refresh token", err)
		}

		return &ConnectorWrapper{connector, refreshToken, clientID, clientSecret, orgID, vmcURL, cspURL}, nil
	} else {
		// set client ID and client secret to env variable so that it can be used by other connectors
		os.Setenv(ClientID, clientID)
		os.Setenv(ClientSecret, clientSecret)

		connector, err := NewClientConnectorByClientID(clientID, clientSecret, vmcURL, cspURL, httpClient)
		if err != nil {
			return nil, HandleCreateError("Client connector using client ID and client secret", err)
		}

		return &ConnectorWrapper{connector, refreshToken, clientID, clientSecret, orgID, vmcURL, cspURL}, nil
	}
}
