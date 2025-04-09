# Example

## Creation of SDDC, Public IP, Cluster, Site Recovery, and SRM Node

This is an example that demonstrates the creation of VMC resources like SDDC,
cluster, public IP, Site Recovery and SRM node.

For Site Recovery activation, a 10-minute delay must be added after SDDC is
created and before Site Recovery can be activated.

To add delay after SDDC has been created, update SDDC resource in the
[`main.tf`](https://github.com/vmware/terraform-provider-vmc/blob/main/examples/main.tf)
with `local-exec` provisioner:

```hcl
resource "vmc_sddc" "sddc_1" {
  sddc_name           = var.sddc_name
  vpc_cidr            = var.vpc_cidr
  num_host            = var.num_hosts
  provider_type       = var.provider_type
  region              = data.vmc_customer_subnets.my_subnets.region
  vxlan_subnet        = var.vxlan_subnet
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"
  sddc_type           = var.sddc_type
  deployment_type     = "SingleAZ"

  host_instance_type = var.host_instance_type

  account_link_sddc_config {
    customer_subnet_ids  = [data.vmc_customer_subnets.my_subnets.ids[0]]
    connected_account_id = data.vmc_connected_accounts.my_accounts.id
  }

  timeouts {
    create = "300m"
    update = "300m"
    delete = "180m"
  }

  # provisioner defined to add 10 minute delay after SDDC creation to enable site recovery activation.
  provisioner "local-exec" {
    command = "sleep 600"
  }
}
```

To run the example:

* Generate an API token using [VMware Cloud on AWS console]
  (<https://vmc.vmware.com/console/>)

* Update the required parameters `api_token` and `org_id` in the
  [`variables.tf`](https://github.com/vmware/terraform-provider-vmc/blob/main/examples/variables.tf)
  with your infrastructure settings.

* Load the provider:

  ```sh
  terraform init
  ```

* Run the plan:

  ```sh
  terraform apply
  ```

  or

  ```sh
  terraform apply -var="api_token=xxxx" -var="org_id=xxxx"
  ```

* Check the state:

  ```sh
  terraform show
  ```

* Delete VMC resources created during apply.

  ```sh
  terraform destroy
  ```

  or

  ```sh
  terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
  ```
