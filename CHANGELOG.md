## 1.5.1 (January 21, 2021)

 BUG FIXES:

 * Moved subnet validation from customized diff to create SDDC block [\#96](https://github.com/vmware/terraform-provider-vmc/pull/96)
 * Added zerocloud check for setting intranet MTU uplink [\#98](https://github.com/vmware/terraform-provider-vmc/pull/98)
 * Updated help documentation for import SDDC functionality [\#97](https://github.com/vmware/terraform-provider-vmc/pull/97)

## 1.5.0 (January 5, 2021)

FEATURES:

 * `vmc_sddc`, `vmc_cluster` Added microsoft licensing configuration to SDDC and cluster resource [\#71](https://github.com/vmware/terraform-provider-vmc/pull/71)
 * Added intranet_uplink_mtu field in `vmc_sddc` resource [\#88](https://github.com/vmware/terraform-provider-vmc/pull/88)

 BUG FIXES:

 * Removed validation on sddc_type field in order to allow empty value. [\#72](https://github.com/vmware/terraform-provider-vmc/pull/72)
 * Removed default values to fix EDRS configuration error for 1NODE SDDC [\#83](https://github.com/vmware/terraform-provider-vmc/pull/83)
 * Added check to store vxlan_subnet information in terraform state file only when skip_creating_vxlan = false [\#86](https://github.com/vmware/terraform-provider-vmc/pull/86)

## 1.4.0 (October 12, 2020)

FEATURES:

* `vmc_sddc` Added I3EN_METAL host instance type support.  [\#42](https://github.com/vmware/terraform-provider-vmc/pull/42)
* `vmc_sddc` Modified code to enable EDRS policy configuration. [\#43](https://github.com/vmware/terraform-provider-vmc/pull/43)
* `vmc_sddc` Added size parameter in resource schema to enable users to deploy large SDDC. [\#59](https://github.com/vmware/terraform-provider-vmc/pull/59)

BUG FIXES: 

* Added check in resourceClusterRead to see if cluster exists and remove cluster information from terraform state file. [\#48](https://github.com/vmware/terraform-provider-vmc/pull/48)
* Added validation check for customer subnet IDs based on the deployment type. [\#54](https://github.com/vmware/terraform-provider-vmc/pull/54)

ENHANCEMENTS:

* Modified Importer State in `vmc_cluster` and `vmc_public_ip` resources for terraform import command. [\#49](https://github.com/vmware/terraform-provider-vmc/pull/49)

## 1.3.0 (June 18, 2020)

FEATURES:

* **New Resource:** `vmc_cluster` Added resource for cluster management. [\#25](https://github.com/vmware/terraform-provider-vmc/pull/25)

BUG FIXES: 

* Modified code to store num_host in resourceSddcRead method [\#39](https://github.com/vmware/terraform-provider-vmc/pull/39)

ENHANCEMENTS:

* Validation check added for MultiAZ SDDC [\#35](https://github.com/vmware/terraform-provider-vmc/pull/29)
* Added detailed error handler functions for CRUD operations on resources and data sources [\#35](https://github.com/vmware/terraform-provider-vmc/pull/29)
* Added documentation for vmc_cluster resource  [\#26](https://github.com/vmware/terraform-provider-vmc/pull/26)

## 1.2.1 (May 04, 2020)

BUG FIXES: 

* Added instructions for delay needed after SDDC creation for site recovery [\#21](https://github.com/vmware/terraform-provider-vmc/pull/21)
* Removed capitalized error messages from code [\#23](https://github.com/vmware/terraform-provider-vmc/pull/23)
* Updated module name in go.mod [\#24](https://github.com/vmware/terraform-provider-vmc/pull/24)

ENHANCEMENTS:

* Updated dependencies version to latest in go.mod [\#20](https://github.com/vmware/terraform-provider-vmc/pull/20)
* Added sample .tf file for each resource in examples folder [\#22](https://github.com/vmware/terraform-provider-vmc/pull/22)

## 1.2.0 (April 03, 2020)

FEATURES:

* **New Resource:** `vmc_site_recovery` Added resource for site recovery management. [\#14](https://github.com/vmware/terraform-provider-vmc/pull/14)
* **New Resource:** `vmc_srm_node` Added resource to add SRM instance after site recovery has been activated. [\#14](https://github.com/vmware/terraform-provider-vmc/pull/14)


## 1.1.1 (March 24, 2020)

BUG FIXES:

* Set ForceNew for host_instance_type to true in order to enforce SDDC redeploy when host_instance_type is changed [\#5](https://github.com/vmware/terraform-provider-vmc/pull/5)
* Fix for re-creating public IP if it is accidentally delete via console. [\#11](https://github.com/vmware/terraform-provider-vmc/pull/11)
* Updated variables.tf for description of fields. [\#10](https://github.com/vmware/terraform-provider-vmc/pull/10)


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
* Set ForceNew for host_instance_type to true in order to enforce SDDC redeploy when host_instance_type is changed. [\#5](https://github.com/vmware/terraform-provider-vmc/pull/5)
