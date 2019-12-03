layout: "vmc"
page_title: "VMC: org"
sidebar_current: "docs-vmc-datasource-org"
description: A organization data source.
---

# vmc_org

This data source provides information about an organization.

## Example Usage

```hcl
data "vmc_org" "my_org" {
	id = ""

}
```

## Argument Reference

* `id` - (Required) ID of the organization