# Example

## Site Recovery Management for SDDC

This is an example that demonstrates Site Recovery management actions like
activation and deactivation after SDDC has been created.

To run the example:

* Generate an API token using
  [VMware Cloud on AWS console](https://vmc.vmware.com/console/).

* Update the required parameters `api_token` and `org_id` in the
  [`variables.tf`](https://github.com/vmware/terraform-provider-vmc/blob/main/examples/site_recovery/variables.tf)
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

  Verify Site Recovery has been activated successfully.

* Check the state:

  ```sh
  terraform show
  ```

* Deactivate Site Recovery:

  ```sh
  terraform destroy
  ```

  or

  ```sh
  terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
  ```
