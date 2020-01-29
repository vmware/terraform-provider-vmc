/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/account_link"
)

func dataSourceVmcConnectedAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcConnectedAccountsRead,

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Required:    true,
			},
			"provider_type": {
				Type:        schema.TypeString,
				Description: "The cloud provider of the SDDC (AWS or ZeroCloud).",
				Optional:    true,
				Default:     "AWS",
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "AWS account number.",
				Optional:    true,
			},
			"ids": {
				Type:        schema.TypeList,
				Description: "The corresponding connected (customer) account UUID this connection is attached to.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVmcConnectedAccountsRead(d *schema.ResourceData, m interface{}) error {

	orgID := d.Get("org_id").(string)
	providerType := d.Get("provider_type").(string)
	accountNumber := d.Get("account_number").(string)

	if orgID == "" {
		return fmt.Errorf("org ID is a required parameter and cannot be empty")
	}

	connector := (m.(*ConnectorWrapper)).Connector
	defaultConnectedAccountsClient := account_link.NewDefaultConnectedAccountsClient(connector)
	accounts, err := defaultConnectedAccountsClient.Get(orgID, &providerType)

	ids := []string{}
	for _, account := range accounts {
		if(accountNumber != "") {
			if(*account.AccountNumber == accountNumber) {
				ids = append(ids, account.Id)
			}
		} else {
			ids = append(ids, account.Id)
		}
	}

	if err != nil {
		return fmt.Errorf("Error while reading accounts from org %q: %v", orgID, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", orgID, providerType))
	d.Set("ids", ids)
	return nil
}
