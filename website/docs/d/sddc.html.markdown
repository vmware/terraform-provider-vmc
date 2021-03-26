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

* `sddc_id` - (Required) ID of the SDDC.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - SDDC identifier.

* `region` - The AWS specific (e.g us-west-2) or VMC specific region (e.g US_WEST_2) of the cloud resources to work in.

* `sddc_name` -  Name of the SDDC.

* `num_host` - The number of hosts.

* `sddc_type` - Denotes the sddc type.

* `sddc_state` - Denotes the sddc state.

* `provider_type` - Possible values are "ZEROCLOUD" and "AWS".

* `skip_creating_vxlan` - Boolean value to skip creating vxlan for compute gateway for SDDC provisioning.

* `sso_domain` - The SSO domain name to use for vSphere users. If not specified, vmc.local will be used.

* `deployment_type` - Denotes if request is for a SingleAZ or a MultiAZ SDDC. Default is SingleAZ.

* `availability_zones` - Availability Zones.

* `vc_url` - VC URL.

* `cloud_username` - Cloud user name.

* `cloud_password` - Cloud password.

* `nsxt_reverse_proxy_url` - NSXT reverse proxy url for managing public IP.
