package vmc

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"vmc_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://vmc.vmware.com/vmc/api",
			},
			"csp_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console.cloud.vmware.com",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmc_sddc": resourceSddc(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"vmc_org":                dataSourceVmcOrg(),
			"vmc_connected_accounts": dataSourceVmcConnectedAccounts(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	refreshToken := d.Get("refresh_token").(string)
	vmcURL := d.Get("vmc_url").(string)
	cspURL := d.Get("csp_url").(string)

	println(vmcURL)
	println(cspURL)

	apiClient, err := vmc.NewVmcClient(refreshToken, vmcURL, cspURL)

	return apiClient, err
}
