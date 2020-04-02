---
layout: "vmc"

page_title: "VMC: vmc_site_recovery"
sidebar_current: "docs-vmc_site_recovery"

description: |-
  Provides a resource to activate and deactivate site recovery for an SDDC.
---

# vmc_site_recovery

Provides a resource to activate and deactivate site recovery for an SDDC.
~> **Note:** Site recovery resource implicitly depends on SSDC resource creation. SDDC must be provisioned before a site recovery can be activated. For details on how to provision an SDDC refer to [vmc_sddc](https://www.terraform.io/docs/providers/vmc/r/sddc.html).

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

The following arguments are supported for vmc_site_recovery resource:

* `sddc_id` - (Required) SDDC identifier.

* `srm_node_extension_key_suffix` - (Optional) Custom extension key suffix for SRM. If not specified, default extension key will be used. 
The custom extension suffix must contain 13 characters or less, be composed of letters, numbers, ., - characters. The extension suffix must begin and end with a letter or number. The suffix is appended to com.vmware.vcDr- to form the full extension key.


## Attributes Reference

In addition to arguments listed above, the following attributes are exported after site recovery activation:

* `site_recovery_state` - Site recovery state. Possible values are: ACTIVATED, ACTIVATING, CANCELED, DEACTIVATED, DEACTIVATING, DELETED, FAILED.

* `srm_node` - Site recovery node created after site recovery activation.

* `vr_node` - VR node created after site recovery activation.
