// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/account_link"
)

func dataSourceVmcConnectedAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcConnectedAccountsRead,

		Schema: map[string]*schema.Schema{
			"provider_type": {
				Type:        schema.TypeString,
				Description: "The cloud provider of the SDDC (AWS or ZeroCloud).",
				Optional:    true,
				Default:     "AWS",
			},
			"account_number": {
				Type:        schema.TypeString,
				Description: "AWS account number.",
				Required:    true,
			},
			"id": {
				Type:        schema.TypeString,
				Description: "The corresponding connected (customer) account UUID this connection is attached to.",
				Computed:    true,
			},
		},
	}
}

func dataSourceVmcConnectedAccountsRead(d *schema.ResourceData, m interface{}) error {
	orgID := (m.(*connector.Wrapper)).OrgID
	providerType := d.Get("provider_type").(string)
	accountNumber := d.Get("account_number").(string)

	connectorWrapper := (m.(*connector.Wrapper)).Connector
	defaultConnectedAccountsClient := account_link.NewConnectedAccountsClient(connectorWrapper)
	accounts, err := defaultConnectedAccountsClient.Get(orgID, &providerType)

	if accountNumber == "" {
		return fmt.Errorf("account number is a required parameter and cannot be empty")
	}
	id := ""
	for _, account := range accounts {
		if *account.AccountNumber == accountNumber {
			id = account.Id
		}
	}

	if err != nil {
		return HandleDataSourceReadError("Connected Accounts", err)
	}

	if id == "" {
		return fmt.Errorf("no connected account found with the account number : %q ", accountNumber)
	}

	d.SetId(id)
	return nil
}
