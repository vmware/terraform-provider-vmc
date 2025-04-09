---
page_title: "VMC: vmc_sddc_group"
description: A resource for adding SDDCs into an SDDC Group.
---

# Resource: vmc_sddc_group

Provides a resource to add SDDCs into an SDDC Group.

## Example

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
}
# Empty data source defined in order to store the org display name and name in terraform state
data "vmc_org" "my_org" {
}

data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = replace(upper(var.sddc_region), "-", "_")
}

resource "vmc_sddc_group" "sddc_group" {
  name            = var.sddc_group_name
  description     = var.sddc_group_description
  sddc_member_ids = [vmc_sddc.sddc_1.id, vmc_sddc.sddc_2.id]
}

resource "vmc_sddc" "sddc_1" {
  sddc_name           = var.sddc1_name
  vpc_cidr            = var.vpc1_cidr
  num_host            = var.sddc_primary_cluster_num_hosts
  provider_type       = var.provider_type
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"
  sddc_type           = var.sddc_type
  deployment_type     = var.deployment_type
  size                = var.size
  host_instance_type  = var.host_instance_type

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }

  microsoft_licensing_config {
    mssql_licensing   = "ENABLED"
    windows_licensing = "DISABLED"
  }

  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}

resource "vmc_sddc" "sddc_2" {
  sddc_name           = var.sddc2_name
  vpc_cidr            = var.vpc2_cidr
  num_host            = var.sddc_primary_cluster_num_hosts
  provider_type       = var.provider_type
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"
  sddc_type           = var.sddc_type
  deployment_type     = var.deployment_type
  size                = var.size
  host_instance_type  = var.host_instance_type

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }

  microsoft_licensing_config {
    mssql_licensing   = "ENABLED"
    windows_licensing = "DISABLED"
  }

  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}
```

## Argument Reference

The following arguments are supported for this resource:

* `name` - (Required) Name of the SDDC Group.

* `description` - (Required)  Short description of the SDDC Group.

* `sddc_member_ids` - (Required) IDs of the SDDCs to be included as members in
  the SDDC Group.
