# Example: Deploy a MultiAZ SDDC for stretched cluster

This is an example that demonstrates MultiAZ SDDC management actions like creating, updating and deleting an existing SDDC.

# Stretched Clusters

The Stretched Clusters feature deploys a single SDDC across two AWS availability zones. 
This option is only available during the SDDC creation. When enabled the default number of ESXi hosts supported in a Stretched Cluster is six. 
Additional hosts can be added later but must be done in pairs across AWS availability zones. The Stretched Clusters feature requires an AWS VPC with two subnets, one subnet per availability zone. The subnets determine an ESXi host placement between the two availability zones.

To run the example:

* Generate an API token using [VMware Cloud on AWS console] (https://vmc.vmware.com/console/)

* Update the required parameters api_token and org_id in [variables.tf](https://github.com/terraform-providers/terraform-provider-vmc/blob/master/examples/sddc/variables.tf) with your infrastructure settings. 

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

Verify the SDDC has been created successfully.

* Check the terraform state

```sh
    terraform show
```

* Delete the SDDC

```sh
    terraform destroy
```

or

```sh
    terraform destroy -var="api_token=xxxx" -var="org_id=xxxx"
```
