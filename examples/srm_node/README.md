# Example

## SRM Instance Management for SDDC

This is an example that demonstrates SRM instance management actions like adding
and deleting an instance after Site Recovery has been activated.

To run the example:

* Generate an API token using
  [VMware Cloud on AWS console](https://vmc.vmware.com/console/).

* Update the required parameters `api_token` and `org_id` in the
  [`variables.tf`](https://github.com/vmware/terraform-provider-vmc/blob/main/examples/srm_node/variables.tf)
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

  Verify SRM instance has been added successfully.

* Check the state:

  ```sh
  terraform show
  ```

* Destroy:

  ```sh
  terraform destroy
  ```

  or

  ```sh
  terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
  ```
