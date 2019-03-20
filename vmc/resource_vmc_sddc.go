package vmc

import (
	"context"
	"fmt"
	"github.com/antihax/optional"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
	"net/http"
)

func resourceSddc() *schema.Resource {
	return &schema.Resource{
		Create: resourceSddcCreate,
		Read:   resourceSddcRead,
		Update: resourceSddcUpdate,
		Delete: resourceSddcDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of this resource",
			},
			"storage_capacity": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sddc_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"account_link_sddc_config": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_subnet_ids": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								// Optional: true,
							},
							Optional: true,
						},
						"connected_account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"vpc_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"num_host": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"sddc_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vxlan_subnet": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			// TODO check the deprecation statement
			"delay_account_link": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			// TODO change default to AWS
			"provider_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "ZEROCLOUD",
			},
			"skip_creating_vxlan": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"sso_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "vmc.local",
			},
			"sddc_template_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"deployment_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "SingleAZ",
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-west-2",
			},
		},
	}
}

func resourceSddcCreate(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	orgID := d.Get("org_id").(string)
	storageCapacity := d.Get("storage_capacity").(int)
	sddcName := d.Get("sddc_name").(string)
	vpcCidr := d.Get("vpc_cidr").(string)
	numHost := d.Get("num_host").(int)
	sddcType := d.Get("sddc_type").(string)
	vxlanSubnet := d.Get("vxlan_subnet").(string)
	accountLinkConfig := &vmc.AccountLinkConfig{
		DelayAccountLink: d.Get("delay_account_link").(bool),
	}
	providerType := d.Get("provider_type").(string)
	skipCreatingVxlan := d.Get("skip_creating_vxlan").(bool)
	ssoDomain := d.Get("sso_domain").(string)
	sddcTemplateID := d.Get("sddc_template_id").(string)
	deploymentType := d.Get("deployment_type").(string)
	region := d.Get("region").(string)
	accountLinkSddcConfig := expandAccountLinkSddcConfig(d.Get("account_link_sddc_config").([]interface{}))

	var awsSddcConfig = &vmc.AwsSddcConfig{
		StorageCapacity:       int64(storageCapacity),
		Name:                  sddcName,
		VpcCidr:               vpcCidr,
		NumHosts:              int32(numHost),
		SddcType:              sddcType,
		VxlanSubnet:           vxlanSubnet,
		AccountLinkConfig:     accountLinkConfig,
		Provider:              providerType,
		SkipCreatingVxlan:     skipCreatingVxlan,
		AccountLinkSddcConfig: accountLinkSddcConfig,
		SsoDomain:             ssoDomain,
		SddcTemplateId:        sddcTemplateID,
		DeploymentType:        deploymentType,
		Region:                region,
	}

	// Create a Sddc
	task, resp, err := vmcClient.SddcApi.OrgsOrgSddcsPost(context.Background(), orgID, *awsSddcConfig)
	if err != nil {
		return fmt.Errorf("Error while creating sddc %s: %v", sddcName, err)
	}

	// Wait until Sddc is created
	sddcID := task.ResourceId
	err = vmc.WaitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}

	// Get Sddc detail
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("Sddc %s was not found", sddcID)
	}

	d.SetId(sddc.Id)
	d.Set("name", sddc.Name)
	d.Set("updated", sddc.Updated)
	d.Set("user_id", sddc.UserId)
	d.Set("updated_by_user_id", sddc.UpdatedByUserId)
	d.Set("created", sddc.Created)
	d.Set("version", sddc.Version)
	d.Set("updated_by_user_name", sddc.UpdatedByUserName)
	d.Set("user_name", sddc.UserName)
	d.Set("sddc_state", sddc.SddcState)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_type", sddc.SddcType)
	d.Set("provider", sddc.Provider)
	d.Set("account_link_state", sddc.AccountLinkState)
	d.Set("sddc_access_state", sddc.SddcAccessState)
	d.Set("sddc_type", sddc.SddcType)

	return resourceSddcRead(d, m)
}

func resourceSddcRead(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	sddc, resp, err := vmcClient.SddcApi.OrgsOrgSddcsSddcGet(context.Background(), orgID, sddcID)
	if err != nil {
		return fmt.Errorf("Error while getting sddc detail %s: %v", sddcID, err)
	}

	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	}

	d.SetId(sddc.Id)
	d.Set("org_id", sddc.OrgId)
	d.Set("sddc_name", sddc.Name)
	d.Set("provider_type", sddc.Provider)
	d.Set("created", sddc.Created)
	return nil
}

func resourceSddcDelete(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)
	task, _, err := vmcClient.SddcApi.OrgsOrgSddcsSddcDelete(context.Background(), orgID, sddcID, nil)
	if err != nil {
		return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
	}
	err = vmc.WaitForTask(vmcClient, orgID, task.Id)
	if err != nil {
		return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
	}
	d.SetId("")
	return nil
}

func resourceSddcUpdate(d *schema.ResourceData, m interface{}) error {
	vmcClient := m.(*vmc.APIClient)
	sddcID := d.Id()
	orgID := d.Get("org_id").(string)

	// Add,remove hosts
	if d.HasChange("num_host") {
		oldTmp, newTmp := d.GetChange("num_host")
		oldNum := oldTmp.(int)
		newNum := newTmp.(int)

		action := "add"
		diffNum := newNum - oldNum

		if newNum < oldNum {
			action = "remove"
			diffNum = oldNum - newNum
		}

		esxConfig := vmc.EsxConfig{
			NumHosts: int32(diffNum),
		}

		actionString := optional.NewString(action)

		// API_CALL
		task, _, err := vmcClient.EsxApi.OrgsOrgSddcsSddcEsxsPost(
			context.Background(), orgID, sddcID, esxConfig, &vmc.OrgsOrgSddcsSddcEsxsPostOpts{Action: actionString})

		if err != nil {
			return fmt.Errorf("Error while deleting sddc %s: %v", sddcID, err)
		}
		err = vmc.WaitForTask(vmcClient, orgID, task.Id)
		if err != nil {
			return fmt.Errorf("Error while waiting for task %s: %v", task.Id, err)
		}
	}

	// Update sddc name
	if d.HasChange("sddc_name") {
		sddcPatchRequest := vmc.SddcPatchRequest{
			Name: d.Get("sddc_name").(string),
		}
		sddc, _, err := vmcClient.SddcApi.OrgsOrgSddcsSddcPatch(context.Background(), orgID, sddcID, sddcPatchRequest)

		if err != nil {
			return fmt.Errorf("Error while updating sddc's name %v", err)
		}
		d.Set("sddc_name", sddc.Name)
	}

	return resourceSddcRead(d, m)
}

func expandAccountLinkSddcConfig(l []interface{}) []vmc.AccountLinkSddcConfig {
	if len(l) == 0 {
		return nil
	}

	var configs []vmc.AccountLinkSddcConfig

	for _, config := range l {
		c := config.(map[string]interface{})
		var subnetIds []string
		for _, subnetID := range c["customer_subnet_ids"].([]interface{}) {

			subnetIds = append(subnetIds, subnetID.(string))
		}

		con := vmc.AccountLinkSddcConfig{
			CustomerSubnetIds:  subnetIds,
			ConnectedAccountId: c["connected_account_id"].(string),
		}

		configs = append(configs, con)
	}
	return configs
}
