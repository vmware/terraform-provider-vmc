package vmc

import (
	"context"
	"fmt"
	"github.com/antihax/optional"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.eng.vmware.com/vapi-sdk/vmc-go-sdk/vmc"
)

func dataSourceVmcCustomerSubnets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVmcCustomerSubnetsRead,

		Schema: map[string]*schema.Schema{
			"org_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Organization identifier.",
				Required:    true,
			},
			"connected_account_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The linked connected account identifier.",
				Optional:    true,
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The region of the cloud resources to work in.",
				Required:    true,
			},
			"sddc_id": &schema.Schema{
				Type:        schema.TypeString,
				Description: "Sddc ID.",
				Optional:    true,
			},
			"customer_available_zones": {
				Type:        schema.TypeList,
				Description: "A list of AWS subnet IDs to create links to in the customer's account.",
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
	vmcClient := m.(*vmc.APIClient)
	orgID := d.Get("org_id").(string)
	accountID := d.Get("connected_account_id").(string)
	sddcID := d.Get("sddc_id").(string)
	region := d.Get("region").(string)
	// providerString := optional.NewString(providerType)

	getOpts := &vmc.OrgsOrgAccountLinkCompatibleSubnetsGetOpts{
		LinkedAccountId: optional.NewString(accountID),
		Region:          optional.NewString(region),
		Sddc:            optional.NewString(sddcID),
	}

	compSubnets, _, err := vmcClient.AccountLinkingApi.OrgsOrgAccountLinkCompatibleSubnetsGet(
		context.Background(), orgID, getOpts)

	ids := []string{}
	for _, value := range compSubnets.VpcMap {
		for _, subnet := range value.Subnets {
			ids = append(ids, subnet.SubnetId)
		}
	}

	// for _, subnet := range subnets.VpcMap["VpcInfoSubnets"].Subnets {
	// 	ids = append(ids, subnet.SubnetId)
	// }
	log.Printf("[DEBUG] Subnet IDs are %v\n", ids)

	if err != nil {
		return fmt.Errorf("Error while reading subnets IDs from org %q: %v", orgID, err)
	}

	d.Set("ids", ids)
	d.Set("customer_available_zones", compSubnets.CustomerAvailableZones)
	d.SetId(fmt.Sprintf("%s-%s", orgID, accountID))
	return nil
}
