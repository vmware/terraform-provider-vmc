---
layout: "vmc"

page_title: "VMC: vmc_sddc"
sidebar_current: "docs-vmc-resource-sddc"

description: |-
  Provides a resource to provision SDDC.
---

# vmc_sddc

Provides a resource to provision SDDC.

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  region               = var.sddc_region
}

resource "vmc_sddc" "sddc_1" {
  sddc_name           = var.sddc_name
  vpc_cidr            = var.vpc_cidr
  num_host            = 3
  provider_type       = "AWS"
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = "SingleAZ"

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  }
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}

```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) Organization identifier.

* `region` - (Required)  The AWS specific (e.g us-west-2) or VMC specific region (e.g US_WEST_2) of the cloud resources to work in.

* `sddc_name` - (Required) Name of the SDDC.

* `storage_capacity` - (Optional) The storage capacity value to be requested for the SDDC primary cluster. 
   This variable is only for R5.METAL. Possible values are 15TB, 20TB, 25TB, 30TB, 35TB per host.

* `num_host` - (Required) The number of hosts.

* `account_link_sddc_config` - (Optional) The account linking configuration object.

* `host_instance_type` -  (Optional) The instance type for the esx hosts in the primary cluster of the SDDC. Possible values : I3_METAL, R5_METAL.

* `vpc_cidr` - (Optional) SDDC management network CIDR. Only prefix of 16, 20 and 23 are supported. Note : Specify a private subnet range (RFC 1918) to be used for 
   vCenter Server, NSX Manager, and ESXi hosts. Choose a range that will not conflict with other networks you will connect to this SDDC.
   Minimum CIDR sizes : /23 for up to 27 hosts, /20 for up to 251 hosts, /16 for up to 4091 hosts.
   Reserved CIDRs : 10.0.0.0/15, 172.31.0.0/16.
 
* `sddc_type` - (Optional) Denotes the sddc type , if the value is null or empty, the type is considered
   as default.

* `vxlan_subnet` - (Optional) A logical network segment that will be created with the SDDC under the compute gateway.

* `delay_account_link` - (Optional)  Boolean flag identifying whether account linking should be delayed
   or not for the SDDC.

* `provider_type` - (Optional)  Determines what additional properties are available based on cloud
   provider. Default value : AWS

* `skip_creating_vxlan` - (Optional) Boolean value to skip creating vxlan for compute gateway for SDDC provisioning.

* `sso_domain` - (Optional) The SSO domain name to use for vSphere users. If not specified, vmc.local will be used.

* `sddc_template_id` - (Optional) If provided, configuration from the template will applied to the provisioned SDDC.

* `deployment_type` - (Optional) Denotes if request is for a SingleAZ or a MultiAZ SDDC. Default : SingleAZ.

* `cluster_id` - (Optional) Cluster identifier.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - SDDC identifier.

* `nsxt_reverse_proxy_url` - NSXT reverse proxy url for managing public IP.
