---
layout: "vmc"
page_title: "VMC: sddc"
sidebar_current: "docs-vmc-datasource-sddc"
description: An sddc data source.
---

# vmc_sddc

The sddc data source provides information about an SDDC.
## Example Usage

```hcl
data "vmc_sddc" "my_sddc" {
  sddc_id               = var.sddc_id
}
```

## Argument Reference

* `org_id` - (Required) Organization identifier.

* `region` - (Optional)  The AWS specific (e.g us-west-2) or VMC specific region (e.g US_WEST_2) of the cloud resources to work in.

* `sddc_id` - (Required) ID of the SDDC.

* `sddc_name` - (Optional) Name of the SDDC.

* `storage_capacity` - (Optional) The storage capacity value to be requested for the sddc primary cluster,
   in GiBs. If provided, instead of using the direct-attached storage, a capacity value amount of
   separable storage will be used. Possible values for R5 metal are 15TB, 20TB, 25TB, 30TB, 35TB.

* `num_host` - (Optional) The number of hosts.

* `account_link_sddc_config` - (Optional) The account linking configuration object.

* `host_instance_type` -  (Optional) The instance type for the esx hosts in the primary cluster of the SDDC.

* `vpc_cidr` - (Optional) AWS VPC IP range. Only prefix of 16 or 20 is currently supported.

* `sddc_type` - (Optional) Denotes the sddc type , if the value is null or empty, the type is considered
   as default.

* `vxlan_subnet` - (Optional) VXLAN IP subnet in CIDR for compute gateway.

* `delay_account_link` - (Optional)  Boolean flag identifying whether account linking should be delayed
   or not for the SDDC.

* `provider_type` - (Optional)  Determines what additional properties are available based on cloud
   provider. Acceptable values are "ZEROCLOUD" and "AWS" with AWS as the default value.

* `skip_creating_vxlan` - (Optional) Boolean value to skip creating vxlan for compute gateway for SDDC provisioning.

* `sso_domain` - (Optional) The SSO domain name to use for vSphere users. If not specified, vmc.local will be used.

* `sddc_template_id` - (Optional) If provided, configuration from the template will applied to the provisioned SDDC.

* `deployment_type` - (Optional) Denotes if request is for a SingleAZ or a MultiAZ SDDC. Default is SingleAZ.

* `cluster_id` - (Optional) Cluster identifier.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - SDDC identifier.

* `nsxt_reverse_proxy_url` - NSXT reverse proxy url for managing public IP.
