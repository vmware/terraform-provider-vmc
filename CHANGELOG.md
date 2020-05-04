## 1.2.1 (May 04, 2020)

BUG FIXES: 

* Added instructions for delay needed after SDDC creation for site recovery [\#21](https://github.com/terraform-providers/terraform-provider-vmc/pull/21)
* Removed capitalized error messages from code [\#23](https://github.com/terraform-providers/terraform-provider-vmc/pull/23)
* Updated module name in go.mod [\#24](https://github.com/terraform-providers/terraform-provider-vmc/pull/24)

ENHANCEMENTS:

* Updated dependencies version to latest in go.mod [\#20](https://github.com/terraform-providers/terraform-provider-vmc/pull/20)
* Added sample .tf file for each resource in examples folder [\#22](https://github.com/terraform-providers/terraform-provider-vmc/pull/22)

## 1.2.0 (April 03, 2020)

FEATURES:

* **New Resource:** `vmc_site_recovery` Added resource for site recovery management. [\#14](https://github.com/terraform-providers/terraform-provider-vmc/pull/14)
* **New Resource:** `vmc_srm_node` Added resource to add SRM instance after site recovery has been activated. [\#14](https://github.com/terraform-providers/terraform-provider-vmc/pull/14)


## 1.1.1 (March 24, 2020)

BUG FIXES:

* Set ForceNew for host_instance_type to true in order to enforce SDDC redeploy when host_instance_type is changed [\#5](https://github.com/terraform-providers/terraform-provider-vmc/pull/5)
* Fix for re-creating public IP if it is accidentally delete via console. [\#11](https://github.com/terraform-providers/terraform-provider-vmc/pull/11)
* Updated variables.tf for description of fields. [\#10](https://github.com/terraform-providers/terraform-provider-vmc/pull/10)


## 1.1.0 (March 10, 2020)

FEATURES:

* **New Resource:** `vmc_sddc`
* **New Resource:** `vmc_public_ip` [\#43](https://github.com/vmware/terraform-provider-vmc/pull/43)
* **New Data Source:** `vmc_org`
* **New Data Source:** `vmc_connected_accounts`
* **New Data Source:** `vmc_customer_subnets`


ENHANCEMENTS:

* vmc_sddc: Added nsxt_reverse_proxy_url to SDDC resource data. [\#23](https://github.com/vmware/terraform-provider-vmc/pull/23)
* vmc_connected_accounts: Added filtering to return AWS account ID associated with the account number provided in configuration. [\#30](https://github.com/vmware/terraform-provider-vmc/pull/30)
* provider.go: Added org_id as a required parameter in terraform schema. [\#38](https://github.com/vmware/terraform-provider-vmc/pull/38)
* data_source_vmc_customer_subnets.go : Added validateFunctions for sddc and customer subnet resources. [\#41](https://github.com/vmware/terraform-provider-vmc/pull/41)
* examples/main.tf : Added expression to convert AWS specific region to VMC region. [\#46](https://github.com/vmware/terraform-provider-vmc/pull/46) 


BUG FIXES:

* Moved main.tf to examples folder. [\#17](https://github.com/vmware/terraform-provider-vmc/pull/17)
* License statement fixed. [\#8](https://github.com/vmware/terraform-provider-vmc/pull/8)
* Implemented connected account data source to return a single account ID associated with the account number. [\#40](https://github.com/vmware/terraform-provider-vmc/pull/40)
* Set ForceNew for host_instance_type to true in order to enforce SDDC redeploy when host_instance_type is changed. [\#5](https://github.com/terraform-providers/terraform-provider-vmc/pull/5)
