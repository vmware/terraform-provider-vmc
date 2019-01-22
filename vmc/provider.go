package vmc

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/vmc-go-sdk/vmc"
)

// Provider for VMware VMC Console APIs. Returns terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"refresh_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"csp_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://console-stg.cloud.vmware.com",
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
	cspURL := d.Get("csp_url").(string)

	apiClient, err := vmc.NewVmcClient(refreshToken, cspURL)

	return apiClient, err
}
