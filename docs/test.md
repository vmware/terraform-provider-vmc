<!--
© Broadcom. All Rights Reserved.
The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
SPDX-License-Identifier: MPL-2.0
-->

<!-- markdownlint-disable first-line-h1 no-inline-html -->

<img src="images/icon-color.png" alt="VMware Cloud on AWS" width="150">

# Testing the Terraform Provider for VMware Cloud on AWS

Testing the Terraform Provider for VMware Cloud on AWS requires having an VMware
Cloud on AWS organization to test against. Generally, the acceptance tests
create real resources, and often cost money to run.

## Configuring Environment Variables

Set required environment variables based on your infrastructure settings.

```sh
$ # clientId and client secret of the test OAuth2.0 app attached to the test organization with at least
$ # "Organization Member" role and service role  on "VMware Cloud on AWS" service that is allowed to deploy SDDCs.
$ # Note: it is recommended to use OAuth2.0 app with the least possible roles (the above mentioned) for testing
$ # purposes.
$ export CLIENT_ID=xxx
$ export CLIENT_SECRET=xxx
$ # Id of a VMC Org in which test SDDC are (to be) placed
$ export ORG_ID=xxxx
$ # Id of an existing SDDC used for SDDC data source (import) test
$ export TEST_SDDC_ID=xxx
$ # Name of above SDDC
$ export TEST_SDDC_NAME=xxx
$ # NSX URL of a non-ZEROCLOUD SDDC, used for real IP testing
$ export NSXT_REVERSE_PROXY_URL=xxx
$ # Account number of a connected to the above Org AWS account, required for test SDDC deployment
$ export AWS_ACCOUNT_NUMBER=xxx
```

## Running the Acceptance Tests

Acceptance tests create real resources, and often cost money to run.

You can run the acceptance tests by running:

```sh
$ make testacc
```

If you want to run against a specific set of tests, run `make testacc` with the
`TESTARGS` parameter containing the run mask. For example:

```sh
$ make testacc TESTARGS="-run=TestAccResourceVmcSddc_basic"
```

Additionally, a limited set of acceptance tests can be run with the ZEROCLOUD
cloud provider, which is much faster and cheaper, while providing decent API
coverage:. For example:

```sh
$ make testacc TESTARGS="-run=TestAccResourceVmcSddcZerocloud"
```
