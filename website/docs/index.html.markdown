---
layout: "vmc"
page_title: "Terraform Provider for VMware Cloud on AWS"
sidebar_current: "docs-vmc-index"
description: |-
  Terraform Provider for VMware Cloud on AWS
---

<img src="https://raw.githubusercontent.com/vmware/terraform-provider-vmc/main/docs/images/icon-color.png" alt="VMware Cloud on AWS" width="150">

# Terraform Provider for VMware Cloud on AWS

The Terraform Provider for [VMware Cloud on AWS][product-documentation] is a plugin for Terraform that allows you to
interact with VMware Cloud on AWS.

## Example Usage

In order to use the provider you need to obtain the authentication token from the Cloud Service Provider by providing
the org scoped API token. The provider client uses Cloud Service Provider (CSP) API to exchange this org scoped API
token/OAuth App ID and Secret for user access token.

```hcl
terraform {
  required_providers {
    vcf = {
      source = "vmware/vmc"
      version = "x.y.z"
    }
  }
}

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = var.sddc_region
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
}

resource "vmc_public_ip" "public_ip_1" {
  nsxt_reverse_proxy_url = vmc_sddc.sddc_1.nsxt_reverse_proxy_url
  display_name = var.public_ip_displayname
}
```

Refer to the provider documentation for information on all of the resources
and data sources supported by this provider. Each includes a detailed
description of the purpose and how to use it.

## Argument Reference

The following arguments are used to configure the provider:

* `api_token` - (Required, in conflict with `client_id` and `client_secret`)
  API token is used to authenticate when calling VMware Cloud Services APIs. 
  This token is scoped within the organization.
* `client_id` - (Required in pair with `client_secret`, in conflict with `api_token`)
  ID of OAuth App associated with the organization. The combination with
  "client_secret" is used to authenticate when calling VMware Cloud Services
  APIs.
* `client_secret` - (Required in pair with `client_id`, in conflict with `api_token`)
  Secret of OAuth App associated with the organization. The combination with
  "client_id" is used to authenticate when calling VMware Cloud Services APIs.
* `org_id` - (Required) Organization Identifier.
* `vmc_url` - (Optional) VMware Cloud on AWS URL. Default: https://vmc.vmware.com
* `csp_url` - (Optional) Cloud Service Provider URL. Default: https://console.cloud.vmware.com

[product-documentation]: https://docs.vmware.com/en/VMware-Cloud-on-AWS/index.html