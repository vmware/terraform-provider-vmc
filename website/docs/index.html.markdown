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
token from the Cloud Service Provider by providing the org scoped refresh token. 
The Terraform provider client uses Cloud Service Provider CSP API 
to exchange this org scoped refresh token for user access token. 


There are also a number of other parameters that can be set to tune how the
provider connects to the VMC Console API. 

Note that in all of the examples you will need to update the `refresh_token`,
`sddc_name`, and `org_id` settings to match those configured in your VMC
environment.

### Example of Configuration with Credentials

```hcl
provider "vmc" {
  refresh_token = "Ih89XXXXX"
}

```

## Argument Reference

The following arguments are used to configure the VMware VMC Provider:

* `org_id` - (Required) Organization identifier.
* `region` - (Required) The region of the cloud resources to work in.
* `num_hosts` - (Optional) The number of hosts.
* `connected_account_id` - (Optional) The linked connected account identifier.
* `sddc_id` - (Optional) SDDC identifier.
* `force_refresh` - (Optional) Boolean value when set to true, forces the mappings for datacenters 
   to be refreshed for the connected account.
* `instance_type` - (Optional) The server instance type to be used.
* `sddc_id` - (Required) SDDC identifier.
* `private_ip` - (Required) Workload VM private IP to be assigned the public IP just allocated.
* `name` - (Required) Workload VM private IPs to be assigned the public IP just allocated.
* `region` - (Required)  The region of the cloud resources to work in.
* `sddc_name` - (Required) Name of the SDDC.
* `storage_capacity` - (Optional) The storage capacity value to be requested for the sddc 
   primary cluster, in GiBs. If provided, instead of using the direct-attached storage, a 
   capacity value amount of seperable storage will be used.
* `num_host` - (Required) The number of hosts.
* `account_link_sddc_config` - (Optional) The account linking configuration object.
* `vpc_cidr` - (Optional) AWS VPC IP range. Only prefix of 16 or 20 is currently supported.
* `sddc_type` - (Optional) Denotes the sddc type , if the value is null or empty, the type 
   is considered as default.
* `vxlan_subnet` - (Optional) VXLAN IP subnet in CIDR for compute gateway.
* `delay_account_link` - (Optional)  Boolean flag identifying whether account linking 
   should be delayed or not for the SDDC.
* `provider_type` - (Optional)  Determines what additional properties are available 
   based on cloud provider.
* `skip_creating_vxlan` - (Optional) Boolean value to skip creating vxlan for 
   compute gateway for SDDC provisioning.
* `sso_domain` - (Optional) The SSO domain name to use for vSphere users. 
   If not specified, vmc.local will be used.
* `sddc_template_id` - (Optional) If provided, configuration from the template will 
   applied to the provisioned SDDC.
* `deployment_type` - (Optional) Denotes if request is for a SingleAZ 
   or a MultiAZ SDDC. Default is SingleAZ.

#### Example vmc.tf file

This file will define the logical topology that Terraform will
create in VMC.

```hcl
#
# The first step is to configure the VMware VMC provider to connect to Cloud Service 
# Provider 

provider "vmc" {
  refresh_token = ""
}

# This part of the example shows some data sources we will need to refer to
# later in the .tf file. They include the org, connected accounts and 
# customer subnets.
data "vmc_org" "my_org" {
  id = ""
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = "${data.vmc_org.my_org.id}"
  connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
  region               = "US_WEST_2"
}

# This shows the settings required to provision an SDDC in the target cloud.
resource "vmc_sddc" "sddc_1" {
  org_id = "${data.vmc_org.my_org.id}"

  sddc_name           = ""
  vpc_cidr            = "10.2.0.0/16"
  num_host            = 1
  provider_type       = "AWS"
  region              = "${data.vmc_customer_subnets.my_subnets.region}"
  vxlan_subnet        = "192.168.1.0/24"
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  deployment_type = "SingleAZ"

  account_link_sddc_config = [
    {
      customer_subnet_ids  = ["${data.vmc_customer_subnets.my_subnets.ids.0}"]
      connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
    },
  ]
  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }
}

```


