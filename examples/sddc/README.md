# Example

## Provision an SDDC

This is an example that demonstrates SDDC management actions like creating,
updating, and deleting an existing SDDC.

To run the example:

* Generate an API token using
  [VMware Cloud on AWS console](https://vmc.vmware.com/console/).

* Update the required parameters `api_token` and `org_id` in the
  [`variables.tf`](https://github.com/vmware/terraform-provider-vmc/blob/main/examples/sddc/variables.tf)
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

  Verify the SDDC has been created successfully.

* Check the state:

  ```sh
  terraform show
  ```

* Delete the SDDC:

  ```sh
  terraform destroy
  ```

  or

  ```sh
  terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
  ```
