* go get gitlab.eng.vmware.com/het/vmc-go-sdk
* Upddate refresh_token field in main.tf
* go build -o terraform-provider-vmc
* terraform init
* terraform plan
* terraform apply 


# Terraform VMC Provider

This is the repository for the Terraform VMC Provider, which one can use with
Terraform to work with [VMware Cloud on AWS][https://vmc.vmware.com/].

This provider is currently implemented using go bindings generated from swagger codegen.   
Please note the binding may change before the public release.

# Using the Provider

git clone https://gitlab.eng.vmware.com/het/terraform-provider-vmc.git
cd terraform-provider-vmc
go get ...
fill in your refresh_token in main.ft
terraform init

## To create a testing sddc
terraform plan
terraform apply

## To delete the sddc
terraform destroy

# Testing the Provider

## Set required environment variable
export REFRESH_TOKEN=xxx

```sh
$ make testacc
```