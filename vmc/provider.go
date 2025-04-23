// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(constants.APIToken, nil),
				ConflictsWith: []string{"client_id", "client_secret"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(constants.ClientID, nil),
				ConflictsWith: []string{"refresh_token"},
				RequiredWith:  []string{"client_secret"},
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				DefaultFunc:   schema.EnvDefaultFunc(constants.ClientSecret, nil),
				ConflictsWith: []string{"refresh_token"},
				RequiredWith:  []string{"client_id"},
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc(constants.OrgID, nil),
			},
			"vmc_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(constants.VmcURL, constants.DefaultVmcURL),
			},
			"csp_url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc(constants.CspURL, constants.DefaultCspURL),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc":          resourceSddc(),
			"vmc_public_ip":     resourcePublicIP(),
			"vmc_site_recovery": resourceSiteRecovery(),
			"vmc_srm_node":      resourceSrmNode(),
			"vmc_cluster":       resourceCluster(),
			"vmc_sddc_group":    resourceSddcGroup(),
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
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	if len(refreshToken) == 0 && len(clientID) == 0 && len(clientSecret) == 0 {
		return nil, fmt.Errorf("must provide value for refresh_token or client_id and client_secret")
	}
	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)
	orgID := d.Get("org_id").(string)
	connectorWrapper := connector.Wrapper{
		RefreshToken: refreshToken,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		OrgID:        orgID,
		VmcURL:       vmcURL,
		CspURL:       cspURL,
	}
	err := connectorWrapper.Authenticate()
	if err != nil {
		return nil, HandleCreateError("Client connector", err)
	}

	return &connectorWrapper, err
}
