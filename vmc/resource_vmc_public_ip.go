// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package vmc

import (
	"fmt"
	"strings"

	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/nsx_vmc_app/infra"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration/nsx_vmc_app/model"

	"github.com/gofrs/uuid/v5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Update: resourcePublicIPUpdate,
		Delete: resourcePublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")
				if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%q), expected public_ip_id,nsxt_reverse_proxy_url", d.Id())
				}
				if err := IsValidUUID(idParts[0]); err != nil {
					return nil, fmt.Errorf("invalid format for public_ip_id : %v", err)
				}
				if err := IsValidURL(idParts[1]); err != nil {
					return nil, fmt.Errorf("invalid format for nsxt_reverse_proxy_url : %v", err)
				}
				d.SetId(idParts[0])
				if err := d.Set("nsxt_reverse_proxy_url", idParts[1]); err != nil {
					return nil, fmt.Errorf("error setting nsxt_reverse_proxy_url: %v", err)
				}

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

func resourcePublicIPCreate(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
	connectorWrapper := m.(*connector.Wrapper)
	connector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	publicIpsClient := infra.NewPublicIpsClient(connector)

	displayName := d.Get("display_name").(string)
	// generate random UUID
	UUIDObject, err := uuid.NewV4()
	if err != nil {
		return HandleCreateError("Public IP", err)
	}
	UUIDStr := UUIDObject.String()

	// set values in public IP model struct
	var publicIPModel = &model.PublicIp{
		DisplayName: &displayName,
		Id:          &UUIDStr,
	}

	// API call to create public IP
	publicIP, err := publicIpsClient.Update(UUIDStr, *publicIPModel)
	if err != nil {
		return HandleCreateError("Public IP", err)
	}

	d.SetId(*publicIP.Id)
	return resourcePublicIPRead(d, m)
}

func resourcePublicIPRead(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
	connectorWrapper := m.(*connector.Wrapper)
	connector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	publicIpsClient := infra.NewPublicIpsClient(connector)
	uuid := d.Id()

	if len(uuid) > 0 {
		publicIP, err := publicIpsClient.Get(uuid)
		if err != nil {
			return HandleReadError(d, "Public IP", uuid, err)
		}
		if err := d.Set("ip", publicIP.Ip); err != nil {
			return fmt.Errorf("error setting ip: %v", err)
		}
		if err := d.Set("display_name", publicIP.DisplayName); err != nil {
			return fmt.Errorf("error setting display_name: %v", err)
		}
	} else {
		displayName := d.Get("display_name").(string)
		if len(displayName) > 0 {
			// get the list of IPs
			publicIPResultList, err := publicIpsClient.List(nil, nil, nil, nil, nil)
			if err != nil {
				return HandleListError("Public IP", err)
			}
			publicIpsList := publicIPResultList.Results
			for _, publicIP := range publicIpsList {
				if displayName == *publicIP.DisplayName {
					if err := d.Set("ip", publicIP.Ip); err != nil {
						return fmt.Errorf("error setting ip: %v", err)
					}

					if err := d.Set("display_name", publicIP.DisplayName); err != nil {
						return fmt.Errorf("error setting display_name: %v", err)
					}
					break
				}
			}
		}
	}
	return nil
}

func resourcePublicIPUpdate(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
	connectorWrapper := m.(*connector.Wrapper)
	connector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	publicIpsClient := infra.NewPublicIpsClient(connector)

	if d.HasChange("display_name") {
		uuid := d.Id()
		displayName := d.Get("display_name").(string)

		// set values in public IP model struct
		var publicIPModel = &model.PublicIp{
			DisplayName: &displayName,
			Id:          &uuid,
		}

		// API call to update public IP
		publicIP, err := publicIpsClient.Update(uuid, *publicIPModel)
		if err != nil {
			return HandleUpdateError("Public IP", err)
		}

		if err := d.Set("display_name", publicIP.DisplayName); err != nil {
			return fmt.Errorf("error setting display_name: %v", err)
		}
	}

	return resourcePublicIPRead(d, m)
}

func resourcePublicIPDelete(d *schema.ResourceData, m interface{}) error {
	nsxtReverseProxyURL := d.Get("nsxt_reverse_proxy_url").(string)
	connectorWrapper := m.(*connector.Wrapper)
	connector, err := getNsxtReverseProxyURLConnector(nsxtReverseProxyURL, connectorWrapper)
	if err != nil {
		return HandleCreateError("NSXT reverse proxy URL connector", err)
	}
	publicIpsClient := infra.NewPublicIpsClient(connector)
	uuid := d.Id()
	forceDelete := true
	err = publicIpsClient.Delete(uuid, &forceDelete)
	if err != nil {
		return HandleDeleteError("Public IP", uuid, err)
	}
	d.SetId("")
	return nil
}
