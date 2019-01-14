package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.eng.vmware.com/het/terraform-provider-vmc/vmc"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return vmc.Provider()
		},
	})
}
