# Terraform provider for VMware Cloud on AWS

This is the repository for the Terraform provider for VMware Cloud, which one can use with
Terraform to work with [VMware Cloud on AWS](https://vmc.vmware.com/).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10+
- [Go](https://golang.org/doc/install) 1.13 (to build the provider plugin)

Developing the Provider
---------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (please check the [requirements] before proceeding).

*Note:* This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH). The instructions that follow assume a directory in your home directory outside of the standard GOPATH (i.e `$HOME/development/terraform-providers/`).

Clone repository to: `$HOME/development/terraform-providers/`

```sh
$ mkdir -p $HOME/development/terraform-providers/; cd $HOME/development/terraform-providers/
$ git clone git@github.com:vmware/terraform-provider-vmc.git
...
```

Enter the provider directory and build the provider

```sh
cd $HOME/development/terraform-providers/terraform-provider-vmc
go get
go build -o terraform-provider-vmc
```

Testing the Provider
---------------------------

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run. 

```sh
$ make testacc
```

# License 

Copyright 2019 VMware, Inc.

The Terraform provider for VMware Cloud on AWS is available under [MPL2.0 license](https://github.com/vmware/terraform-provider-vmc/blob/master/LICENSE.txt).
