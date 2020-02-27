/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/rs/xid"
	"gitlab.eng.vmware.com/golangsdk/vsphere-automation-sdk-go/services/nsxt/vmc-aws-integration/api"
	"gitlab.eng.vmware.com/golangsdk/vsphere-automation-sdk-go/services/nsxt/vmc-aws-integration/model"
)

func resourcePublicIp() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIpCreate,
		Read:   resourcePublicIpRead,
		Delete: resourcePublicIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(300 * time.Minute),
			Update: schema.DefaultTimeout(300 * time.Minute),
			Delete: schema.DefaultTimeout(180 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of this resource",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Display name/notes about this resource",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Public IP associated with the SDDC",
			},
			"ips": {
				Type:        schema.TypeList,
				Description: "The list of all public IPs associated with the SDDC.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePublicIpCreate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	displayName := d.Get("display_name").(string)
	uuid := d.Get("ip_id").(string)

	// if UUID is empty, user is attempting to create public IP
	// else its an update to an existing IPs display name
	if len(uuid) == 0 {
		// generate random UUID
		guid := xid.New()
		uuid = guid.String()
	}

	// set values in public IP model struct
	var publicIpModel = &model.PublicIp{
		DisplayName: &displayName,
		Id:          &uuid,
	}

	// API call to create public IP
	publicIp, err := nsxVmcAwsClient.CreatePublicIp(uuid, *publicIpModel)
	if err != nil {
		return fmt.Errorf("Error while creating public IP : %v", err)
	}

	publicIpId := publicIp.Id
	d.SetId(*publicIpId)

	// Since the same API is used for create and update, set display name if a change is detected
	if d.HasChange("display_name") {
		d.Set("display_name", d.Get("name").(string))
	}

	return resourcePublicIpRead(d, m)
}

func resourcePublicIpRead(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Get("ip_id").(string)

	if len(uuid) > 0 {
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return fmt.Errorf("Error while getting public IPs : %v", err)
		}
		d.Set("ip_id", publicIp.Id)
		d.Set("ip", publicIp.Ip)
		d.Set("display_name", publicIp.DisplayName)
	} else {
		// get the list of IPs
		publicIpResultList, err := nsxVmcAwsClient.ListPublicIps()
		if err != nil {
			return fmt.Errorf("Error while getting list of public IPs associated with the SDDC : %v", err)
		}
		ips := []string{}
		displayName := d.Get("display_name")
		publicIpsList := publicIpResultList.Results
		if publicIpsList != nil {
			for _, publicIp := range publicIpsList {
				ip := *publicIp.Ip
				if displayName == publicIp.DisplayName {
					ips = append(ips, ip)
					d.Set("ip_id", publicIp.Id)
					d.Set("ip", publicIp.Ip)
					d.Set("display_name", publicIp.DisplayName)
				} else {
					ips = append(ips, ip)
				}
			}
		}
		d.Set("ips", ips)
	}
	return nil
}

func resourcePublicIpDelete(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Get("ip_id").(string)
	var forceDelete bool = true
	err := nsxVmcAwsClient.DeletePublicIp(uuid, &forceDelete)
	if err != nil {
		return fmt.Errorf("Error while deleting public IP : %v", err)
	}
	d.SetId("")
	return nil
}
