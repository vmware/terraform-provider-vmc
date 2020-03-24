---
layout: "vmc"

page_title: "VMC: vmc_site_recovery"
sidebar_current: "docs-vmc-resource-site-recovery"

description: |-
  Provides a resource to manage site recovery.
---

# vmc_public_ip

Provides a resource to manage site recovery.
~> **Note:** Site recovery resource implicitly depends on SSDC resource creation. SDDC must be provisioned before a site recovery can be activated. For details on how to provision a SDDC refer to the SDDC documentation.

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

resource "vmc_site_recovery" "site_recovery_1" {
  sddc_id = vmc_sddc.sddc_1.id
  srm_extension_key_suffix = var.site_recovery_srm_extension_key_suffix
}

```

## Argument Reference

The following arguments are supported for vmc_public_ip resource:

* `sddc_id` - (Required) SDDC Identifier.

* `srm_extension_key_suffix` - (Optional) Custom extension key suffix for SRM. If not specified, default extension key will be used. The custom extension suffix must contain 13 characters or less, be composed of letters, numbers, ., - characters only. The extension suffix must begin and end with a letter or number. The suffix is appended to com.vmware.vcDr- to form the full extension key.        
   Length of srm_extension_key_suffix cannot be more than 13 characters.                                   

## Attributes Reference

In addition to arguments listed above, the following attributes are exported after public IP creation:

* `id` - Public IP identifier.

* `ip` - Public IP.

* `display_name` - Display name for public IP.
