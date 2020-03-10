# Terraform provider for VMware Cloud on AWS

This is the repository for the Terraform provider for VMware Cloud, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

## Build the Provider

The instructions outlined below to build the provider are specific to Mac OS or Linux OS only.

Clone repository to: `$GOPATH/src/github.com/provider/`

```sh
mkdir -p $GOPATH/src/github.com/provider/
cd $GOPATH/src/github.com/provider/
git clone https://github.com/terraform-providers/terraform-provider-vmc.git
```

Enter the provider directory and build the provider

```sh
cd $GOPATH/src/github.com/provider/terraform-provider-vmc
go get
go build -o terraform-provider-vmc
```

Using the Provider
----------------------

The instructions and configuration details to run the provider can be found in [examples/README.md](https://github.com/terraform-providers/terraform-provider-vmc/blob/master/examples/README.md)


Testing the Provider
----------------------

Set required environment variables based as per your infrastructure settings

```sh
$ export API_TOKEN=xxx
$ export ORG_ID=xxxx
$ export TEST_SDDC_ID=xxx
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

If you want to run against a specific set of tests, run make testacc with the TESTARGS parameter containing the run mask as per below:

```sh
$ make testacc TESTARGS="-run=TestAccResourceVmcSddc_basic"
```

# License

Copyright 2019 VMware, Inc.

The Terraform provider for VMware Cloud on AWS is available under [MPL2.0 license](https://github.com/terraform-providers/terraform-provider-vmc/blob/master/LICENSE).
