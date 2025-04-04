---
page_title: "VMC: vmc_sddc"
description: The data source for an SDDC.
---

# Data Source: vmc_sddc

The SDDC data source retrieves information about an SDDC.

## Example Usage

```hcl
data "vmc_sddc" "my_sddc" {
  sddc_id = var.sddc_id
}
```

## Argument Reference

* `org_id` - (Required) The organization identifier.

* `sddc_id` - (Required) The SDDC identifier.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The SDDC identifier.

* `region` - The AWS specific (*e.g.*, `us-west-2`) or VMC specific region
  (*e.g.*, `US_WEST_2`) of the cloud resources to work in.

* `sddc_name` - The name of the SDDC.

* `num_host` - The number of ESX hosts.

* `sddc_type` - The SDDC type.

* `sddc_state` - The SDDC state.

* `provider_type` - Allowed values include `AWS` and `ZEROCLOUD`.

* `skip_creating_vxlan` - Specifies to skip creating VXLAN for compute gateway
  for SDDC provisioning.

* `sso_domain` - The SSO domain name to use for vSphere users.

* `deployment_type` - The deployment type. One of: `SingleAZ` or `MultiAZ`.

* `availability_zones` - The availability zones.

* `vc_url` - The vCenter instance URL.

* `cloud_username` - The cloud username.

* `nsxt_reverse_proxy_url` - The NSX reverse proxy URL for managing public IP.

* `nsxt_cloudadmin` - The NSX `admin` user for direct access.

* `nsxt_cloudadmin_password` - The NSX `admin` user password for direct access.

* `nsxt_cloudaudit` - The NSX `audit` user  for direct access.

* `nsxt_cloudaudit_password` - The NSX `audit` user password for direct access.

* `nsxt_private_url` - The NSX private URL.

* `nsxt_public_url` - Same as `nsxt_reverse_proxy_url`.
