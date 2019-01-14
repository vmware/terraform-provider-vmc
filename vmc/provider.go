package vmc

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/het/vmc-go-sdk/vmc"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc": resourceSddc(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vmc_org": dataSourceVmcOrg(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	refreshToken := d.Get("refresh_token").(string)

	apiClient := vmc.NewVMCClient(refreshToken)

	return apiClient, nil
}
