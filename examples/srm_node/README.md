# Example: SRM instance management for SDDC

This is an example that demonstrates SRM instance management actions like adding and deleting an instance after site recovery has been activated.

To run the example:

* Generate an API token using [VMware Cloud on AWS console] (https://vmc.vmware.com/console/)

* Update the required parameters api_token and org_id in [variables.tf](https://github.com/terraform-providers/terraform-provider-vmc/blob/master/examples/srm_node/variables.tf) with your infrastructure settings. 

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

Verify SRM instance has been added successfully.

* Check the terraform state

```sh
    terraform show
```

* Delete SRM instance

```sh
    terraform destroy
```

or

```sh
    terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
```
