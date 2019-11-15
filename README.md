# Terraform VMC Provider

This is the repository for the Terraform VMC Provider, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

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
* id
* sddc_name

Note if you wnat to connect to the staging environment, uncomment the vmc_url and csp_url under vmc settings.

```
provider "vmc" {
  refresh_token = ""

  # for staging environment only
  # vmc_url       = "https://stg.skyscraper.vmware.com/vmc/api"
  # csp_url       = "https://console-stg.cloud.vmware.com"
}

data "vmc_org" "my_org" {
  id = ""
}

data "vmc_connected_accounts" "my_accounts" {
  org_id = "${data.vmc_org.my_org.id}"
}

data "vmc_customer_subnets" "my_subnets" {
  org_id               = "${data.vmc_org.my_org.id}"
  connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
  region               = "us-west-2"
}

resource "vmc_sddc" "sddc_1" {
  org_id = "${data.vmc_org.my_org.id}"

  # storage_capacity    = 100
  sddc_name           = ""
  vpc_cidr            = "10.2.0.0/16"
  num_host            = 1
  provider_type       = "AWS"
  region              = "${data.vmc_customer_subnets.my_subnets.region}"
  vxlan_subnet        = "192.168.1.0/24"
  delay_account_link  = false
  skip_creating_vxlan = false
  sso_domain          = "vmc.local"

  # sddc_template_id = ""
  deployment_type = "SingleAZ"

  account_link_sddc_config = [
    {
      customer_subnet_ids  = ["${data.vmc_customer_subnets.my_subnets.ids.0}"]
      connected_account_id = "${data.vmc_connected_accounts.my_accounts.ids.0}"
    },
  ]
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
$ export ORG_ID=xxx
$ export TEST_SDDC_ID=xxx 
$ make testacc
```