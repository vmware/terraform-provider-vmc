# Changelog

## [1.0.0 (Unreleased)]

(https://github.com/vmware/terraform-provider-vmc/tree/HEAD)

[Full Changelog](https://github.com/vmware/terraform-provider-vmc/compare/87cfcf8f21600ef6198389c569d1e988ca30b5a9...HEAD)

FEATURES:

* resource/vmc_sddc: Support for CRUD operations on Software Defined Data Center on VMC.
* data_source/vmc_org: Support to retrieve details about the organization.
* data_source/vmc_connected_accounts: Support to retrieve a list of connected accounts in the given organization.
* data_source/vmc_customer_subnets: Support to retrieve a customer's compatible subnets for account linking.


ENHANCEMENTS:

* resource/vmc_sddc: Added nsxt_reverse_proxy_url to SDDC resource data. [\#23](https://github.com/vmware/terraform-provider-vmc/pull/23)
* data_source/vmc_connected_accounts: Added filtering to return AWS account ID associated with the account number provided in configuration. [\#30](https://github.com/vmware/terraform-provider-vmc/pull/30)
* resource/vmc_public_ip: Added public IP resource. [\#43](https://github.com/vmware/terraform-provider-vmc/pull/43)

BUG FIXES:

* Moved main.tf to examples folder. [\#17](https://github.com/vmware/terraform-provider-vmc/pull/17)
* License statement fixed. [\#8](https://github.com/vmware/terraform-provider-vmc/pull/8)
