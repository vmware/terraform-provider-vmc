---
page_title: "VMC: vmc_public_ip"
description: A resource for managing public IPs.
---

# Resource: vmc_public_ip

Provides a resource to manage public IPs.

~> **Note:** Public IP resource implicitly depends on SDDC resource creation.
SDDC must be provisioned before a public IP can be created.

## Example Usage

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
}

resource "vmc_public_ip" "public_ip_1" {
  nsxt_reverse_proxy_url = vmc_sddc.sddc_1.nsxt_reverse_proxy_url
  display_name           = var.public_ip_displayname
}
```

## Argument Reference

The following arguments are supported for this resource:

* `nsxt_reverse_proxy_url` - (Required) The NSX reverse proxy URL for managing
  public IP. Computed after SDDC creation.

* `display_name` - (Optional) Display name for public IP.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - Public IP identifier.

* `ip` - Public IP.

* `display_name` - Display name for public IP.

## Import

Import the resource using the `id` and `nsxt_reverse_proxy_url`.

`$ terraform import vmc_public_ip.public_ip_1 nsxt_reverse_proxy_url,id`

For example:

`$ terraform import vmc_public_ip.public_ip_1 'https://nsx-44-228-76-55.rp.vmwarevmc.com/vmc/reverse-proxy/api/orgs/{orgI}/sddcs/afe7a0fd-3f0a-48b2-9ddb-0489c22732ae/sks-nsxt-manager,8d730ad4-aa6b-4f9f-9679-ec17beeaceaf'`
