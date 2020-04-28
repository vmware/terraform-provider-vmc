# Example: Public IP management for SDDC

This is an example that demonstrates public IP management actions like creating and deleting public IP after SDDC has been created.

To run the example:

* Generate an API token using [VMware Cloud on AWS console] (https://vmc.vmware.com/console/)

* Update the variables required parameters api_token and org_id in [variables.tf](https://github.com/terraform-providers/terraform-provider-vmc/blob/master/examples/public_ip/variables.tf) with your infrastructure settings. 

* Load the provider

```sh
    terraform init
```

* Execute the plan

```sh
   terraform apply
```

or

```sh
   terraform apply -var="api_token=xxxx" -var="org_id=xxxx"
```

Verify public IP has been created successfully.

* Check the terraform state

```sh
    terraform show
```

* Delete public IP

```sh
    terraform destroy
```

or

```sh
    terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
```
