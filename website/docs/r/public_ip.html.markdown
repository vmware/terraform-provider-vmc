---
layout: "vmc"

page_title: "VMC: vmc_public_ip"
sidebar_current: "docs-vmc-resource-public-ip"

description: |-
  Provides a resource to manage public IPs.
---

# vmc_public_ip

Provides a resource to manage public IPs.

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

resource "vmc_public_ip" "public_ip_1" {
  nsxt_reverse_proxy_url = vmc_sddc.sddc_1.nsxt_reverse_proxy_url
  display_name = var.public_ip_displayname
}

```

## Argument Reference

The following arguments are supported:

* `org_id` - (Required) Organization identifier.

* `nsxt_reverse_proxy_url` - (Required) NSXT reverse proxy url for managing public IP. Computed after SDDC creation.

* `display_name` - (Optional) Display name for public IP.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - IP identifier.

* `ip` - Public IP.
