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
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" || !IsValidUUID(idParts[1]) {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected nsxt_reverse_proxy_url,public_ip_id", d.Id())
				}
				d.Set("nsxt_reverse_proxy_url", idParts[0])
				d.SetId(idParts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"nsxt_reverse_proxy_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "NSX API public endpoint url used for public IP resource management",
			},
			"ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP associated with the SDDC",
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Display name/notes about this resource",
			},
		},
	}
}

func resourcePublicIpCreate(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyUrl := d.Get("nsxt_reverse_proxy_url").(string)
	connector, err := getNSXTReverseProxyUrlConnector(nsxtReverseProxyUrl)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
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
		return HandleCreateError("Public IP", err)
	}

	d.SetId(*publicIp.Id)
	return resourcePublicIpRead(d, m)
}

func resourcePublicIpRead(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyUrl := d.Get("nsxt_reverse_proxy_url").(string)
	connector, err := getNSXTReverseProxyUrlConnector(nsxtReverseProxyUrl)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Id()

	if len(uuid) > 0 {
		publicIp, err := nsxVmcAwsClient.GetPublicIp(uuid)
		if err != nil {
			return HandleReadError(d, "Public IP", uuid, err)
		}
		d.Set("ip", publicIp.Ip)
		d.Set("display_name", publicIp.DisplayName)
	} else {
		displayName := d.Get("display_name").(string)
		if len(displayName) > 0 {
			// get the list of IPs
			publicIpResultList, err := nsxVmcAwsClient.ListPublicIps()
			if err != nil {
				return HandleListError("Public IP", err)
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
	nsxtReverseProxyUrl := d.Get("nsxt_reverse_proxy_url").(string)
	connector, err := getNSXTReverseProxyUrlConnector(nsxtReverseProxyUrl)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
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
			return HandleUpdateError("Public IP", err)
		}

		d.Set("display_name", publicIp.DisplayName)
	}

	return resourcePublicIpRead(d, m)
}

func resourcePublicIpDelete(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyUrl := d.Get("nsxt_reverse_proxy_url").(string)
	connector, err := getNSXTReverseProxyUrlConnector(nsxtReverseProxyUrl)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	nsxVmcAwsClient := api.NewDefaultNsxVmcAwsIntegrationClient(connector)
	uuid := d.Id()
	forceDelete := true
	err = nsxVmcAwsClient.DeletePublicIp(uuid, &forceDelete)
	if err != nil {
		return HandleDeleteError("Public IP", uuid, err)
	}
	d.SetId("")
	return nil
}

func getNSXTReverseProxyUrlConnector(nsxtReverseProxyUrl string) (client.Connector, error) {
	apiToken := os.Getenv(APIToken)
	if len(nsxtReverseProxyUrl) == 0 {
		return nil, fmt.Errorf("NSX reverse proxy url is required for public IP resource creation")
	}
	if strings.Contains(nsxtReverseProxyUrl, SksNSXTManager) {
		nsxtReverseProxyUrl = strings.Replace(nsxtReverseProxyUrl, SksNSXTManager, "", -1)
	}
	httpClient := http.Client{}
	connector, err := NewClientConnectorByRefreshToken(apiToken, nsxtReverseProxyUrl, DefaultCSPUrl, httpClient)
	if err != nil {
		return nil, HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	return connector, nil
}
