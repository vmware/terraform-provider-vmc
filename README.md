# Terraform provider for VMware Cloud on AWS

This is the repository for the Terraform provider for VMware Cloud, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

# Requirements


- [Terraform](https://www.terraform.io/downloads.html) 0.12+
- [Go](https://golang.org/doc/install) 1.16 (to build the provider plugin)


# Building the Provider

The instructions outlined below are specific to Mac OS or Linux OS only.

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements](https://github.com/vmware/terraform-provider-vmc#requirements) before proceeding).

First, you will want to clone the repository to : `$GOPATH/src/github.com/vmware/terraform-provider-vmc`

```sh
mkdir -p $GOPATH/src/github.com/vmware
cd $GOPATH/src/github.com/vmware
git clone git@github.com:vmware/terraform-provider-vmc.git
```

After the clone is complete, you can enter the provider directory and build the provider.

```sh
cd $GOPATH/src/github.com/vmware/terraform-provider-vmc
go get
go build -o terraform-provider-vmc
```

After the build is complete, if your terraform running folder does not match your GOPATH environment, you need to copy the `terraform-provider-vmc` executable to your running folder and re-run `terraform init` to make terraform aware of your local provider executable.


# Using the Provider

To use a released provider in your Terraform environment, run [`terraform init`](https://www.terraform.io/docs/commands/init.html) and Terraform will automatically install the provider. To specify a particular provider version when installing released providers, see the [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#version-provider-versions).

To instead use a custom-built provider in your Terraform environment (e.g. the provider binary from the build instructions above), follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-plugins) After placing the custom-built provider into your plugins directory,  run `terraform init` to initialize it.

For either installation method, documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/vmc/index.html).


## Controlling the provider version

Note that you can also control the provider version. This requires the use of a
`provider` block in your Terraform configuration if you have not added one
already.

The syntax is as follows:

```sh
provider "vmc" {
  version = "~> 1.0"
  ...
}
```

Version locking uses a pessimistic operator, so this version lock would mean
anything within the 1.x namespace, including or after 1.0.0. [Read
more][provider-vc] on provider version control.

[provider-vc]: https://www.terraform.io/docs/configuration/providers.html#provider-versions


# Automated Installation (Recommended)

Download and initialization of Terraform providers is with the “terraform init” command. This applies to the VMC provider as well. Once the provider block for the VMC provider is specified in your .tf file, “terraform init” will detect a need for the provider and download it to your environment.
You can list versions of providers installed in your environment by running “terraform version” command:

```sh
$ terraform version
Terraform v0.12.20
+ provider.vmc (unversioned)
```


# Manual Installation

**NOTE:** Unless you are [developing](#developing-the-provider) or require a
pre-release bugfix or feature, you will want to use the officially released
version of the provider (see [the section above](#using-the-provider)).

**NOTE:** Note that if the provider is manually copied to your running folder (rather than fetched with the “terraform init” based on provider block), Terraform is not aware of the version of the provider you’re running. It will appear as “unversioned”:

```sh
$ terraform version
Terraform v0.12.20
+ provider.vmc (unversioned)
```

Since Terraform has no indication of version, it cannot upgrade in a native way, based on the “version” attribute in provider block.
In addition, this may cause difficulties in housekeeping and issue reporting.


# Developing the Provider

**NOTE:** Before you start work on a feature, please make sure to check the
[issue tracker][gh-issues] and existing [pull requests][gh-prs] to ensure that
work is not being duplicated. For further clarification, you can also ask in a
new issue.

[gh-issues]: https://github.com/vmware/terraform-provider-vmc/issues
[gh-prs]: https://github.com/vmware/terraform-provider-vmc/pulls

See [the section above](#building-the-provider) for details on building the
provider.


# Testing the Provider

Set required environment variables based as per your infrastructure settings

```sh
$ export API_TOKEN=xxx
$ export ORG_ID=xxxx
$ export TEST_SDDC_ID=xxx
$ export TEST_SDDC_NAME=xxx
$ export NSXT_REVERSE_PROXY_URL=xxx
$ export AWS_ACCOUNT_NUMBER=xxx
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

The Terraform provider for VMware Cloud on AWS is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-vmc/blob/master/LICENSE).
