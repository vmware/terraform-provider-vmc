---
layout: "vmc"
page_title: "VMC: connected_accounts"
sidebar_current: "docs-vmc-datasource-connected-accounts"
description: A connected accounts data source.
---

# vmc_connected_accounts

The connected accounts data source get a list of connected accounts.

## Example Usage

```hcl
data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}
```

## Argument Reference

* `org_id` - (Computed) Organization identifier.

* `account_number` - (Optional) AWS account number.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `ids` - The corresponding connected (customer) account UUID this connection is attached to.
