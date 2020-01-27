---
layout: "vmc"
page_title: "VMC: customer_subnets"
sidebar_current: "docs-vmc-datasource-customer-subnets"
description: A customer subnets data source.
---

# vmc_customer_subnets

The customer subnets data source provides information about customer's compatible subnets for account linking.
## Example Usage

```hcl
data "vmc_customer_subnets" "my_subnets" {
	org_id = "${data.vmc_org.my_org.id}"
	region = "us-west-2"
}
```

## Argument Reference

* `org_id` - (Required) Organization identifier.

* `region` - (Required) The region of the cloud resources to work in.

* `num_hosts` - (Optional) The number of hosts.

* `connected_account_id` - (Optional) The linked connected account identifier.

* `sddc_id` - (Optional) SDDC identifier.

* `force_refresh` - (Optional) Boolean value when set to true, forces the mappings for datacenters to be refreshed for the connected account.

* `instance_type` - (Optional) The server instance type to be used.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `customer_available_zones` - A list of AWS availability zones.

* `ids` - A list of AWS subnet IDs to create links to in the customer's account.