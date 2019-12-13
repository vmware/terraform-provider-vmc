/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc"
)

func dataSourceVmcOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcOrgRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Required:    true,
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
	orgID := d.Get("id").(string)
	if orgID == "" {
		return fmt.Errorf("org ID is a required parameter and cannot be empty")
	}

	connector := (m.(*ConnectorWrapper)).Connector
	orgClient := vmc.NewDefaultOrgsClient(connector)
	org, err := orgClient.Get(orgID)

	if err != nil {
		return fmt.Errorf("Error while reading org information for %s: %v", orgID, err)
	}
	d.SetId(org.Id)
	d.Set("display_name", org.DisplayName)
	d.Set("name", org.Name)

	return nil
}
