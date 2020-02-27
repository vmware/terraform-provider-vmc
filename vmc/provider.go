/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/golangsdk/vsphere-automation-sdk-go/runtime/protocol/client"
)

type ConnectorWrapper struct {
	client.Connector
	RefreshToken string
	VmcURL       string
	CspURL       string
}

func (c *ConnectorWrapper) authenticate() error {
	var err error
	httpClient := http.Client{}
	c.Connector, err = NewVmcConnectorByRefreshToken(c.RefreshToken, c.VmcURL, c.CspURL, httpClient)
	if err != nil {
		return err
	}
	return nil
}

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("API_TOKEN", nil),
			},
			"vmc_url": {
				Type:     schema.TypeString,
				Optional: true,
				// Default:  "https://vmc.vmware.com",
				Default: "https://nsx-3-219-173-255.rp.vmwarevmc.com/vmc/reverse-proxy/api/orgs/54937bce-8119-4fae-84f5-e5e066ee90e6/sddcs/9f5f1fc4-b6a7-471b-b614-45f416fcaca7/cloud-service/api/v1/infra/public-ips/1ad34",
			},
			"csp_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console.cloud.vmware.com",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc":     resourceSddc(),
			"vmc_publicip": resourcePublicIp(),
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
	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)
	httpClient := http.Client{}
	connector, err := NewVmcConnectorByRefreshToken(refreshToken, vmcURL, cspURL, httpClient)
	if err != nil {
		return nil, fmt.Errorf("Error creating connector : %v ", err)
	}

	return &ConnectorWrapper{connector, refreshToken, vmcURL, cspURL}, nil
}
