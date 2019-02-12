# Terraform VMC Provider

This is the repository for the Terraform VMC Provider, which one can use with
Terraform to work with (VMware Cloud on AWS)[https://vmc.vmware.com/].

This provider is currently implemented using go bindings generated from swagger codegen.   
Please note the binding may change before the public release.

# Using the Provider

```sh
git clone https://gitlab.eng.vmware.com/het/terraform-provider-vmc.git
cd terraform-provider-vmc
go get ...
```

fill in your refresh_token in main.ft

```sh
terraform init
```

## To create a testing sddc

```sh
terraform plan
terraform apply
```

## To delete the sddc

```sh
terraform destroy
```

# Testing the Provider

## Set required environment variable

```sh
$ export REFRESH_TOKEN=xxx
$ make testacc
```