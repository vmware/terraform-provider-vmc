package vmc

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"

	"gitlab.eng.vmware.com/het/vmc-go-sdk/vmc"

	"net/http"
)

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of this resource",
			},
			"sddc_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"num_host": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"provider_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ZEROCLOUD",
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "US_WEST_2",
			},
			"created": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	orgID := d.Get("org_id").(string)
	sddcName := d.Get("sddc_name").(string)
	numHost := d.Get("num_host").(int)
	providerType := d.Get("provider_type").(string)
	region := d.Get("region").(string)
	var awsSddcConfig = &vmc.AwsSddcConfig{
		Name:     sddcName,
		NumHosts: int32(numHost),
		Provider: providerType,
		Region:   region,
	}

	// Create a Sddc
	task, resp, err := vmcClient.SddcApi.OrgsOrgSddcsPost(context.TODO(), orgID, *awsSddcConfig)
	if err != nil {
		return fmt.Errorf("Error while creating sddc %s: %v", sddcName, err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	err = waitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}

	// Get Sddc detail
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.TODO(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Sddc %s was not found", sddcID)
	}

	d.SetId(sddc.Id)
	d.Set("name", sddc.Name)
	d.Set("created", sddc.Created)

	return nil
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.TODO(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return fmt.Errorf("Sddc %s was not found", sddcID)
	}

	d.Set("org_id", sddc.Id)
	d.Set("sddc_name", sddc.Name)
	d.Set("provider_type", sddc.Provider)
	d.Set("created", sddc.Created)

	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}
