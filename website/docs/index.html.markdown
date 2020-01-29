---
layout: "vmc"
page_title: "Provider: VMC"
sidebar_current: "docs-vmc-index"
description: |-
  The Terraform Provider for VMware Cloud
---

# terraform-provider-vmware-cloud

The terraform-provider-vmware-cloud gives the VMC administrator a way to automate features
of VMware Cloud on AWS using the VMC API.

More information on VMC can be found on the [VMC Product
Page](https://cloud.vmware.com/vmc-aws)

Please use the navigation to the left to read about available data sources and
resources.

## Basic Configuration of the terraform-provider-vmware-cloud

In order to use the terraform-provider-vmware-cloud you need to obtain the authentication
token from the Cloud Service Provider by providing the org scoped API token. 
The Terraform provider client uses Cloud Service Provider CSP API 
to exchange this org scoped API token for user access token. 

There are also a number of other parameters that can be set to tune how the
provider connects to the VMC Console API. 

Note that in all of the examples you will need to update the `api_token` and `org_id` settings 
in the variables.tf file to match those configured in your VMC environment.


## Argument Reference

The following arguments are used to configure the VMware VMC Provider:

* `api_token` - (Required) API token is used to authenticate when calling VMware Cloud Services APIs. 
   This token is scoped within the organization.
*  `org_id` - (Required) Organization Identifier.
*  `vmc_url` - (Required) VMware Cloud on AWS URL.
*  `csp_url` - (Required) Cloud Service Provider URL.

#### Example main.tf file

This file will define the logical topology that Terraform will
create in VMC.

```hcl
#
# The first step is to configure the provider to connect to Cloud Service 
# Provider. 

provider "vmc" {
  refresh_token = var.api_token
}

data "vmc_org" "my_org" {
  id =  var.org_id
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = data.vmc_org.my_org.id
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = data.vmc_org.my_org.id
  connected_account_id = data.vmc_connected_accounts.my_accounts.ids[0]
  region               = var.sddc_region
}

resource "vmc_sddc" "sddc_1" {
  org_id = data.vmc_org.my_org.id

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



