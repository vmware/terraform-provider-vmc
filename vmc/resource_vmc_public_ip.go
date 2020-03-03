/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/api"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/model"
)

func resourcePublicIp() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIpCreate,
		Read:   resourcePublicIpRead,
		Update: resourcePublicIpUpdate,
		Delete: resourcePublicIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP associated with the SDDC",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Display name/notes about this resource",
			},
		},
	}
}

func resourcePublicIpCreate(d *schema.ResourceData, m interface{}) error {
	connector, err := getNSXTReverseProxyConnector()
	if err != nil {
		return fmt.Errorf("Error getting connector for reverse proxy url : %v", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	displayName := d.Get("display_name").(string)
	// generate random UUID
	uuid := uuid.NewV4().String()

	// set values in public IP model struct
	var publicIpModel = &model.PublicIp{
		DisplayName: &displayName,
		Id:          &uuid,
	}

	// API call to create public IP
	publicIp, err := nsxVmcAwsClient.CreatePublicIp(uuid, *publicIpModel)
	if err != nil {
		return fmt.Errorf("Error creating public IP : %v", err)
	}

	d.SetId(*publicIp.Id)
	return resourcePublicIpRead(d, m)
}

func resourcePublicIpRead(d *schema.ResourceData, m interface{}) error {
	connector, err := getNSXTReverseProxyConnector()
	if err != nil {
		return fmt.Errorf("Error getting connector for reverse proxy url : %v", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Id()

	if len(uuid) > 0 {
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return fmt.Errorf("Error getting public IP with ID %s : %v", uuid, err)
		}
		d.Set("ip", publicIp.Ip)
		d.Set("display_name", publicIp.DisplayName)
	} else {
		displayName := d.Get("display_name").(string)
		if len(displayName) > 0 {
			// get the list of IPs
			publicIpResultList, err := nsxVmcAwsClient.ListPublicIps()
			if err != nil {
				return fmt.Errorf("Error getting list of public IPs : %v", err)
			}
			publicIpsList := publicIpResultList.Results
			if publicIpsList != nil {
				for _, publicIp := range publicIpsList {
					if displayName == *publicIp.DisplayName {
						d.Set("ip", publicIp.Ip)
						d.Set("display_name", publicIp.DisplayName)
						break
					}
				}
			}
		}
	}
	return nil
}

func resourcePublicIpUpdate(d *schema.ResourceData, m interface{}) error {
	connector, err := getNSXTReverseProxyConnector()
	if err != nil {
		return fmt.Errorf("Error getting connector for reverse proxy url : %v", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)

	if d.HasChange("display_name") {
		uuid := d.Id()
		displayName := d.Get("display_name").(string)

		// set values in public IP model struct
		var publicIpModel = &model.PublicIp{
			DisplayName: &displayName,
			Id:          &uuid,
		}

		// API call to update public IP
		publicIp, err := nsxVmcAwsClient.CreatePublicIp(uuid, *publicIpModel)
		if err != nil {
			return fmt.Errorf("Error while updating public IP's display name : %v", err)
		}

		d.Set("display_name", publicIp.DisplayName)
	}

	return resourcePublicIpRead(d, m)
}

func resourcePublicIpDelete(d *schema.ResourceData, m interface{}) error {
	connector, err := getNSXTReverseProxyConnector()
	if err != nil {
		return fmt.Errorf("Error getting connector for reverse proxy url : %v", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Id()
	var forceDelete bool = true
	err = nsxVmcAwsClient.DeletePublicIp(uuid, &forceDelete)
	if err != nil {
		return fmt.Errorf("Error deleting public IP with ID %s : %v", uuid, err)
	}
	d.SetId("")
	return nil
}

func getNSXTReverseProxyConnector() (client.Connector, error) {
	apiToken := os.Getenv(APIToken)
	nsxtReverseProxyURL := os.Getenv(NSXTReverseProxyUrl)
	if len(nsxtReverseProxyURL) == 0 {
		return nil, fmt.Errorf("NSXT reverse proxy url is required for Public IP resource creation.")
	}
	if strings.Contains(nsxtReverseProxyURL, SksNSXTManager) {
		nsxtReverseProxyURL = strings.Replace(nsxtReverseProxyURL, SksNSXTManager, "", -1)
	}
	httpClient := http.Client{}
	connector, err := NewClientConnectorByRefreshToken(apiToken, nsxtReverseProxyURL, DefaultCSPUrl, httpClient)
	if err != nil {
		return nil, fmt.Errorf("Error creating connector : %v ", err)
	}
	return connector, nil
}
