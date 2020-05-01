/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-vmc/vmc"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return vmc.Provider()
		},
	})
}
