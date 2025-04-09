---
page_title: "VMC: vmc_sddc"
description: A resource for provisioning an SDDC.
---

# Resource: vmc_sddc

Provides a resource to provision an SDDC.

## Deploying a SingleAZ SDDC

For the `deployment_type` of `SingleAZ`, the `sddc_type` can be `1NODE` with
`num_host` argument set to `1`. The `sddc_type` for `num_host` set to 2 or
greater is `DEFAULT`.

## Example

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
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
  num_host            = 1
  provider_type       = "AWS"
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"
  deployment_type     = "SingleAZ"
  sddc_type           = "1NODE"
  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }
  microsoft_licensing_config {
    mssql_licensing   = "ENABLED"
    windows_licensing = "DISABLED"
  }
}
```

## Modifying an Elastic DRS Policy

In a new SDDC, elastic DRS uses the **Default Storage Scale-Out** policy, adding
hosts only when storage utilization exceeds the threshold of 75%. For two-host
SDDCs, only the **Default Storage Scale-Out** policy is available. Elastic DRS
is not supported for single ESX host (`1Node`) SDDCs.

You can select a different policy if it provides better support for your
workload VMs by updating the resource using the following arguments :

* `edrs_policy_type` - (Optional) The EDRS policy type. This can either be
  `cost`, `performance`, `storage-scaleup` or `rapid-scaleup`. Defaults to
  storage-scaleup.

* `enable_edrs` - (Optional) Enable EDRS.

* `min_hosts` - (Optional) The minimum number of ESX hosts that the cluster can
  scale in to.

* `max_hosts` - (Optional) The maximum number of ESX hosts that the cluster can
  scale out to.

~> **Note:** When the EDRS policy is disabled (*i.e.*, `enable_edrs = false`)
for `performance`, `cost` or `rapid-scaleup`, the policy type changes to the
default, `storage-scaleup`.

~> **Note:** The EDRS policy properties can be modified only after an SDDC has
been deployed.

## Example

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
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
  deployment_type     = "SingleAZ"

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }
  microsoft_licensing_config {
    mssql_licensing   = "ENABLED"
    windows_licensing = "DISABLED"
  }
  edrs_policy_type = "cost"
  enable_edrs      = true
  min_hosts        = 3
  max_hosts        = 8
}
```

## Deploying a MultiAZ SDDC (Stretched Cluster)

For deployment type `MultiAZ`, a single SDDC can be deployed across two AWS
availability zones.

When enabled the default number of ESX hosts supported in a `MultiAZ` SDDC is 6.
Additional hosts can be added later but must be done in pairs across AWS
availability zones. The `MultiAZ` SDDC requires an AWS VPC with two subnets, one
subnet per availability zone.

## Example

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
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
  num_host            = 6
  provider_type       = var.provider_type
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"
  deployment_type     = "MultiAZ"
  host_instance_type  = var.host_instance_type
  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0], data.vmc_customer_subnets.my_subnets.ids[1]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }
  microsoft_licensing_config {
    mssql_licensing   = "ENABLED"
    windows_licensing = "DISABLED"
  }
}
```

## Argument Reference

The following arguments are supported for this resource:

* `org_id` - (Required) The organization identifier.

* `region` - (Required) The AWS specific (*e.g.*, `us-west-2`) or VMC specific
  region (*e.g.*, `US_WEST_2`) of the cloud resources to work in.

* `sddc_name` - (Required) The name of the SDDC.

* `num_host` - (Required) The number of ESX hosts in the primary cluster of the
  SDDC.

* `size` - (Optional) The size of the vCenter and NSX appliances. `large` or
  `LARGE` SDDC size corresponds to a large vCenter appliance and large NSX
  appliance. `medium` or `MEDIUM` SDDC size corresponds to medium vCenter
  appliance and medium NSX appliance. Defaults to `medium`.

* `account_link_sddc_config` - (Optional) The account linking configuration
  object.

* `host_instance_type` - (Optional) The instance type for the ESX hosts in the
  primary cluster of the SDDC. Allows values include: `I3_METAL`, `I3EN_METAL`,
  `I4I_METAL`, and `R5_METAL`. Defaults to `I3_METAL`. Currently, `I3EN_METAL`
  does not support `1NODE` and 2 node SDDC deployment.

* `vpc_cidr` - (Optional) SDDC management network CIDR. Only prefix of `16`,
  `20` and `23` are supported.

~> **Note:** Specify a private subnet range (RFC 1918) to be used for vCenter,
NSX Manager, and ESX hosts. Choose a range that will not conflict with other
networks you will connect to this SDDC. Minimum CIDR sizes: `/23` for up to 27
hosts; `/20` for up to 251 hosts, and `/16` for up to 4091 hosts.

~> **Note:** Reserved CIDRs: `10.0.0.0/15` and `172.31.0.0/16`.

* `sddc_type` - (Optional) Specifies the SDDC type, if the value is `null` or
  empty, the type is considered as default.

* `vxlan_subnet` - (Optional) A logical network segment that will be created
  with the SDDC under the compute gateway.

* `delay_account_link` - (Optional) Specifics whether account linking should be
  delayed or not for the SDDC.

* `provider_type` - (Optional) Determines what additional properties are
  available based on cloud provider. Defaults to `AWS`.

* `skip_creating_vxlan` - (Optional) Specifies to skip creating VXLAN for
  compute gateway for SDDC provisioning.

* `sso_domain` - (Optional) The SSO domain name to use for vSphere users. If not
  specified, `vmc.local` will be used.

* `sddc_template_id` - (Optional) If provided, configuration from the template
  will be applied to the provisioned SDDC.

* `deployment_type` - (Optional) Specifies if the type is for a `SingleAZ` or a
  `MultiAZ` SDDC. Defaults to `SingleAZ`.

* `cluster_id` - (Optional) The cluster identifier.

* `microsoft_licensing_config` - (Optional) Indicates the desired licensing
  support, if any, of Microsoft software.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The SDDC identifier.

* `cluster_info` - Information about cluster such as the id, name, state, and
  host instance type.

* `sddc_size` - The size information of vCenter appliance and NSX appliance.

* `intranet_uplink_mtu` - The Uplink MTU of direct connect, sddc-grouping, and
  outposts traffic in an edge tier-0 router port. This field can be updated only
  after an SDDC is created. Range: `1500 - 8900`. Defaults to `1500`.

* `nsxt_reverse_proxy_url` - The NSX reverse proxy URL for managing public IP.

* `nsxt_cloudadmin` - The NSX `admin` user for direct access.

* `nsxt_cloudadmin_password` - The NSX `admin` user password for direct access.

* `nsxt_cloudaudit` - The NSX `audit` user for direct access.

* `nsxt_cloudaudit_password` - The NSX `audit` user password for direct access.

* `nsxt_private_url` - The NSX private URL.

* `nsxt_public_url` - Same as `nsxt_reverse_proxy_url`

## Import

Import the resource using the `id`.

`$ terraform import vmc_sddc.sddc_1 id`

For example:

`$ terraform import vmc_sddc.sddc_1 afe7a0fd-3f0a-48b2-9ddb-0489c22732ae`

~> **Note:** Running plan/apply after importing an SDDC causes the SDDC to be
re-created. This is due to a limitation in the current `GET` and `UPDATE` SDDC
APIs. Hence, the import functionality is only partially supported.
