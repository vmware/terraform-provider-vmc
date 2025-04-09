---
page_title: "VMC: vmc_srm_node"
description: A resource for adding an instance to SDDC after site recovery has been activated.
---

# Resource:  vmc_srm_node

Provides a resource to add an instance to SDDC after site recovery has been
activated.

~> **Note:** This resource depends on site recovery resource creation. Site
recovery must be activated to add this resource.

## Example Usage

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
}

resource "vmc_srm_node" "srm_node_1" {
  sddc_id                       = vmc_sddc.sddc_1.id
  srm_node_extension_key_suffix = var.srm_node_srm_extension_key_suffix
  depends_on                    = [vmc_site_recovery.site_recovery_1]
}
```

## Argument Reference

The following arguments are supported for this resource:

* `sddc_id` - (Required) The SDDC identifier.

* `srm_node_extension_key_suffix` - (Required) Custom extension key suffix for
  SRM. If not specified, default extension key will be used. The custom
  extension suffix must contain 13 characters or fewer, be composed of letters,
  numbers, `.`, and `-` characters. The extension suffix must begin and end with
  a letter or number. The suffix is appended to `com.vmware.vcdr-` to form the
  full extension key.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `site_recovery_state` - The site recovery state. Allowed values include:
  `ACTIVATED`, `ACTIVATING`, `CANCELED`, `DEACTIVATED`, `DEACTIVATING`,
  `DELETED`, and `FAILED`.

* `srm_node` - The site Recovery node information.

* `vr_node` - The vSphere Replication node information.

## Import

Import the resource using the `id` and `sddc_id`.

`$ terraform import vmc_srm_node.srm_node_1 id,sddc_id`

For example:

`$ terraform import vmc_srm_node.srm_node_1 7aad97e9-9a4f-4e43-8817-5c8d8c0e87a5,afe7a0fd-3f0a-48b2-9ddb-0489c22732ae`
