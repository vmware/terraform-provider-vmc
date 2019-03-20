package vmc

import (
	"context"
	"fmt"

	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

func dataSourceVmcOrg() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcOrgRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Unique ID of this resource",
				Required:    true,
			},
			"display_name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The display name of this resource",
				Computed:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The Name of this resource",
				Computed:    true,
			},
		},
	}
}

func dataSourceVmcOrgRead(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	objID := d.Get("id").(string)
	var obj vmc.Organization
	obj, resp, err := vmcClient.OrgsApi.OrgsOrgGet(context.Background(), objID)
	if err != nil {
		return fmt.Errorf("Error while reading ns group %s: %v", objID, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("NS group %s was not found", objID)
	}
	d.SetId(obj.Id)
	d.Set("display_name", obj.DisplayName)
	d.Set("name", obj.Name)

	return nil
}
