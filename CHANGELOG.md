# Release History

## 1.15.4

> Release Date: 2025-04-21

CHORE:

- Updated `go` from 1.23.7 to v1.23.8. [#313](https://github.com/vmware/terraform-provider-vmc/pull/313)
- Updated `github.com/hashicorp/terraform-plugin-sdk/v2` from 2.36.0 to 2.36.1. [#297](https://github.com/vmware/terraform-provider-vmc/pull/297)
- Updated `golang.org/x/oauth2` from 0.26.0 to 0.23.0. [#300](https://github.com/vmware/terraform-provider-vmc/pull/300), [#307](https://github.com/vmware/terraform-provider-vmc/pull/307)
- Updated `golang.org/x/net` from 0.34.0 to 0.36.0. [#301](https://github.com/vmware/terraform-provider-vmc/pull/301)
- Updated `github.com/vmware/vsphere-automation-sdk-go` from 0.7.0 to 0.8.0. [#303](https://github.com/vmware/terraform-provider-vmc/pull/303)
- Updated `github.com/golang-jwt/jwt/v4` from 4.5.1 to 4.5.2. [#304](https://github.com/vmware/terraform-provider-vmc/pull/304)
- Updated `github.com/gofrs/uuid/v5` from 5.3.1 to 5.3.2. [#305](https://github.com/vmware/terraform-provider-vmc/pull/305)
- Updated `goreleaser` configuration to v2. [#309](https://github.com/vmware/terraform-provider-vmc/pull/309)
- Updated `golangci-lint` configuration to v2. [#317](https://github.com/vmware/terraform-provider-vmc/pull/317)
- Updated GitHub Actions workflows. [#310](https://github.com/vmware/terraform-provider-vmc/pull/310), [#316](https://github.com/vmware/terraform-provider-vmc/pull/316)
- Migrated documentation from legacy path. [#306](https://github.com/vmware/terraform-provider-vmc/pull/306)
- Applied documentation HCL formatting. [#314](https://github.com/vmware/terraform-provider-vmc/pull/314)
- Removed unused files. [#308](https://github.com/vmware/terraform-provider-vmc/pull/308), [#311](https://github.com/vmware/terraform-provider-vmc/pull/311)

## 1.15.3

> Release Date: 2025-02-11

CHORE:

- Updated `github.com/gofrs/uuid/v5` from 5.3.0 to 5.3.1. [#295](https://github.com/vmware/terraform-provider-vmc/pull/295)
- Updated `github.com/hashicorp/terraform-plugin-sdk/v2` from 2.35.0 to 2.36.0. [#293](https://github.com/vmware/terraform-provider-vmc/pull/293)
- Updated `golang.org/x/net` from 0.31.0 to 0.33.0. [#292](https://github.com/vmware/terraform-provider-vmc/pull/292)
- Updated `golang.org/x/oauth2` from 0.24.0 to 0.25.0 [#291](https://github.com/vmware/terraform-provider-vmc/pull/291)
- Updated `golang.org/x/crypto` from 0.29.0 to 0.31.0. [#290](https://github.com/vmware/terraform-provider-vmc/pull/290)

## 1.15.2

> Release Date: 2024-12-09

DOCUMENTATION:

- Added install, build, and test documentation. [#265](https://github.com/vmware/terraform-provider-vmc/pull/265)

CHORE:

- Updated copyright and SPDX. [#246](https://github.com/vmware/terraform-provider-vmc/pull/246)
- Update `NOTICE`. [#241](https://github.com/vmware/terraform-provider-vmc/pull/241)
- Updated Code of Conduct. [#242](https://github.com/vmware/terraform-provider-vmc/pull/242)
- Removed superfluous files. [#260](https://github.com/vmware/terraform-provider-vmc/pull/260)

- Go:
    - Updated Go to v1.22.7. [#284](https://github.com/vmware/terraform-provider-vmc/pull/284)

- Library Dependencies:
    - Updated `github.com/stretchr/testify` from 1.10.0. [#283](https://github.com/vmware/terraform-provider-vmc/pull/283)
    - Updated `github.com/hashicorp/terraform-plugin-sdk/v2` to 2.34.0. [#238](https://github.com/vmware/terraform-provider-vmc/pull/238)
      - Updated `github.com/golang-jwt/jwt/v4`  to 4.5.1. [#279](https://github.com/vmware/terraform-provider-vmc/pull/279)
    - Updated `github.com/gofrs/uuid/v5` to 5.3.0. [#257](https://github.com/vmware/terraform-provider-vmc/pull/257)

- Specific Packages:
    - Updated `golang.org/x/oauth2` to 0.24.0. [#280](https://github.com/vmware/terraform-provider-vmc/pull/280)
    - Updated `golang.org/x/net` to 0.23.0. [#233](https://github.com/vmware/terraform-provider-vmc/pull/233)

- SDKs:
    - Updated `github.com/vmware/vsphere-automation-sdk-go/services/vmc/draas` to 0.7.0. [#220](https://github.com/vmware/terraform-provider-vmc/pull/220)
    - Updated `github.com/vmware/vsphere-automation-sdk-go/services/nsxt-vmc-aws-integration` to 0.8.0. [#228](https://github.com/vmware/terraform-provider-vmc/pull/228)

- Linting:
    - Fixed grammar. [#286](https://github.com/vmware/terraform-provider-vmc/pull/286)
    - Fixed variable naming. [#268](https://github.com/vmware/terraform-provider-vmc/pull/268)
    - Fixed non-constant format strings. [#267](https://github.com/vmware/terraform-provider-vmc/pull/267)
    - Fixed staticcheck ST1005. [#273](https://github.com/vmware/terraform-provider-vmc/pull/273)
    - Fixed nosec g101. [#270](https://github.com/vmware/terraform-provider-vmc/pull/270)
    - Fixed indent-error-flow. [#269](https://github.com/vmware/terraform-provider-vmc/pull/269)

## 1.15.1

> Release Date: 2024-04-05

CHORE:

- Updated `google.golang.org/protobuf` to 1.33.0. [#217](https://github.com/vmware/terraform-provider-vmc/pull/217)

## 1.15.0

> Release Date: 2024-01-25

ENHANCEMENT:

- Updated `vsphere-automation-sdk-go/services/vmc` to v0.14.0.
- Added support for disk less instances c6i/m7i. [\#214](https://github.com/vmware/terraform-provider-vmc/pull/214)

## 1.14

> Release Date: 2023-11-24

BREAKING CHANGES:

- Removed support for r5.metal instances. [\#205](https://github.com/vmware/terraform-provider-vmc/pull/205)
- Removed `storage_capacity` field from sddc and cluster resources schema. [\#205](https://github.com/vmware/terraform-provider-vmc/pull/205)

BUG FIXES:

- Fixed failure during SDDC group creation. [\#208](https://github.com/vmware/terraform-provider-vmc/pull/210)

## 1.13.3

> Release Date: 2023-08-21

BUG FIXES:

- Fixed errors when creating and deleting multiple srm nodes. [\#175](https://github.com/vmware/terraform-provider-vmc/pull/195)

## 1.13.2

> Release Date: 2023-08-08

BUG FIXES:

- Fixing errors when reading customer subnets. [\#191](https://github.com/vmware/terraform-provider-vmc/pull/191)

## 1.13.1

> Release Date: 2023-08-04

BUG FIXES:

- Fixed allowing usage of `microsoft_license_config` upon SDDC creation. Reading `microsoft_license_config.academic_license` field. [\#190](https://github.com/vmware/terraform-provider-vmc/pull/190)

DOCUMENTATION:

- Updated documentation. [\#186](https://github.com/vmware/terraform-provider-vmc/pull/186)

CHORE:

- Updated `google.golang.org/grpc`1.53.0. [\#184](https://github.com/vmware/terraform-provider-vmc/pull/184)

## 1.13.0

> Release Date: 2023-02-23

FEATURES:

- Added support for OAuth2.0 app authentication. [\#173](https://github.com/vmware/terraform-provider-vmc/pull/173)

CHORE:

- Fixes for security vulnerabilities.

## 1.12.1

> Release Date: 2023-02-20

BUG FIXES:

- Fixed destroying `sddc_group `times out and then fails on subsequent attempts. [\#172](https://github.com/vmware/terraform-provider-vmc/pull/172)
- Removeed restrictions on 6 hosts minimum in Multi-AZ SDDCs.  [\#171](https://github.com/vmware/terraform-provider-vmc/pull/171)

## 1.12.0

> Release Date: 2022-11-09

FEATURES:

- Added support for SDDC Groups.  [\#163](https://github.com/vmware/terraform-provider-vmc/pull/163)

## 1.11.0

> Release Date: 2022-11-17

FEATURES:

- `num_hosts` property on SDDC now shows and controls the number of hosts in the primary cluster (created by default with an SDDC). Previously there was no way to scale up/down the number of hosts in the primary cluster
- Cluster operations like Create/Update/Destroy on more than one cluster can now be initiated simultaneously, without the need for depends_on=[]
- Error reporting improvements

BUG FIXES:

- EDRS settings are not truly "Optional" [\#151](https://github.com/vmware/terraform-provider-vmc/pull/151)
- Lack of multi-cluster SDDC support in "resourceSddcUpdate" function [\#155](https://github.com/vmware/terraform-provider-vmc/pull/155)
- `vmc_cluster` resource tries to create new clusters simultaneously and fails [\#160](https://github.com/vmware/terraform-provider-vmc/pull/160)

## 1.10.1

> Release Date: 2022-09-20

BUG FIXES:

- Fixed defaults for `enable_edrs`, `edrs_policy_type`, `max_hosts`, and `min_hosts` should not be set to null. [\#94](https://github.com/vmware/terraform-provider-vmc/issues/94)

## 1.10.0

> Release Date: 2022-07-12

FEATURES:

- `vmc_sddc` Added I4I_METAL host instance type support.

BUG FIXES:

- Allow `min_hosts` as low as 2 as per VMC service backend default value. [\#147](https://github.com/vmware/terraform-provider-vmc/pull/147)
- Removed references to R5 host instance type in examples as new deployments are unavailable. [\#146](https://github.com/vmware/terraform-provider-vmc/pull/146)

## 1.9.3

> Release Date: 2022-06-07

BUG FIXES:

- Fixed nil dereference when doing "terraform plan" in some environments. (https://github.com/vmware/terraform-provider-vmc/pull/142)

## 1.9.2

> Release Date: 2022-06-07

BUG FIXES:

- Fixed nil dereference when doing "terraform plan" in some environments (https://github.com/vmware/terraform-provider-vmc/pull/141)

CHORE:

- Updated `hashicorp/terraform-plugin-sdk` to v2.11.0. (https://github.com/vmware/terraform-provider-vmc/pull/140)

## 1.9.1

> Release Date: 2022-03-28

BUG FIXES:

- Example files in `/examples/..` now have a `required_providers` declaration, fixing "terraform init".

## 1.9.0

> Release Date: 2022-03-16

CHORE:

- Updated `github.com/vmware/vsphere-automation-sdk-go/services/vmc` to 0.8.0. [\#130](https://github.com/vmware/terraform-provider-vmc/pull/130)

## 1.8.0

> Release Date: 2021-11-21

ENHANCEMENT:

 - Added firect access to NSX Manager. [\#116](https://github.com/vmware/terraform-provider-vmc/pull/116)

BUG FIXES:

 - Deprecated  `esx_hosts` property in `SddcResourceConfig`. [\#119](https://github.com/vmware/terraform-provider-vmc/pull/119)

CHORE:

- Updated `github.com/vmware/vsphere-automation-sdk-go/services/vmc` to 0.6.0. [\#130](https://github.com/vmware/terraform-provider-vmc/pull/130)
  [\#116](https://github.com/vmware/terraform-provider-vmc/pull/116), [\#118](https://github.com/vmware/terraform-provider-vmc/pull/118)

## 1.7.0

> Release Date: 2021-08-05

BUG FIXES:

- Fixed static check failures caused by references to deprecated functions. [\#109](https://github.com/vmware/terraform-provider-vmc/pull/109)

ENHANCEMENT:

- Upgraded `hashicorp/terraform-plugin-sdk` to v2. [\#108](https://github.com/vmware/terraform-provider-vmc/pull/108)

## 1.6.0

> Release Date: 2021-05-17

BUG FIXES:

- Fix for updating multiple params in SDDC resource. [\#101](https://github.com/vmware/terraform-provider-vmc/pull/101)

FEATURES:

- Added support for SDDC data source. [\#105](https://github.com/vmware/terraform-provider-vmc/pull/105)
- Added support for M1 Mac. [\#107](https://github.com/vmware/terraform-provider-vmc/pull/107)

## 1.5.1

> Release Date: 2021-01-21

BUG FIXES:

- Moved subnet validation from customized diff to create SDDC block. [\#96](https://github.com/vmware/terraform-provider-vmc/pull/96)
- Added zerocloud check for setting intranet MTU uplink. [\#98](https://github.com/vmware/terraform-provider-vmc/pull/98)

DOCUMENTATION:

- Updated help documentation to specify current limitations in import SDDC functionality. [\#97](https://github.com/vmware/terraform-provider-vmc/pull/97)

## 1.5.0

> Release Date: 2021-01-05

BUG FIXES:

- Removed validation on `sddc_type` field in order to allow empty value. [\#72](https://github.com/vmware/terraform-provider-vmc/pull/72)
- Removed default values to fix EDRS configuration error for 1NODE SDDC. [\#83](https://github.com/vmware/terraform-provider-vmc/pull/83)
- Added check to store `vxlan_subnet` information in terraform state file only when `skip_creating_vxlan = false`. [\#86](https://github.com/vmware/terraform-provider-vmc/pull/86)

FEATURES:

- `vmc_sddc`, `vmc_cluster` Added Microsoft licensing configuration to SDDC and cluster resource. [\#71](https://github.com/vmware/terraform-provider-vmc/pull/71)
- Added intranet_uplink_mtu field in `vmc_sddc` resource. [\#88](https://github.com/vmware/terraform-provider-vmc/pull/88)

## 1.4.0

> Release Date: 2020-10-12

BUG FIXES:

- Added check in resourceClusterRead to see if cluster exists and remove cluster information from terraform state file. [\#48](https://github.com/vmware/terraform-provider-vmc/pull/48)
- Added validation check for customer subnet IDs based on the deployment type. [\#54](https://github.com/vmware/terraform-provider-vmc/pull/54)

ENHANCEMENTS:

- Modified Importer State in `vmc_cluster` and `vmc_public_ip` resources for terraform import command. [\#49](https://github.com/vmware/terraform-provider-vmc/pull/49)

FEATURES:

- `vmc_sddc` Added I3EN_METAL host instance type support.  [\#42](https://github.com/vmware/terraform-provider-vmc/pull/42)
- `vmc_sddc` Modified code to enable EDRS policy configuration. [\#43](https://github.com/vmware/terraform-provider-vmc/pull/43)
- `vmc_sddc` Added size parameter in resource schema to enable users to deploy large SDDC. [\#59](https://github.com/vmware/terraform-provider-vmc/pull/59)

## 1.3.0

> Release Date: 2020-06-18

BUG FIXES:

- Modified code to store num_host in resourceSddcRead method. [\#39](https://github.com/vmware/terraform-provider-vmc/pull/39)

FEATURES:

- - **New Resource:** - `vmc_cluster` Added resource for cluster management. [\#25](https://github.com/vmware/terraform-provider-vmc/pull/25)

ENHANCEMENTS:

- Validation check added for Multi-AZ SDDC. [\#35](https://github.com/vmware/terraform-provider-vmc/pull/29)
- Added detailed error handler functions for CRUD operations on resources and data sources. [\#35](https://github.com/vmware/terraform-provider-vmc/pull/29)
- Added documentation for vmc_cluster resource.  [\#26](https://github.com/vmware/terraform-provider-vmc/pull/26)

## 1.2.1

> Release Date: 2020-05-04

BUG FIXES:

- Added instructions for delay needed after SDDC creation for site recovery. [\#21](https://github.com/vmware/terraform-provider-vmc/pull/21)
- Removed capitalized error messages from code. [\#23](https://github.com/vmware/terraform-provider-vmc/pull/23)
- Updated module name in `go.mod`. [\#24](https://github.com/vmware/terraform-provider-vmc/pull/24)

ENHANCEMENTS:

- Updated dependencies version to latest in `go.mod`. [\#20](https://github.com/vmware/terraform-provider-vmc/pull/20)
- Added sample `.tf` file for each resource in examples folder. [\#22](https://github.com/vmware/terraform-provider-vmc/pull/22)

## 1.2.0

> Release Date: 2020-04-03

FEATURES:

- - **New Resource:** - `vmc_site_recovery` Added resource for site recovery management. [\#14](https://github.com/vmware/terraform-provider-vmc/pull/14)
- - **New Resource:** - `vmc_srm_node` Added resource to add SRM instance after site recovery has been activated. [\#14](https://github.com/vmware/terraform-provider-vmc/pull/14)

## 1.1.1

> Release Date: 2020-03-24

BUG FIXES:

- Set ForceNew for `host_instance_type` to `true` in order to enforce SDDC redeploy when `host_instance_type` is changed. [\#5](https://github.com/vmware/terraform-provider-vmc/pull/5)
- Fix for re-creating public IP if it is accidentally delete via console. [\#11](https://github.com/vmware/terraform-provider-vmc/pull/11)
- Updated `variables.tf` for description of fields. [\#10](https://github.com/vmware/terraform-provider-vmc/pull/10)

## 1.1.0

> Release Date: 2020-03-10

FEATURES:

- **New Resource:** - `vmc_sddc`
- **New Resource:** - `vmc_public_ip` [\#43](https://github.com/vmware/terraform-provider-vmc/pull/43)
- **New Data Source:** - `vmc_org`
- **New Data Source:** - `vmc_connected_accounts`
- **New Data Source:** - `vmc_customer_subnets`

ENHANCEMENTS:

- `vmc_sddc`: Added nsxt_reverse_proxy_url to SDDC resource data. [\#23](https://github.com/vmware/terraform-provider-vmc/pull/23)
- `vmc_connected_accounts`: Added filtering to return AWS account ID associated with the account number provided in configuration. [\#30](https://github.com/vmware/terraform-provider-vmc/pull/30)
- `provider.go`: Added `org_id` as a required parameter in terraform schema. [\#38](https://github.com/vmware/terraform-provider-vmc/pull/38)
- `data_source_vmc_customer_subnets.go`: Added `validateFunctions` for sddc and customer subnet resources. [\#41](https://github.com/vmware/terraform-provider-vmc/pull/41)
- `examples/main.tf`: Added expression to convert AWS specific region to VMC region. [\#46](https://github.com/vmware/terraform-provider-vmc/pull/46)

BUG FIXES:

- Moved `main.tf` to `examples` folder. [\#17](https://github.com/vmware/terraform-provider-vmc/pull/17)
- License statement fixed. [\#8](https://github.com/vmware/terraform-provider-vmc/pull/8)
- Implemented connected account data source to return a single account ID associated with the account number. [\#40](https://github.com/vmware/terraform-provider-vmc/pull/40)
- Set ForceNew for `host_instance_type` to true in order to enforce SDDC redeploy when `host_instance_type` is changed. [\#5](https://github.com/vmware/terraform-provider-vmc/pull/5)
