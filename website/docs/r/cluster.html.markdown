---
layout: "vmc"

page_title: "VMC: vmc_cluster"
sidebar_current: "docs-vmc-resource-cluster"

description: |-
  Provides a resource to manage clusters.
---

# vmc_cluster

Provides a resource to manage clusters.
~> **Note:** Cluster resource implicitly depends on SDDC resource creation. SDDC must be provisioned before a cluster can be created. For details on how to provision a SDDC refer to [vmc_sddc](https://www.terraform.io/docs/providers/vmc/r/sddc.html).

## Example Usage

```hcl

provider "vmc" {
  refresh_token = var.api_token
  org_id = var.org_id
}

 resource "vmc_cluster" "Cluster-1"{
    sddc_id = vmc_sddc.sddc_1.id
    num_hosts= var.num_hosts
}


```

## Argument Reference

The following arguments are supported for vmc_cluster resource:

* `sddc_id` - (Required) SDDC identifier.

* `num_hosts` - (Required) Number of hosts in the cluster. The number of hosts must be between 3 - 16 hosts for a cluster.

* `host_cpu_cores_count` - (Optional) Customize CPU cores on hosts in a cluster. Specify number of cores to be enabled on hosts in a cluster.

* `host_instance_type` - (Optional) The instance type for the esx hosts added to this cluster. Possible values are: I3_METAL, I3EN_METAL and R5_METAL. Default value: I3_METAL.

* `storage_capacity` - (Optional) For EBS-backed instances i.e: for R5_METAL only, the requested storage capacity in GiB.

## Attributes Reference

In addition to arguments listed above, the following attributes are exported after cluster creation:

* `id` - Cluster identifier.

* `cluster_info` - Information about cluster like name, state, host instance type and cluster identifier.

## Import

Cluster resource can be imported using the `id` and `sddc_id`, e.g.

`$ terraform import vmc_cluster.cluster_1 id,sddc_id`

- id = Cluster Identifier
- sddc_id = SDDC Identifier

`$ terraform import vmc_cluster.cluster_1 afe7a0fd-3f0a-48b2-9ddb-0489c22732ae,45495963-d24d-469b-830a-9003bfe132b5`