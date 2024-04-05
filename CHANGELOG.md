## 1.15.1 (Apr 5, 2024)

BUG FIXES:
* bump protobuf due to a security vulnerability

## 1.15 (Jan 25, 2024)

ENHANCEMENT:
* bumping vsphere-automation-sdk-go/services/vmc to v0.14.0
* Added support for disk less instances c6i/m7i [\#214](https://github.com/vmware/terraform-provider-vmc/pull/214)

## 1.14

Breaking changes:
* Removing support for r5.metal instances [\#205](https://github.com/vmware/terraform-provider-vmc/pull/205)
* removing `storage_capacity` field from sddc and cluster resources schema [\#205](https://github.com/vmware/terraform-provider-vmc/pull/205)

BUG FIXES:

* Failure during SDDC group creation [\#208](https://github.com/vmware/terraform-provider-vmc/pull/210)

## 1.13.3 (Aug 21, 2023)

BUG FIXES:

* Fixing errors when creating and deleting multiple srm nodes [\#175](https://github.com/vmware/terraform-provider-vmc/pull/195)

## 1.13.2 (Aug 8, 2023)

BUG FIXES:

* Fixing errors when reading Customer Subnets [\#191](https://github.com/vmware/terraform-provider-vmc/pull/191)


## 1.13.1 (Aug 4, 2023)

BUG FIXES:

* Allowing usage of microsoft_license_config upon SDDC creation. Reading microsoft_license_config.academic_license field [\#190](https://github.com/vmware/terraform-provider-vmc/pull/190)

ENHANCEMENT:
* Bump google.golang.org/grpc from 1.51.0 to 1.53.0 [\#184](https://github.com/vmware/terraform-provider-vmc/pull/184)
* Updates to documentation [\#186](https://github.com/vmware/terraform-provider-vmc/pull/186)

## 1.13 (Feb 23, 2023)

FEATURES:

* Added support for OAuth2.0 app authentication [\#173](https://github.com/vmware/terraform-provider-vmc/pull/173)

Fixes for security vulnerabilities.

## 1.12.1 (Feb 3, 2023)

BUG FIXES:

* Destroying sddc_group times out and then fails on subsequent attempts [\#172](https://github.com/vmware/terraform-provider-vmc/pull/172)
* Remove restrictions on 6 hosts minimum in MultiAZ SDDCs  [\#171](https://github.com/vmware/terraform-provider-vmc/pull/171)

## 1.12.0 (Nov 9, 2022)

FEATURES:

* Added support for SDDC Groups

## 1.11.0 (Oct 17, 2022)

FEATURES:

* `num_hosts` property on SDDC now shows and controls the number of hosts in the primary cluster (created by default with an SDDC). Previously there was no way to scale up/down the number of hosts in the primary cluster
* Cluster operations like Create/Update/Destroy on more than one cluster can now be initiated simultaneously, without the need for depends_on=[] 
* Error reporting improvements

BUG FIXES:

* EDRS settings are not truly "Optional" [\#151](https://github.com/vmware/terraform-provider-vmc/pull/151)
* Lack of multi-cluster SDDC support in "resourceSddcUpdate" function [\#155](https://github.com/vmware/terraform-provider-vmc/pull/155)
* vmc_cluster resource tries to create new clusters simultaneously and fails [\#160](https://github.com/vmware/terraform-provider-vmc/pull/160)

## 1.10.1 (Sep 20, 2022)

BUG FIXES:
Defaults for enable_edrs, edrs_policy_type, max_hosts, min_hosts should not be set to null https://github.com/vmware/terraform-provider-vmc/issues/94

## 1.10.0 (Jul 12, 2022)

FEATURES:

* `vmc_sddc` Added I4I_METAL host instance type support.

BUG FIXES:

* Allow min_hosts as low as 2 as per VMC service backend default value [\#147](https://github.com/vmware/terraform-provider-vmc/pull/147)
* Removed references to R5 host instance type in examples as new deployments are unavailable [\#146](https://github.com/vmware/terraform-provider-vmc/pull/146)

## 1.9.3 (Jun 7, 2022)

BUG FIXES:
* nil derreference when doing "terraform plan" in some environments (https://github.com/vmware/terraform-provider-vmc/pull/142)

## 1.9.2 (Jun 7, 2022)

BUG FIXES:
* nil derreference when doing "terraform plan" in some environments (https://github.com/vmware/terraform-provider-vmc/pull/141)

ENHANCEMENT:

* Upgrade to TF plugin SDK v2.11.0 due to CVE-2022-30323 (https://github.com/vmware/terraform-provider-vmc/pull/140)

## 1.9.1 (Mar 28, 2022)

BUG FIXES:
* Example files in /examples/.. dir now have a required_providers declaration, fixing "terraform init"

## 1.9.0 (Mar 16, 2022)

ENHANCEMENT:
* Upgrade Provider to use VMC SDK 0.8.0 [\#130](https://github.com/vmware/terraform-provider-vmc/pull/130)

## 1.8.0 (Nov 10, 2021)

ENHANCEMENT:
 * Upgrade Provider to use VMC SDK 0.6.0 [\#116](https://github.com/vmware/terraform-provider-vmc/pull/116), [\#118](https://github.com/vmware/terraform-provider-vmc/pull/118)
 * Direct access to NSX Manager[\#116](https://github.com/vmware/terraform-provider-vmc/pull/116)

BUG FIXES:
 * Property "esx_hosts" in SddcResourceConfig is deprecated. [\#119](https://github.com/vmware/terraform-provider-vmc/pull/119)

## 1.7.0 (Aug 5, 2021)

ENHANCEMENT:

 * Upgrade to TF plugin SDK v2 [\#108](https://github.com/vmware/terraform-provider-vmc/pull/108)

BUG FIXES:

 * Fix for static check failures caused by references to deprecated functions [\#109](https://github.com/vmware/terraform-provider-vmc/pull/109)

## 1.6.0 (May 17, 2021)

FEATURES:

 * Support for SDDC data source [\#105](https://github.com/vmware/terraform-provider-vmc/pull/105)
 * Support for M1 mac [\#107](https://github.com/vmware/terraform-provider-vmc/pull/107)

 BUG FIXES:

 * Fix for updating multiple params in SDDC resource [\#101](https://github.com/vmware/terraform-provider-vmc/pull/101)

## 1.5.1 (January 21, 2021)

 BUG FIXES:

 * Moved subnet validation from customized diff to create SDDC block [\#96](https://github.com/vmware/terraform-provider-vmc/pull/96)
 * Added zerocloud check for setting intranet MTU uplink [\#98](https://github.com/vmware/terraform-provider-vmc/pull/98)
 * Updated help documentation to specify current limitations in import SDDC functionality [\#97](https://github.com/vmware/terraform-provider-vmc/pull/97)

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
