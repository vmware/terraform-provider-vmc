/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/connector"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/account_link"
)

func dataSourceVmcCustomerSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcCustomerSubnetsRead,

		Schema: map[string]*schema.Schema{
			"connected_account_id": {
				Type:        schema.TypeString,
				Description: "The linked connected account identifier.",
				Optional:    true,
			},
			"region": {
				Type:        schema.TypeString,
				Description: "The VMC region of the cloud resources to work in. (e.g. US_WEST_2)",
				Required:    true,
				ValidateFunc: validation.All(
					validation.NoZeroValues,
				),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == strings.ReplaceAll(strings.ToUpper(new), "-", "_")
				},
			},
			"num_hosts": {
				Type:        schema.TypeInt,
				Description: "The number of hosts .",
				Optional:    true,
			},
			"sddc_id": {
				Type:        schema.TypeString,
				Description: "Sddc ID.",
				Optional:    true,
			},
			"sddc_type": {
				Type:        schema.TypeString,
				Description: "Sddc Type.",
				Optional:    true,
			},
			"force_refresh": {
				Type:        schema.TypeBool,
				Description: "When true, forces the mappings for datacenters to be refreshed for the connected account.",
				Optional:    true,
			},
			"instance_type": {
				Type:        schema.TypeString,
				Description: "The server instance type to be used.",
				Optional:    true,
			},
			"customer_available_zones": {
				Type:        schema.TypeList,
				Description: "A list of AWS availability zones.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"vpc_map": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"ids": {
				Type:        schema.TypeList,
				Description: "A list of AWS subnet IDs to create links to in the customer's account.",
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVmcCustomerSubnetsRead(d *schema.ResourceData, m interface{}) error {
	orgID := m.(*connector.Wrapper).OrgID
	accountID := d.Get("connected_account_id").(string)
	sddcID := d.Get("sddc_id").(string)
	region := d.Get("region").(string)
	numHosts := int64(d.Get("num_hosts").(int))
	sddcType := d.Get("sddc_type").(string)
	forceRefresh := d.Get("force_refresh").(bool)
	instanceType := d.Get("instance_type").(string)

	connectorWrapper := (m.(*connector.Wrapper)).Connector
	compatibleSubnetsClient := account_link.NewCompatibleSubnetsClient(connectorWrapper)
	compatibleSubnets, err := compatibleSubnetsClient.Get(orgID, accountID, &region, &sddcID, &forceRefresh, &instanceType, &sddcType, &numHosts)
	ids := []string{}
	for _, value := range compatibleSubnets.VpcMap {
		for _, subnet := range value.Subnets {
			ids = append(ids, *subnet.SubnetId)
		}
	}
	log.Printf("[DEBUG] Subnet IDs are %v\n", ids)

	if err != nil {
		return HandleDataSourceReadError("Customer Subnets", err)
	}

	d.Set("ids", ids)
	d.Set("customer_available_zones", compatibleSubnets.CustomerAvailableZones)
	d.SetId(fmt.Sprintf("%s-%s", orgID, accountID))
	return nil
}
