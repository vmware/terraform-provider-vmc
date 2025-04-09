---
page_title: "VMC: vmc_cluster"
description: A resource managing clusters.
---

# Resource: vmc_cluster

Provides a resource to manage clusters.

~> **Note:** Cluster resource implicitly depends on SDDC resource creation. SDDC
must be provisioned before a cluster can be created.

## Example Usage

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
}

resource "vmc_cluster" "Cluster-1" {
  sddc_id   = vmc_sddc.sddc_1.id
  num_hosts = var.num_hosts
  microsoft_licensing_config {
    mssql_licensing   = "DISABLED"
    windows_licensing = "ENABLED"
  }
}
```

## Modifying an Elastic DRS Policy

For a new cluster, elastic DRS uses the **Default Storage Scale-Out** policy,
adding hosts only when storage utilization exceeds the threshold of 75%.

You can select a different policy if it provides better support for your
workload VMs by updating the resource using the following arguments:

* `edrs_policy_type` - (Optional) The EDRS policy type. This can either be
  `cost`, `performance`, `storage-scaleup` or `rapid-scaleup`. Defaults to
  `storage-scaleup`.

* `enable_edrs` - (Optional) Enable EDRS.

* `min_hosts` - (Optional) The minimum number of ESX hosts that the cluster can
  scale in to.

* `max_hosts` - (Optional) The maximum number of ESX hosts that the cluster can
  scale out to.

~> **Note:** When the EDRS policy is disabled (*i.e.*, `enable_edrs = false`)
for `performance`, `cost` or `rapid-scaleup`, the policy type changes to the
default, `storage-scaleup`.

~> **Note:** The EDRS policy properties can be modified only after a cluster has
been created.

## Example

```hcl
provider "vmc" {
  refresh_token = var.api_token
  org_id        = var.org_id
}

data "vmc_connected_accounts" "my_accounts" {
  account_number = var.aws_account_number
}

data "vmc_customer_subnets" "my_subnets" {
  connected_account_id = data.vmc_connected_accounts.my_accounts.id
  region               = var.sddc_region
}
resource "vmc_cluster" "Cluster-1" {
  sddc_id          = vmc_sddc.sddc_1.id
  num_hosts        = var.num_hosts
  edrs_policy_type = "cost"
  enable_edrs      = true
  min_hosts        = 3
  max_hosts        = 8
}
```

## Argument Reference

The following arguments are supported for this resource:

* `sddc_id` - (Required) SDDC identifier.

* `num_hosts` - (Required) Number of ESX hosts in the cluster. The number of ESX
  hosts must be between 2-16 hosts for a cluster.

* `host_cpu_cores_count` - (Optional) Customize CPU cores on ESX hosts in a
  cluster. Specify number of cores to be enabled on ESX hosts in a cluster.

* `host_instance_type` - (Optional) The instance type for the ESX hosts added to
  this cluster. Allowed values include: `I3_METAL`, `I3EN_METAL`, `I4I_METAL`,
  and `R5_METAL`. Defaults to `I3_METAL`.

* `microsoft_licensing_config` - (Optional) Indicates the desired licensing
  support, if any, of Microsoft software.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported:

* `id` - The cluster identifier.

* `cluster_info` - Information about cluster such as name, state, host instance
  type, and cluster identifier.

## Import

Import the using the `id` and `sddc_id`.

`$ terraform import vmc_cluster.cluster_1 id,sddc_id`

For example:

`$ terraform import vmc_cluster.cluster_1 afe7a0fd-3f0a-48b2-9ddb-0489c22732ae,45495963-d24d-469b-830a-9003bfe132b5`
