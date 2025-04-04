---
page_title: "VMC: vmc_connected_accounts"
description: The data source for connected accounts.
---

# Data Source: vmc_connected_accounts

The connected accounts data source retrieves a list of connected accounts.

## Example Usage

```hcl
data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}
```

## Argument Reference

* `org_id` - (Computed) The organization identifier.

* `account_number` - (Required) The AWS account number.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The corresponding connected (customer) account UUID this connection is
  attached to.
