package vmc

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

func dataSourceVmcConnectedAccounts() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcConnectedAccountsRead,

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Required:    true,
			},
			"provider_type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The cloud provider of the SDDC (AWS or ZeroCloud).",
				Optional:    true,
				Default:     "AWS",
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVmcConnectedAccountsRead(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	orgID := d.Get("org_id").(string)
	providerType := d.Get("provider_type").(string)
	// var obj vmc.Organization
	providerString := optional.NewString(providerType)
	accounts, _, err := vmcClient.AccountLinkingApi.OrgsOrgAccountLinkConnectedAccountsGet(
		context.Background(), orgID, &vmc.OrgsOrgAccountLinkConnectedAccountsGetOpts{Provider: providerString})

	ids := []string{}
	for _, account := range accounts {
		ids = append(ids, account.Id)
	}

	log.Printf("[DEBUG] Connected accounts are %v\n", accounts)

	if err != nil {
		return fmt.Errorf("Error while reading accounts from org %q: %v", orgID, err)
	}

	d.SetId(fmt.Sprintf("%s-%s", orgID, providerType))
	d.Set("ids", ids)
	return nil
}
