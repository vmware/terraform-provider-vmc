# Terraform VMC Provider

This is the repository for the Terraform VMC Provider, which one can use with
Terraform to work with (VMware Cloud on AWS)[https://vmc.vmware.com/].

This provider is currently implemented using go bindings generated from swagger codegen.   
Please note the binding may change before the public release.

# Use the provider

## Requirements

* Install [Terraform 0.10.1+](https://learn.hashicorp.com/terraform/getting-started/install.html)
* Install [Go 1.9](https://golang.org/doc/install) (to build the provider plugin)
* Set GOPATH to $HOME/go. For detail see [here](https://github.com/golang/go/wiki/SettingGOPATH)

## Build the Provider

Clone repository to: `$GOPATH/src/gitlab.eng.vmware.com/vapi-sdk/`

```sh
mkdir -p $GOPATH/src/gitlab.eng.vmware.com/vapi-sdk
cd $GOPATH/src/gitlab.eng.vmware.com/vapi-sdk
git clone https://gitlab.eng.vmware.com/vapi-sdk/terraform-provider-vmc.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/gitlab.eng.vmware.com/vapi-sdk/terraform-provider-vmc
go get
go build -o terraform-provider-vmc
```

## Load the provider
```sh
terraform init
```

## Connect to VMC and create a testing sddc

Update following fields in the [main.tf](main.tf) with your infra settings

* refresh_token
* csp_url
* id
* sddc_name

```
provider "vmc" {
  refresh_token = ""
  csp_url       = "https://console-stg.cloud.vmware.com"
  # for production 
  # cap_url     = "https://console.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  # org id listed in UI. e.g: 05e0a625-3293-41bb-a01f-35e762781c2a
  id = ""
}

resource "vmc_sddc" "sddc_1" {
  org_id        = "${data.vmc_org.my_org.id}"
  sddc_name     = ""
  num_host      = 4
  provider_type = "ZEROCLOUD"
  region        = "US_WEST_1"
  vpc_cidr      = "10.0.0.0/17"
}
```

## Try a dry run

```sh
terraform plan
```

Check if the terraform plan looks good

## Execute the plan

```sh
terraform apply
```

Verified the sddc is created

## Add/Remove hosts

Update the "num_host" field in [main.tf](main.tf) to expected number.   
Review and execute the plan

```sh
terraform plan
terraform apply
```

Verified the hosts are added/removed successfully.

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