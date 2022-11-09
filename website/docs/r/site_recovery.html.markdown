---
layout: "vmc"

page_title: "VMC: vmc_site_recovery"
sidebar_current: "docs-vmc_site_recovery"

description: |-
  Provides a resource to activate and deactivate site recovery for SDDC.
---

# vmc_site_recovery

Provides a resource to activate and deactivate site recovery for SDDC.

~> **Note:** Site recovery resource implicitly depends on SDDC resource creation. SDDC must be provisioned before a site recovery can be activated. 
A 10-minute delay must be added to SDDC resource before site recovery can be activated.
This delay is added using the local-exec provisioner. For details on how to provision SDDC refer to [vmc_sddc](https://www.terraform.io/docs/providers/vmc/r/sddc.html).

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

resource "vmc_sddc" "sddc_1" {
  sddc_name           = var.sddc_name
  vpc_cidr            = var.vpc_cidr
  num_host            = var.sddc_num_hosts
  provider_type       = "AWS"
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = "SingleAZ"

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }

  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
  
   # provisioner defined to add 10 minute delay after SDDC creation to enable site recovery activation.
   provisioner "local-exec" {
     command = "sleep 600"     
   } 
}


resource "vmc_site_recovery" "site_recovery_1" {
  sddc_id = vmc_sddc.sddc_1.id
  srm_extension_key_suffix = var.site_recovery_srm_extension_key_suffix
}


```

## Argument Reference

The following arguments are supported for vmc_site_recovery resource:

* `sddc_id` - (Required) SDDC identifier.

* `srm_node_extension_key_suffix` - (Optional) Custom extension key suffix for SRM. If not specified, default extension key will be used. 
The custom extension suffix must contain 13 characters or less, be composed of letters, numbers, ., - characters. 
The extension suffix must begin and end with a letter or number. The suffix is appended to com.vmware.vcDr- to form the full extension key.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported after site recovery activation:

* `site_recovery_state` - Site recovery state. Possible values are: ACTIVATED, ACTIVATING, CANCELED, DEACTIVATED, DEACTIVATING, DELETED, FAILED.

* `srm_node` - Site recovery node created after site recovery activation.

* `vr_node` - VR node created after site recovery activation.

## Import

Site recovery resource can be imported using the `sddc_id` , e.g.

`$ terraform import vmc_site_recovery.site_recovery_1 sddc_id`

- sddc_id = SDDC Identifier

`$ terraform import vmc_site_recovery.site_recovery_1 afe7a0fd-3f0a-48b2-9ddb-0489c22732ae`