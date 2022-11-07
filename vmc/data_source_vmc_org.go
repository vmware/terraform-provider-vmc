/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc"
)

func dataSourceVmcOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcOrgRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Computed:    true,
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Computed:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "The Name of this resource",
				Computed:    true,
			},
		},
	}
}

func dataSourceVmcOrgRead(d *schema.ResourceData, m interface{}) error {
	orgID := (m.(*connector.Wrapper)).OrgID
	connectorWrapper := (m.(*connector.Wrapper)).Connector
	orgClient := vmc.NewOrgsClient(connectorWrapper)
	org, err := orgClient.Get(orgID)
	if err != nil {
		return HandleDataSourceReadError(d, "VMC Organization", err)
	}
	d.SetId(orgID)
	d.Set("display_name", org.DisplayName)
	d.Set("name", org.Name)

	return nil
}
