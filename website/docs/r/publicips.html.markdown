---
layout: "vmc"
page_title: "VMC: vmc_publicips"
sidebar_current: "docs-vmc-resource-publicips"
description: |- 
  
---

# vmc_publicips

Provides a resource to allocate public IPs for an SDDC.

## Example Usage

```hcl
data "vmc_org" "my_org" {
	id = ""
}
resource "vmc_sddc" "sddc_1" {
	org_id = "${data.vmc_org.my_org.id}"

	# storage_capacity    = 100
	sddc_name = "SDDC_Name"

	vpc_cidr      = "10.2.0.0/16"
	num_host      = 3
	provider_type = "Provider_type"

	region = "US_WEST_2"

	vxlan_subnet = "192.168.1.0/24"

	delay_account_link  = false
	skip_creating_vxlan = false
	sso_domain          = "vmc.local"

	deployment_type = "SingleAZ"
}

resource "vmc_publicips" "publicip_1" {
	org_id = "${data.vmc_org.my_org.id}"
	sddc_id = "${data.vmc_sddc.sddc_1.id}"
	name     = "DefaultIP1"
	private_ip = "10.105.167.133"
}
```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) Organization identifier.

* `sddc_id` - (Required) SDDC identifier.

* `private_ip` - (Required) Workload VM private IP to be assigned the public IP just allocated.

* `name` - (Required) Workload VM private IPs to be assigned the public IP just allocated.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `allocation ID` - IP allocation identifier.

* `public_ip` - public IP allocated to the SDDC.

* `dnat_rule_id` - dnat rule identifier.

* `snat_rule_id` - snat rule identifier.